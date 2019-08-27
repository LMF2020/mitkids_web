package service

import (
	"mitkid_web/model"
	"mitkid_web/utils/log"
)

// 查询学生profile信息
func (s *Service) GetChildProfileById(account *model.AccountInfo) (profile *model.ChildProfilePoJo, err error) {
	profile = &model.ChildProfilePoJo{}

	profile.AccountId = account.AccountId
	profile.PhoneNumber = account.PhoneNumber
	profile.AccountName = account.AccountName
	profile.Gender = account.Gender
	profile.Age = account.Age
	profile.Birth = account.Birth
	profile.Address = account.Address
	profile.Province = account.Province
	profile.City = account.City
	profile.District = account.District
	profile.Email = account.Email

	if child, err := s.dao.GetChildProfileById(account.AccountId); err != nil {
		log.Logger.WithField("child account", account.AccountName).Error("学生profile信息不存在")
	} else if child != nil {
		// 学校信息
		profile.School = child.School
	}
	return
}

/**
更新学生profile信息,
但是, 手机号码,密码不能更新，需要另外的接口单独更新
*/
func (s *Service) UpdateChildProfile(profile model.ChildProfilePoJo) (err error) {

	accountInfo := model.AccountInfo{
		AccountId:   profile.AccountId,
		AccountName: profile.AccountName,
		Gender:      profile.Gender,
		Age:         profile.Age,
		Birth:       profile.Birth,
		Address:     profile.Address,
		Province:    profile.Province,
		City:        profile.City,
		District:    profile.Province,
		Email:       profile.Email,
	}

	accountProfile := model.AccountChild{
		AccountId: profile.AccountId,
		School:    profile.School,
	}

	// 提交事务
	tx := s.dao.DB.Begin()

	if err = s.dao.UpdateChildAccount(accountInfo); err != nil {
		tx.Rollback()
		return
	}
	if err = s.dao.UpdateChildProfile(accountProfile); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	return

}

// create account child profile
func (s *Service) CreateChildProfile (id string) error {
	return s.dao.AddChildProfile(id)
}