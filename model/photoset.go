package model

import (
	"context"
	"database/sql"
	"fmt"
	"site_parsing/pagetype"
)

// ParsingIdModel структура
type PhotoSetModel struct {
	dataBase      *sql.DB
	modelCategory *CategoryModel
	modelModel    *ModelModel
	modelPhoto    *PhotoModel
}

// NewParsingIdModel конструктор модели возвращающий указатель на структуру ParsingIdModel
func NewPhotoSetModel(DB *sql.DB) *PhotoSetModel {
	return &PhotoSetModel{
		dataBase:      DB,
		modelPhoto:    NewPhotoModel(DB),
		modelCategory: NewCategoryModel(DB),
		modelModel:    NewModelModel(DB),
	}
}

// CheckID метод по проверки слайса структур после парсинга на наличие в БД
func (p *PhotoSetModel) CheckID(photoSets []pagetype.PhotoSet) (int, error) {
	var id string
	var err error
	fmt.Println("Массив вошедший для проверки в модель", photoSets)
	for key, photoSet := range photoSets {
		photoId := photoSet.Id
		err := p.dataBase.QueryRow("SELECT id FROM parsing_site.photoset where id= $1", photoId).Scan(&id)
		if id == photoId {
			fmt.Printf("на этом id выборка закончилась", id)
			return key, err
		}

		if err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				err = fmt.Errorf("Ошибка в БД при проверки повторения id", err)
				return key, err
			}
		}
	}
	return len(photoSets), err
}

// CheckAndAddId метод по добавлению новых значений в БД
func (p *PhotoSetModel) AddPhotoset(photoId string, dirName string) error {
	err := p.checkPhotoSet(photoId)
	if err != nil {
		err = fmt.Errorf("Ошибка базы данных: %s", err)
		return err
	}
	_, err = p.dataBase.Exec("insert into parsing_site.photoset (id, address_dir, adding_time) values ($1,$2, now())", photoId, dirName)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления в таблицу фотосет: %s", err)
		return err
	}
	return nil
}
func (p *PhotoSetModel) checkPhotoSet(photoId string) error {
	var id string
	err := p.dataBase.QueryRow("SELECT id FROM parsing_site.photoset where id= $1", photoId).Scan(&id)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		err = fmt.Errorf("Ошибка базы данных при проверки %s", err)
		return err
	}
	return err
}

func (p *PhotoSetModel) AddData(photoSetId, address, dirName string, photoSet pagetype.PhotoSet) error {

	ctx := context.Background()
	tx, err := p.dataBase.BeginTx(ctx, nil)
	if err != nil {
		err := fmt.Errorf("Ошибка метода транзакции в БД %v\n\t", err)
		return err
	}
	defer tx.Rollback()

	err = p.AddPhotoset(photoSetId, dirName)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления фотосета в БД %v\n\t", err)
		return err
	}
	err = p.modelPhoto.AddPhoto(photoSetId, address)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления фото в БД %v\n\t", err)
		return err
	}
	for _, valueModel := range photoSet.Models {

		err = p.modelModel.AddModel(valueModel.Name, photoSetId)
		if err != nil {
			err := fmt.Errorf("Ошибка добавления модели в БД %v\n\t", err)
			return err
		}
	}
	for _, valueCategory := range photoSet.Categories {
		err = p.modelCategory.AddCategory(valueCategory.Name, photoSetId)
		if err != nil {
			err := fmt.Errorf("Ошибка категории в БД %v\n\t", err)
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return err
}
