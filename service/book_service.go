package service

import "mitkid_web/model"

func (s *Service) ListBookByLevel(level int) (books []model.Book, err error) {
	var Query model.Book
	if level == -1 { // 查询全部
		Query = model.Book{}
	} else {
		Query = model.Book{ BookLevel: level }
	}
	return s.dao.GetBookList(Query)
}

// 根据code获取教材
func (s *Service) GetBookById(bookCode string) (book *model.Book, err error) {
	return s.dao.GetBookById(bookCode)
}
