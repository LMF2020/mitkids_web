package service

import (
	"errors"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

func (s *Service) AddChildToClass(id string, childId string) (err error) {
	return s.dao.AddChildToClass(id, childId)
}

func (s *Service) AddChildsToClass(id string, childIds []string) (err error) {
	return s.dao.AddChildsToClass(id, childIds)
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
func (s *Service) ApplyJoiningClass(childId, classId string) error {
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
		return errors.New("已有加入班级的申请")
	}
	joinCls, err = s.dao.GetJoiningClass(classId, childId, consts.JoinClassSuccess)
	if joinCls != nil && err == nil {
		return errors.New("已加入班级，不能重复申请")
	}

	// 插入申请记录
	err = s.dao.AddChildToClass(classId, childId)
	if err != nil {
		return err
	}

	log.Logger.WithField("class id", classId).WithField("child id", childId).Info("正在申请加入班级")

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
	}

	joinCls, err = s.dao.GetJoiningClass(classId, childId, consts.JoinClassInProgress)
	if joinCls != nil && err == nil {
		if err = s.dao.DeleteJoiningClass(childId, joinCls.ClassId); err != nil {
			return errors.New("撤销失败")
		}
	}
	return nil
}

// 根据ClassID获取学生列表
func (s *Service) ListClassChildByClassId(cid string) (ChildIds []string, err error) {
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
		err =  errors.New("学生账号不存在")
	} else if child != nil && child.AccountRole != consts.AccountRoleChild {
		err =  errors.New("学生账号不存在")
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
		return
	}
	if c.ChildNumber == c.Capacity || c.ChildNumber > c.Capacity {
		return errors.New("班级学生数量已满")
	}

	err = s.UpdateJoinClassStatus(classId, childId, consts.JoinClassSuccess)
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
		return
	}
	err = s.UpdateJoinClassStatus(classId, childId, consts.JoinClassFail)
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
