package model

type LoginCredentials struct {
	AccountId     string    `json:"account_id" form:"account_id"`
	Password      string    `json:"password" form:"password"`
	PhoneNumber   string    `json:"phone_number" form:"phone_number"`
}

//func (user UserInfo) string () string {
//	return user.UserId + "," + user.UserName + "," + user.UserType
//}
