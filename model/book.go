package model

import "mitkid_web/consts"

type Book struct {
	BookCode  string `json:"book_code" form:"book_code" gorm:"primary_key" validate:"required"`
	BookName  string `json:"book_name" form:"book_name" validate:"required"`
	BookLink  string `json:"book_link" form:"book_link"`
	BookLevel int    `json:"book_level" form:"book_level" validate:"required"`
	BookPhase int    `json:"book_phase" form:"book_phase"`  // 阶段
	BookUrl   string `json:"book_url" form:"book_url" gorm:"-"`	// MitKids/Level_01/Phase_01/Unit_01/Lesson_01
	BookTitle string `json:"book_title" form:"book_title" gorm:"-"`  // Unit1 Lesson1
}

// 定义表名
func (book *Book) TableName() string {
	return consts.TABLE_BOOK
}
