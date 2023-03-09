package conveyor

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"path/filepath"
	"site_parsing/model"
	"site_parsing/pagetype"
)

// DBController структура пакета controller
type DBConsignor struct {
	modelPhotoSet             *model.PhotoSetModel
	modelPhoto                *model.PhotoModel
	modelNameCategoryPhotoset *model.ModelModel
	modelCategory             *model.CategoryModel
}

// NewPDBController конструктор пакета controller
func NewDBConsignor(DB *sql.DB) *DBConsignor {
	return &DBConsignor{
		modelPhoto:                model.NewPhotoModel(DB),
		modelPhotoSet:             model.NewPhotoSetModel(DB),
		modelNameCategoryPhotoset: model.NewModelModel(DB),
		modelCategory:             model.NewCategoryModel(DB),
	}
}

// CheckIDOnly метод получающий структуру FirstPage, проверяющий наличие в БД id, в случае если на странице все значения новые,
//к текущей странице прибавляется еще одна для парсинга
func (d *DBConsignor) CheckIDOnly(pageLen int, photoSets []pagetype.PhotoSet) ([]pagetype.PhotoSet, int, error) {
	var m = 1
	key, err := d.modelPhotoSet.CheckID(photoSets)
	if err != nil {
		err = fmt.Errorf("Ошибка в базе данных при проверки", err)
		return nil, key, err
	}
	if key == pageLen {
		return photoSets, 0, err
	}
	if key == 0 {
		fmt.Printf("Нет обновлений")
		photoSets = photoSets[0:key]
		return photoSets, m, err
	}
	if key < pageLen {
		photoSets = photoSets[0:key]
		return photoSets, m, err
	}
	return photoSets, m, err
}

// PhotoSetAdd метод по отправки информации из структуры PhotoSet в БД
func (d *DBConsignor) AddPhotoSets(photoSets []pagetype.PhotoSet) error {
	for _, photoSet := range photoSets {
		photoSetId := photoSet.Id
		dirName := filepath.Join(pagetype.Dir, photoSetId) + string(filepath.Separator)
		address := filepath.Join(dirName, photoSetId)
		err := d.modelPhotoSet.AddData(photoSetId, address, dirName, photoSet)
		if err != nil {
			err := fmt.Errorf("Ошибка записи значений в таблицу фото", err)
			return err
		}
	}
	return nil
}

// PhotoAddingTime метод по добавлению времени в таблицу Photo после выставления воттермарки
func (d *DBConsignor) AddPhotoTimeSaving(photoSetId string, photo pagetype.Photo) error {
	err := d.modelPhoto.AddTime(photoSetId, photo)
	if err != nil {
		fmt.Printf("Ошибка записи значений в таблицу времени", err)
		return err
	}
	return err
}
