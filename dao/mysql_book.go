package dao

import (
	"github.com/jinzhu/gorm"
	"mitkid_web/model"
)

// 根据条件查询book
func (d *Dao) GetBookList(query model.Book) (books []model.Book, err error) {

	books = []model.Book{}
	if err = d.DB.Where(query).Find(&books).Error; err == gorm.ErrRecordNotFound {
		err = nil
		books = nil
	}
	return
}

// 根据ID获取book
func (d *Dao) GetBookById(bookCode string) (book *model.Book, err error) {
	book = &model.Book{}
	if err := d.DB.Where("book_code = ?", bookCode).First(book).Error; err == gorm.ErrRecordNotFound {
		err = nil
		book = nil
	}
	return
}
