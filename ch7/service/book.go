package service

import "github.com/weenxin/ulitmate_go_notebook_reading/ch7/model"

func (m *Manager) AddBook(book *model.Book) error {
	if err := m.db.Create(book).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) DeleteBook(bookId uint) error {
	if err := m.db.Delete(&model.Book{}, bookId).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) UpdateBook(book *model.Book) error {
	if err := m.db.Updates(book).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) ListBooks(pageNumber, pageSize int) ([]*model.Book, error) {
	var books []*model.Book
	if err := m.db.Limit(pageNumber).Offset(pageSize).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (m *Manager) GetBook(bookId uint) (*model.Book, error) {
	var book model.Book
	if err := m.db.First(&book, bookId).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
