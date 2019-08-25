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

	var joinCls []model.JoinClass
	joinCls, _ = s.dao.ListJoiningClass(childId, consts.JoinClassInProgress)
	if joinCls != nil && err == nil { // 存在记录
		return errors.New("已有加入班级的申请")
	}
	joinCls, _ = s.dao.ListJoiningClass(childId, consts.JoinClassSuccess)
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
func (s *Service) CancelJoiningClass(childId, classId string) error {
	var joinCls []model.JoinClass
	c, err := s.dao.GetClassById(classId)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("班级不存在")
	}
	joinCls, _ = s.dao.ListJoiningClass(childId, consts.JoinClassSuccess)
	if joinCls != nil && err == nil {
		return errors.New("审批成功，不能撤销")
	}

	joinCls, _ = s.dao.ListJoiningClass(childId, consts.JoinClassInProgress)
	if joinCls != nil && err == nil {
		if len(joinCls) > 0 {
			cls := joinCls[0]
			if err = s.dao.DeleteJoiningClass(childId, cls.ClassId); err != nil {
				return errors.New("撤销失败")
			}
		}
	}
	return nil
}
