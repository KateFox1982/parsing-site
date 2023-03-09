package model

import (
	"database/sql"
	"fmt"
)

// NameCategoryPhotosetModel структура
type CategoryModel struct {
	dataBase *sql.DB
}

// NewNameCategoryPhotosetModel конструктор модели возвращающий указатель на структуру NameCategoryPhotosetModel
func NewCategoryModel(DB *sql.DB) *CategoryModel {
	return &CategoryModel{
		dataBase: DB,
	}
}

// CategorySave для сохранения категорий в БД
func (p *CategoryModel) AddCategory(category, photosetId string) error {
	categoryId, err := p.checkCategory(category)
	if err != nil {
		err = fmt.Errorf("Ошибка выбора категории из БД", err)
		return err
	}

	if categoryId == 0 {
		err := p.dataBase.QueryRow("insert into parsing_site.category(name) values ($1) RETURNING id", category).Scan(&categoryId)
		if err != nil {
			err := fmt.Errorf("Ошибка добавления в таблицу категория", err)
			return err
		}
	}
	err = p.linkPhotoSet(photosetId, categoryId)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления в таблицу фотосет-категория", err)
		return err
	}
	return err
}

func (p *CategoryModel) linkPhotoSet(photosetId string, id int) error {
	_, err := p.dataBase.Exec("insert into parsing_site.photoset_category (fk_photoset, fk_category) values ($1, $2)", photosetId, id)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления в таблицу фотосет-категория", err)
		return err
	}
	return err
}
func (p *CategoryModel) checkCategory(category string) (int, error) {
	var id int
	err := p.dataBase.QueryRow("select id from parsing_site.category where name=$1", category).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return id, nil
		} else {
			err = fmt.Errorf("Ошибка базы данных при проверки", err)
			return id, err
		}
	}

	return id, err
}
