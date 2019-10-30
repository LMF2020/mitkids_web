package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

func (s *Service) AddChildToClass(id string, childId string) (err error) {
	return s.dao.AddChildToClass(id, childId, consts.JoinClassSuccess)
}

func (s *Service) AddChildsToClass(id string, childIds []string) (err error) {
	return s.dao.AddChildsToClass(id, childIds, consts.JoinClassSuccess)
}

/**
学生端申请加入班级，需要满足以下条件：
1.  加入的 class 未开课
2.  加入的class 学生人数小于 < capacity
3.  该学生没有加入其它班级，并且没有申请加入其它班级（即join_class状态为2或者1的）
4.  如满足如上条件，mk_join_class 创建一条记录并设置状态为 1：申请中

学生端撤销加入班级申请，
1.  如果审批成功，不能撤销，只能撤销未审批的
2.  如满足如上条件，直接删除该学生 join_class 记录

管理员同意申请
-- 1. join_class设置状态为2：申请成功
-- 2. 另外 class 实际人数 + 1
管理员拒绝申请
-- 1. join_class设置状态为3：申请失败 --- 3天后job会删掉状态为3的所有记录
管理员撤销申请
-- 1.管理员有权限撤销学生的申请，并且有权限撤销加入的学生
-- 2. 撤销申请：直接删除该学生 join_class 记录
*/

// 申请加入班级
func (s *Service) ApplyJoiningClass(childId, classId string, ctx *gin.Context, plansMap map[int]int) error {
	c, err := s.dao.GetClassById(classId)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("班级不存在")
	}
	if c.Status == consts.ClassInProgress {
		return errors.New("班级已开课")
	}
	if c.Status == consts.ClassEnd {
		return errors.New("班级已关闭")
	}
	if c.ChildNumber >= c.Capacity {
		return errors.New("班级学生人数已满")
	}
	joinCls, err := s.dao.GetJoiningClass(classId, childId, consts.JoinClassInProgress)
	if joinCls != nil && err == nil { // 存在记录
		return errors.New("您仍有申请在处理")
	}
	joinCls, err = s.dao.GetJoiningClass(classId, childId, consts.JoinClassSuccess)
	if joinCls != nil && err == nil {
		return errors.New("已加入班级，不能重复申请")
	}

	// 处理失败的case,允许继续申请
	joinCls, err = s.dao.GetJoiningClass(classId, childId, consts.JoinClassFail)
	if joinCls != nil && err == nil {
		err = s.addOrUpdateChildToClass(childId, classId, plansMap, false)
		return err
	}
	// 插入申请记录
	err = s.addOrUpdateChildToClass(childId, classId, plansMap, true)
	if err != nil {
		return err
	}
	log.Logger.WithField("class id", classId).WithField("child id", childId).Info("申请加入班级")
	return nil
}

func (s *Service) addOrUpdateChildToClass(childId, classId string, plansMap map[int]int, justAdd bool) (err error) {
	if err := s.checkAndUpdatePlanClassUsed(plansMap, childId, classId); err != nil {
		return err
	}
	if justAdd {
		err = s.dao.AddChildToClass(classId, childId, consts.JoinClassInProgress)
	} else {
		err = s.dao.UpdateJoinClassStatus(childId, classId, consts.JoinClassInProgress)
	}
	return err
}

