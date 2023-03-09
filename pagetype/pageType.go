package pagetype

import (
	"time"
)

// Photo Структура данных фото
type Photo struct {
	Id       int       `json:"id" :"id"`
	Photo    string    `json:"photo" :"photo"`
	Address  string    `json:"address" :"address"`
	TimeDown time.Time `json:"timeDown" :"timeDown"`
	TimeSave time.Time `json:"timeSave" :"timeSave"`
}

// PhotoSet структура данных фотосета
type PhotoSet struct {
	Id         string     `json:"id" :"id"`
	Photo      Photo      `json:"photo" :"photo"`
	Models     []Model    `json:"models" :"models"`
	Categories []Category `json:"categories" :"categories"`
	Photos     []Photo    `json:"photos" :"photos"`
	URL        string     `json:"url" :"url"`
}

// Model структура данных модели
type Model struct {
	Id   int    `json:"id" :"id"`
	Name string `json:"name" :"name"`
}

// Category струтрура данных структур
type Category struct {
	Id   int    `json:"id" :"id"`
	Name string `json:"name" :"name"`
}

var Dir = "photoset"

const Watermarks = "logoza.png"
