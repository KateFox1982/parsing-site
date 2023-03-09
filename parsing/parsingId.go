package parsing

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"site_parsing/pagetype"
	"strconv"
	"strings"
)

var watermarks = "logoza.png"

// Parsing структура
type Parser struct {
}

// NewParsing контроллер парсинга  возвращающий указатель на структуру Parsing
func NewParser() *Parser {
	return &Parser{}
}

var website = "https://babesource.com/"

//pageOpen непубличная функция скачивающая страницу сайта
func openPage(pageNum int) (*goquery.Document, error) {
	var webPage string
	numPageString := strconv.Itoa(pageNum)
	if pageNum < 1 {
		webPage = website
	} else {
		webPage = website + "page" + numPageString + ".html"
	}
	res, err := http.Get(webPage)
	if err != nil {
		err := fmt.Errorf("Ошибка получения страницы в функции pageOpen", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err := fmt.Errorf("Ошибка в получении страницы", err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		err := fmt.Errorf("Ошибка метода goquery.NewDocumentFromReader", err)
		return doc, err
	}
	return doc, err
}

// ParsingPage метод по пасингу страницы и поллучению слайса структуры PhotoSet
func (p *Parser) ParsePage(pageNum int) ([]pagetype.PhotoSet, error) {
	var photosetId string
	var newPhotoSets []pagetype.PhotoSet
	doc, err := openPage(pageNum)
	if err != nil {
		err := fmt.Errorf("Ошибка чтения страницы", err)
		return nil, err
	}
	//парсинг основной страницы полуение id, основного фото и url фотосета
	doc.Find(".main-content__card-link").Each(func(j int, tr *goquery.Selection) {
		linkPhotoset, _ := tr.Attr("href")
		tr.Find("img").Each(func(i int, tor *goquery.Selection) {
			linkPhoto, _ := tor.Attr("data-src")
			photosetId = strings.Split(linkPhoto, "/")[4]
			d := pagetype.Photo{Photo: linkPhoto}
			c := pagetype.PhotoSet{Id: photosetId, URL: linkPhotoset, Photo: d}
			newPhotoSets = append(newPhotoSets, c)
		})
	})
	return newPhotoSets, err
}

// GetPhotoset метод по парсингу каждого фотосета, и получению структуры NameCategory
func (p *Parser) GetPhotoset(photoSet pagetype.PhotoSet) ([]pagetype.Model, []pagetype.Category, []pagetype.Photo, error) {
	var sliceName []string
	photoSetLink := photoSet.URL
	res, err := http.Get(photoSetLink)
	if err != nil {
		err := fmt.Errorf("Ошибка чтения страницы в методе GetPhotoset", err)
		return nil, nil, nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err := fmt.Errorf("Ошибка чтения страницы в методе GetPhotoset", err)
		return nil, nil, nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		err := fmt.Errorf("Ошибка чтения страницы goquery в методе GetPhotoset", err)
		return nil, nil, nil, err
	}
	//парсинг имен моделей
	name := doc.Find("title").Text()
	stringNames := strings.Split(name, "-")
	for _, value := range stringNames {
		sliceName = append(sliceName, value)
	}
	for key, valueSlice := range sliceName {
		if key == 0 {
			stringNames := strings.Split(valueSlice, ",")
			for _, valueNames := range stringNames {
				stringName := strings.Split(valueNames, "&")
				for _, valueName := range stringName {
					nameModel := strings.TrimSpace(valueName)
					fmt.Printf("Имена моделей %v\n\t", nameModel)
					model := pagetype.Model{Name: nameModel}
					photoSet.Models = append(photoSet.Models, model)
				}
			}
		}
	}
	//парсинг категорий
	doc.Find(".aside-setting__wrapper-category").Each(func(i int, tg *goquery.Selection) {
		tg.Find("a").Each(func(j int, tag *goquery.Selection) {
			//	Ищем категории
			linkCategories, _ := tag.Attr("href")
			stringCategories := strings.Split(linkCategories, "/")[4]
			category := pagetype.Category{Name: stringCategories}
			photoSet.Categories = append(photoSet.Categories, category)
		})
	})
	//парсинг фото внутри фотосета
	doc.Find(".box-massage__card-link").Each(func(i int, tag *goquery.Selection) {
		if i < 3 {
			link, _ := tag.Attr("href")
			fmt.Printf("Линк фото внутри фотосета %s %s\n", link, i)
			photo := pagetype.Photo{Photo: link}
			photoSet.Photos = append(photoSet.Photos, photo)
		}
	})

	return photoSet.Models, photoSet.Categories, photoSet.Photos, err
}
