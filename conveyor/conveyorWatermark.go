package conveyor

import (
	"fmt"
	"path/filepath"
	"site_parsing/pagetype"
	"site_parsing/waterMarks"
	"strconv"
	"sync"
)

type WaterMarkConsignor struct {
	waterMark *waterMarks.WaterMark
}

func NewWaterMarkConsignor() *WaterMarkConsignor {
	return &WaterMarkConsignor{
		waterMark: waterMarks.NewWaterMark(),
	}
}

// SetWaterMarkPhotoset метод по отправки фотосета для выставления воттермарки
func (w *WaterMarkConsignor) SetPhotosetWaterMark(photoSet *pagetype.PhotoSet) {
	for key, photo := range photoSet.Photos {
		keyString := strconv.Itoa(key + 1)
		wg := new(sync.WaitGroup)
		wg.Add(1)
		address := filepath.Join(pagetype.Dir, photoSet.Id, keyString)
		go w.waterMark.SetPhotoSetWaterMark(wg, photo.Photo, address)
		wg.Wait()
		fmt.Println("Горутины завершили выполнение")
	}
}

// SetWaterMarkPhoto метод по выставлению воттермарки на фото
func (w *WaterMarkConsignor) SetPhotoWaterMark(photoSetId string, photoURL string) pagetype.Photo {
	var time pagetype.Photo
	wg := new(sync.WaitGroup)
	wg.Add(1)
	address := filepath.Join(pagetype.Dir, photoSetId, photoSetId) + ".jpg"
	w.waterMark.SetPhotoWaterMark(wg, photoURL, address, &time)
	wg.Wait()
	return time
}
