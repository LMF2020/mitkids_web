package controllers

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"mitkid_web/consts"
	"mitkid_web/consts/errorcode"
	"mitkid_web/controllers/api"
	"mitkid_web/model"
	"mitkid_web/utils"
	"mitkid_web/utils/cache"
	"mitkid_web/utils/fileUtils"
	"mitkid_web/utils/log"
	"net/http"
)

// 设置默认昵称
func _setDefaultName(phone string, role uint) string {
	if s.IsRoleChild(int(role)) {
		return fmt.Sprintf("学员%s", phone)
	}

	if s.IsRoleChineseTeacher(int(role)) {
		return fmt.Sprintf("中教%s", phone)
	}

	if s.IsRoleForeTeacher(int(role)) {
		return fmt.Sprintf("外教%s", phone)
	}

	if s.IsRoleCorp(int(role)) {
		return fmt.Sprintf("合作%s", phone)
	}

	return ""
}

func CreateAccount(c *gin.Context, role uint) {

	var account model.AccountInfo
	if err := c.ShouldBind(&account); err == nil {

		account.AccountRole = role
		account.AccountStatus = consts.AccountStatusNormal
		account.AccountType = consts.AccountTypePaid

		// 参数校验：手机号,验证码,年龄,密码,性别
		if err := utils.ValidateParam(account); err != nil {
			api.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if _tmpAcc, err := s.GetAccountByPhoneNumber(account.PhoneNumber); err != nil {
			api.Fail(c, http.StatusInternalServerError, "系统内部错误")
			return
		} else if _tmpAcc != nil {
			api.Fail(c, errorcode.USER_ALREADY_EXIS, "手机号已注册")
			return
		}

		// 注册验证码校验：
		if account.Code == "" {
			api.Fail(c, http.StatusBadRequest, "验证码不能为空")
			return
		}

		codeKey := fmt.Sprintf(consts.CodeRegPrefix, account.PhoneNumber) // 注册验证码前缀
		it, _ := cache.Client.Get(codeKey)
		if it == nil || it.Key != codeKey || string(it.Value) != account.Code {
			api.Fail(c, errorcode.VERIFY_CODE_ERR, "验证码错误")
			return
		}

		// 设置默认姓名
		account.AccountName = _setDefaultName(account.PhoneNumber, role)

		// 创建账号信息
		if err := s.CreateAccount(&account); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}

		log.Logger.WithField("account", account).Info("create account successfully")

		api.Success(c, "账号创建成功")
	} else {
		log.Logger.WithField("account", account).Error("create account failed")
		api.Fail(c, http.StatusBadRequest, err.Error())
	}
}

// 管理员创建账号
func AdminCreateAccount(c *gin.Context) {

	var account model.AccountInfo
	if err := c.ShouldBind(&account); err == nil {

		//account.AccountRole = role
		account.AccountStatus = consts.AccountStatusNormal
		account.AccountType = consts.AccountTypePaid

		// 参数校验
		if account.AccountRole == 0 {
			api.Fail(c, http.StatusBadRequest, "缺少角色参数")
			return
		}

		if !s.IsRoleTeacher(int(account.AccountRole)) {
			api.Fail(c, http.StatusBadRequest, "角色不正确")
			return
		}

		if account.AccountName == "" {
			api.Fail(c, http.StatusBadRequest, "缺少教师姓名")
			return
		}

		if account.AccountRole == consts.AccountRoleTeacher && account.PhoneNumber == "" {
			api.Fail(c, http.StatusBadRequest, "缺少中教手机号")
			return
		}

		if account.AccountRole == consts.AccountRoleForeignTeacher && account.Email == "" {
			api.Fail(c, http.StatusBadRequest, "缺少外教邮箱")
			return
		}

		if account.Email != "" && !utils.VerifyEmailFormat(account.Email) {
			api.Fail(c, http.StatusBadRequest, "请填写正确的邮箱")
			return
		}

		if account.Password == "" { // 密码未提供, 默认密码123456
			account.Password = "123456"
		}

		if account.Gender != consts.AccountGenderMale &&  account.Gender != consts.AccountGenderFemale {
			api.Fail(c, http.StatusBadRequest, "缺少教师性别")
			return
		}

		if account.Age == 0 {
			api.Fail(c, http.StatusBadRequest, "缺少教师年龄")
			return
		}

		if account.PhoneNumber != "" { // 如果电话号码不为空，需要校验电话号码
			if _tmpAcc, err := s.GetAccountByPhoneNumber(account.PhoneNumber); err != nil {
				api.Fail(c, http.StatusInternalServerError, "系统内部错误")
				return
			} else if _tmpAcc != nil {
				api.Fail(c, errorcode.USER_ALREADY_EXIS, "手机号已注册")
				return
			}
		}

		// 设置默认姓名
		// account.AccountName = _setDefaultName(account.PhoneNumber, role)

		// 创建账号信息
		if err := s.CreateAccount(&account); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}

		log.Logger.WithField("account", account).Info("create account successfully")

		api.Success(c, "教师账号创建成功")
	} else {
		log.Logger.WithField("account", account).Error("admin create account failed")
		api.Fail(c, http.StatusBadRequest, err.Error())
	}
}

// 学生、教师头像下载
func UserAvatarDownloadHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	imgUrl, err := s.DownloadAvatar(accountId)
	if err != nil {
		api.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	api.Success(c, imgUrl)
}

// 学生个人资料更新
func AccountPicUpdateHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	accountId := claims["AccountId"].(string)
	urlMap := make(map[string]string)
	file, header, err := c.Request.FormFile("file")
	if file != nil {
		if err != nil {
			api.Fail(c, http.StatusBadRequest, "头像更新失败")
			return
		}
		//文件的名称
		filename := header.Filename
		avatarUrl, err := fileUtils.UpdateUserPic(accountId, filename, file)
		if err != nil {
			api.Fail(c, http.StatusBadRequest, "头像更新失败")
			return
		}
		account := model.AccountInfo{AccountId: accountId,
			AvatarUrl: avatarUrl}
		if err = s.UpdateAccountInfo(account); err != nil {
			api.Fail(c, http.StatusInternalServerError, err.Error())
			return
		}
		urlMap["avatar_url"] = avatarUrl
		api.Success(c, urlMap)
		return
	}
	api.Success(c, urlMap)
	return
}