func (s *Service) checkAndUpdatePlanClassUsed(plansMap map[int]int, accountId, classId string) (err error) {
	planIds := make([]int, len(plansMap))
	i := 0
	countUserClass := 0
	for k, v := range plansMap {
		planIds[i] = k
		countUserClass += v
		i++
	}
	var countOC int = 0
	if countOC, err = s.CountClassOccurs(classId); err != nil {
		return err
	}
	if countUserClass != countOC {
		return errors.New("plan 数量和班级课时不符合")
	}
	plans, err := s.ListPlanByPlanIds(planIds)
	if err != nil {
		//api.Fail(c, http.StatusBadRequest, err)
		return err
	}
	if len(planIds) != len(plans) {
		for _, planItem := range plans {
			if _, ok := plansMap[planItem.PlanId]; ok {
				delete(plansMap, planItem.PlanId)
			}
		}
		NonexistPlans := make([]int, 0, 0)
		for k, _ := range plansMap {
			NonexistPlans = append(NonexistPlans, k)
		}
		return errors.New(fmt.Sprintf("plan_ids:%v 不存在", NonexistPlans))
	}
	for _, planItem := range plans {
		if plansMap[planItem.PlanId]+planItem.UsedClass > planItem.PlanTotalClass {
			//api.Failf(c, http.StatusBadRequest, )
			return errors.New(fmt.Sprintf("plans:%d一共有%d课时,已经使用%d课时,无法再分配%d", planItem.PlanId, planItem.PlanTotalClass, planItem.UsedClass, plansMap[planItem.PlanId]))
		}
	}
	if err := s.BatchUpdatePlanUsedClass(accountId, plansMap); err != nil {
		return err
	}
	if err := s.BatchCreateClassPlanS(accountId, classId, plansMap); err != nil {
		return err
	}

	return nil
}

// 撤销申请加入班级
func (s *Service) CancelJoiningClass(childId, classId string) (err error) {
	c, err := s.dao.GetClassById(classId)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("班级不存在")
	}

	if child, err := s.GetAccountById(childId); err != nil {
		return errors.New("系统查询失败")
	} else if child == nil && err == nil {
		return errors.New("学生账号不存在")
	} else if child != nil && child.AccountRole != consts.AccountRoleChild {
		return errors.New("学生账号不存在")
	}

	joinCls, err := s.dao.GetJoiningClass(classId, childId, consts.JoinClassSuccess)
	if joinCls != nil && err == nil {
		return errors.New("审批成功，不能撤销")
	} else if err != nil {
		return err
	}

	joinCls, err = s.dao.GetJoiningClass(classId, childId, consts.JoinClassInProgress)
	if joinCls != nil && err == nil {
		//删除 预约占用的plan
		var classPlans []model.ClassPlan
		if classPlans, err = s.ListClassPlansByClassIdAndAccountId(classId, childId); err != nil {
			return errors.New("撤销失败")
		}
		planMap := make(map[int]int)
		for _, plan := range classPlans {
			planMap[plan.PlanId] = -plan.UsedClass
		}
		if err := s.BatchUpdatePlanUsedClass(childId, planMap); err != nil {
			return errors.New("撤销失败")
		}
		if err := s.DeleteClassPlansByClassIdAndAccountId(classId, childId); err != nil {
			return errors.New("撤销失败")
		}
		if err = s.dao.DeleteJoiningClass(childId, joinCls.ClassId); err != nil {
			return errors.New("撤销失败")
		}
	}
	return nil
}

// 根据ClassID获取学生列表id
func (s *Service) ListClassChildIdsByClassId(cid string) (ChildIds []string, err error) {
	return s.dao.ListClassChildIdsByClassId(cid)
}

// 根据ClassID获取学生列表
func (s *Service) ListClassChildByClassId(cid string) (ChildIds []model.AccountInfo, err error) {
	return s.dao.ListClassChildByClassId(cid)
}

func (s *Service) UpdateJoinClassStatus(studentId, classId string, status int) error {
	return s.dao.UpdateJoinClassStatus(studentId, classId, status)
}
func (s *Service) checkJoiningClass(childId, classId string) (c *model.Class, child *model.AccountInfo, join *model.JoinClass, err error) {
	c, err = s.GetClassById(classId)
	if err != nil {
		return
	}
	if c == nil {
		err = errors.New("班级不存在")
		return
	}
	if c.Status == consts.ClassInProgress {
		err = errors.New("班级已开课")
		return
	}
	if c.Status == consts.ClassEnd {
		err = errors.New("班级已关闭")
		return
	}

	if child, err := s.GetAccountById(childId); err != nil {
		err = errors.New("系统查询失败")
	} else if child == nil && err == nil {
		err = errors.New("学生账号不存在")
	} else if child != nil && child.AccountRole != consts.AccountRoleChild {
		err = errors.New("学生账号不存在")
	}

	join, err = s.GetJoinClassById(classId, childId)
	if err != nil {
		return
	}
	if join == nil {
		err = errors.New("学生申请加入班级不存在")
		return
	}
	return
}

