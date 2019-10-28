package service

import "mitkid_web/model"

func (s *Service) BatchCreateClassPlanS(aid, cid string, planMap map[int]int) error {
	return s.dao.BatchCreateClassPlanS(aid, cid, planMap)
}
func (s *Service) ListClassPlansByClassIdAndAccountId(cid, aid string) (list []model.ClassPlan, err error) {
	return s.dao.ListClassPlansByClassIdAndAccountId(cid, aid)
}
func (s *Service) DeleteClassPlansByClassIdAndAccountId(cid, aid string) error {
	return s.dao.DeleteClassPlansByClassIdAndAccountId(cid, aid)
}
