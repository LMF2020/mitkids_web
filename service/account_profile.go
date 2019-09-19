package service

import (
	"bufio"
	"encoding/base64"
	"errors"
	"mime/multipart"
	"mitkid_web/consts"
	"mitkid_web/model"
	"mitkid_web/utils"
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
	profile.AvatarUrl = account.AvatarUrl

	if role == consts.AccountRoleChild {
		profile.School = account.School
	}
	return
}

// 下载头像
func (s *Service) DownloadAvatar(accountId string) (imgUrl string, err error) {

	account, err := s.GetAccountById(accountId)

	if account == nil && err == nil {
		imgUrl = ""
		err = nil
		return
	} else if account != nil {
		imgUrl = account.AvatarUrl
		return
	}
	return
}

// 上传头像
func (s *Service) UploadAvatar(accountId string, imgFile multipart.File, fileHeader *multipart.FileHeader) (err error) {
	// image size, image type
	if fileHeader.Size >= consts.KB * 15 {
		err = errors.New("头像大小不能超过15kb")
		return
	}
	if !utils.VerifyImageFormat(fileHeader.Filename) {
		err = errors.New("头像必须是图片格式")
		return
	}
	buf := make([]byte, fileHeader.Size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	// if create a new image instead of loading from file, encode the image to buffer instead with png.Encode()

	// png.Encode(&buf, image)

	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	// save into db_account

	err = s.dao.UpdateAvatar(accountId, imgBase64Str)

	return nil

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

	if err = s.dao.UpdateAccountInfo(accountInfo); err != nil {
		//tx.Rollback()
		return err
	}

	//tx.Commit()

	return

}
