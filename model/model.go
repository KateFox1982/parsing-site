package model

import (
	"database/sql"
	"fmt"
)

// NameCategoryPhotosetModel структура
type ModelModel struct {
	dataBase *sql.DB
}

// NewNameCategoryPhotosetModel конструктор модели возвращающий указатель на структуру NameCategoryPhotosetModel
func NewModelModel(DB *sql.DB) *ModelModel {
	return &ModelModel{
		dataBase: DB,
	}
}

// NameSave метод для сохранения имен в БД
func (n *ModelModel) AddModel(name, photosetId string) error {
	modelId, err := n.checkModel(name)
	if err != nil {
		err = fmt.Errorf("Ошибка выбора категории из БД", err)
		return err
	}
	if modelId == 0 {
		err = n.dataBase.QueryRow("insert into parsing_site.model (name, adding_time) values ($1, now()) RETURNING id", name).Scan(&modelId)
		if err != nil {
			err := fmt.Errorf("Ошибка добавления в таблицу моделей", err)
			return err
		}
	}
	err = n.linkPhotoSet(photosetId, modelId)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления в таблицу фотосет-модель", err)
		return err
	}
	return err
}
func (n *ModelModel) linkPhotoSet(photosetId string, modelId int) error {
	_, err := n.dataBase.Exec("insert into parsing_site.photoset_model (fk_photoset, fk_model) values ($1, $2)", photosetId, modelId)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления в таблицу фотосет-модель", err)
		return err
	}
	return err
}
func (n *ModelModel) checkModel(name string) (int, error) {
	var modelId int
	err := n.dataBase.QueryRow("SELECT id FROM parsing_site.model where name = $1", name).Scan(&modelId)

	if err != nil {
		if err == sql.ErrNoRows {
			return modelId, nil
		}
		err = fmt.Errorf("Ошибка базы данных при проверки", err)
		return modelId, err
	}
	return modelId, err
}
