package service

import "mitkid_web/model"

func (s *Service) GetAccountByPhoneNumber(number string) (account *model.AccountInfo, err error) {
	return s.dao.GetAccountByPhoneNumber(number)
}

func (s *Service) GetAccountById(id string) (account *model.AccountInfo, err error) {
	return s.dao.GetAccountById(id)
}
