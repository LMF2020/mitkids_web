package service

func (s *Service) AddChildToClass(id string, childId string) (err error) {
	return s.dao.AddChildToClass(id, childId)
}

func (s *Service) AddChildsToClass(id string, childIds []string) (err error) {
	return s.dao.AddChildsToClass(id, childIds)
}
