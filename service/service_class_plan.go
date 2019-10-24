package service

func (s *Service) BatchCreateClassPlanS(cid string, planMap map[int]int) error {
	return s.dao.BatchCreateClassPlanS(cid, planMap)
}
