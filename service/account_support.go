package service

import (
	"errors"
	"mitkid_web/model"
)

func (s *Service) AddContact(contact *model.Contact) (err error) {
	if exists, _ := s.dao.GetContact(model.Contact{PhoneNumber: contact.PhoneNumber}); exists != nil {
		err = errors.New("联系人已存在")
		return
	}
	err = s.dao.AddContact(contact)
	return
}

func (s *Service) UpdateContact(contact *model.Contact) (err error) {
	if exists, _ := s.dao.GetContact(model.Contact{PhoneNumber: contact.PhoneNumber}); exists == nil {
		err = errors.New("联系人不存在")
		return
	}
	err = s.dao.UpdateContact(contact)
	return
}

func (s *Service) ListContacts(contact *model.Contact) (result []model.Contact, err error) {
	result, err = s.dao.ListContact(*contact)
	return
}
