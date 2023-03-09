package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"site_parsing/conveyor"
	"site_parsing/pagetype"
	"site_parsing/parsing"
)

const ErrNum = 100

func main() {
	if err := run(); err != nil {
		fmt.Printf("Ошибка:%s \n", err)
	}
}

func run() (err error) {
	//парметры БД: имя пользователя, пароль, имя БД, использование SSL
	DataSourceName := "user=fox password=123 dbname=fix sslmode=disable "
	//соединение с БД postgres
	DB, err := sql.Open("postgres", DataSourceName)
	//err ошибка соединения
	if err != nil {
		log.Printf("Получить ошибку о postgres присоединении: %s", err)
		return err
	}
	//defer Close отсрочка закрытия БД
	defer DB.Close()
	//начальная страница
	pageNum := 1
	//обращение к конструктору пакета parsing
	parser := parsing.NewParser()
	var photoSets []pagetype.PhotoSet
	//Ошибка которая показывает то, что при проверки в БД нашлось совпадение
	var num int
	//обращение к конструктору пакета  controller
	dbConsignor := conveyor.NewDBConsignor(DB)
	for num == 0 && pageNum < 3 {
		var newPhotoSets []pagetype.PhotoSet
		num, newPhotoSets, err = startAddNewPage(parser, dbConsignor, pageNum)
		if err != nil {
			err = fmt.Errorf("Ошибка функции startAddNewPage", err)
			return err
		}
		photoSets = append(photoSets, newPhotoSets...)
		pageNum++
	}
	//создаем новые директории
	err = setNewDir(photoSets)
	if err != nil {
		err = fmt.Errorf("Ошибка создания папок", err)
		return err
	}
	//получения слайсов Категорий-Имен, и слайса Фотографий внутри фотосета
	for key, photoSet := range photoSets {
		photoSet.Models, photoSet.Categories, photoSet.Photos, err = parser.GetPhotoset(photoSet)
		if err != nil {
			err = fmt.Errorf("Ошибка записи имени в БД", err)
			return err
		}
		photoSets[key].Models = photoSet.Models
		photoSets[key].Categories = photoSet.Categories
		photoSets[key].Photos = photoSet.Photos
	}
	err = dbConsignor.AddPhotoSets(photoSets)
	if err != nil {
		err = fmt.Errorf("Ошибка записи категории в БД", err)
		return err
	}
	waterMarkConsignor := conveyor.NewWaterMarkConsignor()
	//отправка для скачивания фотосессий и фото с выставлением watermark и получением времени скачивания
	for key, photoSet := range photoSets {
		waterMarkConsignor.SetPhotosetWaterMark(&photoSet)
		photo := waterMarkConsignor.SetPhotoWaterMark(photoSet.Id, photoSet.Photo.Photo)
		photoSets[key].Photo = photo
		err = dbConsignor.AddPhotoTimeSaving(photoSet.Id, photo)
		if err != nil {
			err = fmt.Errorf("Ошибка добавления времени в БД", err)
			return err
		}
	}
	return err
}

//startAddNewPage- функция по парсингу следующих страниц в случае если на первой странице все значения новые
func startAddNewPage(parser *parsing.Parser, dbConsignor *conveyor.DBConsignor, pageNum int) (int, []pagetype.PhotoSet, error) {
	var num int
	var newPhotoSets []pagetype.PhotoSet
	var err error
	newPhotoSets, err = parser.ParsePage(pageNum)
	if err != nil {
		err := fmt.Errorf("Ошибка чтения сайта", err)
		return ErrNum, newPhotoSets, err
	}
	lenPage := len(newPhotoSets)
	newPhotoSets, num, err = dbConsignor.CheckIDOnly(lenPage, newPhotoSets)
	if err != nil {
		err := fmt.Errorf("Ошибка в базе данных при проверки id", err)
		return ErrNum, newPhotoSets, err
	}
	return num, newPhotoSets, err
}

//setNewDir функция для создания папок фотосессий
func setNewDir(photoSets []pagetype.PhotoSet) error {
	for _, photoSet := range photoSets {
		photosetId := photoSet.Id
		dirName := filepath.Join(pagetype.Dir, photosetId) + string(filepath.Separator)
		err := os.Mkdir(dirName, 0740)

		if err != nil {
			if err == fs.ErrExist {
				continue
			} else {
				err = fmt.Errorf("Ошибка в создании папок", err)
				return err
			}
		}
	}
	return nil
}