//admin 同意
func (s *Service) ApproveJoiningClass(classId, childId string) (err error) {
	c, _, join, err := s.checkJoiningClass(childId, classId)
	if err != nil {
		return
	}
	if join.Status == consts.JoinClassSuccess {
		return errors.New("约课申请已经被批准过")
	}
	if c.ChildNumber == c.Capacity || c.ChildNumber > c.Capacity {
		return errors.New("班级学生数量已满")
	}
	plans, err := s.ListClassPlansByClassIdAndAccountId(classId, childId)
	if err != nil {
		return
	}
	count := 0
	for _, plan := range plans {
		count += plan.UsedClass
	}
	countCo, err := s.CountClassOccurs(classId)
	if err != nil {
		return
	}
	if countCo != count {
		return errors.New("批准失败，学生约课后,课程数量有变更，建议学生重新约课")
	}
	err = s.UpdateJoinClassStatus(childId, classId, consts.JoinClassSuccess)
	if err != nil {
		return
	}
	err = s.UpdateClassChildNum(classId, 1)
	if err != nil {
		return
	}
	return
}

// 根据学生ID查询申请班级
func (s *Service) GetJoinClassById(classId, studentId string) (join *model.JoinClass, err error) {
	return s.dao.GetJoinClassById(classId, studentId)
}

func (s *Service) UpdateClassChildNum(classId string, update int) (err error) {
	return s.dao.UpdateClassChildNum(classId, update)
}

//admin 拒绝
func (s *Service) RefuseJoiningClass(classId, childId string) (err error) {
	_, _, join, err := s.checkJoiningClass(childId, classId)
	if err != nil {
		return
	}
	if join.Status == consts.JoinClassFail {
		return errors.New("约课申请已经被拒绝过")
	}

	//todo 删除plan 占用
	err = s.UpdateJoinClassStatus(childId, classId, consts.JoinClassFail)
	if err != nil {
		return
	}
	if join.Status == consts.JoinClassSuccess {
		err = s.UpdateClassChildNum(classId, -1)
	}
	return
}

//admin 修改状态为申请中
func (s *Service) ChangeToApplyJoiningClass(classId, childId string) (err error) {
	_, _, join, err := s.checkJoiningClass(childId, classId)
	if err != nil {
		return
	}
	err = s.UpdateJoinClassStatus(classId, childId, consts.JoinClassInProgress)
	if err != nil {
		return
	}
	if join.Status == consts.JoinClassSuccess {
		err = s.UpdateClassChildNum(classId, -1)
	}
	return

}

func (s *Service) PageListApplyClassChild(pageNumber, pageSize, status int, query string) ([]model.ApplyClassChild, error) {
	offset := (pageNumber - 1) * pageSize
	return s.dao.PageListApplyClassChild(offset, pageSize, status, query)
}

func (s *Service) CountApplyClassChild(status int, query string) (int, error) {
	return s.dao.CountApplyClassChild(status, query)
}

// 删除学生约课申请记录
func (s *Service) DeleteJoiningClasses(classId string, studentIds []string) (err error) {
	return s.dao.DeleteJoiningClasses(classId, studentIds)
}

// 删除学生约课申请记录
func (s *Service) DeleteJoiningClassesByClassId(classId string) (err error) {
	return s.dao.DeleteJoiningClassesByClassId(classId)
}
