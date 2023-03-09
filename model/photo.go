package model

import (
	"database/sql"
	"fmt"
	"site_parsing/pagetype"
)

//PhotoModel структура
type PhotoModel struct {
	dataBase *sql.DB
}

// NewPhotoModel конструктор модели возвращающий указатель на структуру PhotoModel
func NewPhotoModel(DB *sql.DB) *PhotoModel {
	return &PhotoModel{
		dataBase: DB,
	}
}

func (p *PhotoModel) AddTime(photosetId string, photo pagetype.Photo) error {
	_, err := p.dataBase.Exec("update parsing_site.photo set saving_time=$1 where fk_photoset=$2", photo.TimeSave, photosetId)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления времени %v\n\t", err)
		return err
	}
	_, err = p.dataBase.Exec("update parsing_site.photo set download_time=$1 where fk_photoset=$2", photo.TimeDown, photosetId)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления времени %v\n\t", err)
		return err
	}
	return err
}

// AddPhoto метод добавдения фотов БД
func (p *PhotoModel) AddPhoto(photoSetId, address string) error {
	var id int
	err := p.checkPhoto(address)
	if err != nil {
		err = fmt.Errorf("Ошибка поиска в таблице фото", err)
		return err
	}
	err = p.dataBase.QueryRow("insert into parsing_site.photo ( address_file, receiving_time, fk_photoset) values ($1,now(), $2)  RETURNING id", address, photoSetId).Scan(&id)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления в таблицу фото", err)
		return err
	}
	err = p.linkPhotoSet(photoSetId, id)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления fk_photo %v\n\t", err)
		return err
	}
	return err
}

func (p *PhotoModel) linkPhotoSet(photoSetId string, id int) error {
	_, err := p.dataBase.Exec("update parsing_site.photoset set fk_photo=$1 where id=$2", id, photoSetId)
	if err != nil {
		err := fmt.Errorf("Ошибка добавления fk_photo %v\n\t", err)
		return err
	}
	return err
}
func (p *PhotoModel) checkPhoto(address string) error {
	var id int
	err := p.dataBase.QueryRow("SELECT id FROM parsing_site.photo where address_file= $1", address).Scan(&id)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		err = fmt.Errorf("Ошибка забы данных при проверки", err)
		return err
	}
	return err
}
