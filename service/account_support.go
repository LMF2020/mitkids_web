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

func (s *Service) PageListContacts(contact *model.Contact, pageNumber, pageSize int) (result []model.Contact, err error) {
	offset := (pageNumber - 1) * pageSize
	result, err = s.dao.PageListContact(*contact, offset, pageSize)
	return
}

func (s *Service) TotalContact(contact *model.Contact) (count int, err error) {
	count, err = s.dao.TotalContact(*contact)
	return
}
