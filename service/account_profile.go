package service

import (
	"mitkid_web/consts"
	"mitkid_web/model"
)

// 根据role 查询账号profile信息
func (s *Service) GetProfileByRole(account *model.AccountInfo, role int) (profile *model.ProfilePoJo, err error) {
	profile = &model.ProfilePoJo{}

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

	if role == consts.AccountRoleChild {
		profile.School = account.School
	}

	return
}

/**
更新账户profile信息,
但是, 手机号码,密码不能更新，需要另外的接口单独更新
*/
func (s *Service) UpdateProfileByRole(profile model.ProfilePoJo, role int) (err error) {

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

	if role == consts.AccountRoleChild {
		accountInfo.School = profile.School
	}

	// 提交事务
	//tx := s.dao.DB.Begin()

	if err = s.dao.UpdateAccount(accountInfo); err != nil {
		//tx.Rollback()
		return err
	}

	//tx.Commit()

	return

}