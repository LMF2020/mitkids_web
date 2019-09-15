package model

import "mitkid_web/consts"

type Book struct {
	BookCode  string `json:"book_code" form:"book_code" gorm:"primary_key" validate:"required"`
	BookName  string `json:"book_name" form:"book_name" validate:"required"`
	BookLink  string `json:"book_link" form:"book_link"`
	BookLevel int    `json:"book_level" form:"book_level" validate:"required"`
}

// 定义表名
func (book *Book) TableName() string {
	return consts.TABLE_BOOK
}
