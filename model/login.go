package model

type LoginForm struct {
	PhoneNumber   string    `json:"phone_number" form:"phone_number" binding:"required"`
	LoginType	  int		`json:"login_type" form:"login_type" binding:"required"`
	Password      string    `json:"password" form:"password"`
	Code      	  string    `json:"code" form:"code"`
}