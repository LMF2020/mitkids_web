package model

type LoginCredentials struct {
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	AccountType uint   `form:account_type json:"account_type"`
}

//func (user UserInfo) string () string {
//	return user.UserId + "," + user.UserName + "," + user.UserType
//}
