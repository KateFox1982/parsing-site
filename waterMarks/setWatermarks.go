package waterMarks

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"site_parsing/pagetype"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WaterMark struct {
}

func NewWaterMark() *WaterMark {
	return &WaterMark{}
}

func (w *WaterMark) SetPhotoWaterMark(wg *sync.WaitGroup, link, address string, timePhoto *pagetype.Photo) {
	img, wmb, err := openSiteAndWatemark(link)
	if err != nil {
		log.Fatalf("Ошибка открытия страницы  или вотермарки в openSiteAndWatemark", err)
	}
	watermark, err := png.Decode(wmb)
	defer wmb.Close()
	if err != nil {
		log.Fatalf("Ошибка при декорировании логотипа в методе SetWatermark", err)
	}
	offset := image.Pt(236, 328)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
	timePhoto.TimeDown = time.Now()
	imgw, err := os.Create(address)
	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}
	timePhoto.TimeSave = time.Now()
	err = jpeg.Encode(imgw, m, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {

	}
	defer imgw.Close()
	defer wg.Done()
	return
}

func (w *WaterMark) SetPhotoSetWaterMark(wg *sync.WaitGroup, link, address string) {
	img, wmb, err := openSiteAndWatemark(link)
	if err != nil {
		log.Fatalf("Ошибка открытия страницы  или вотермарки в openSiteAndWatemark", err)
	}
	bgDimensions := img.Bounds().Max
	markFit := resizeImage(pagetype.Watermarks, "3000x3000")
	markDimensions := markFit.Bounds().Max
	bgAspectRatio := math.Round(float64(bgDimensions.X) / float64(bgDimensions.Y))
	xPos, yPos := calcWaterMarkPosition(bgDimensions, markDimensions, bgAspectRatio)
	watermark, err := png.Decode(wmb)
	defer wmb.Close()
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	offset := image.Pt(xPos, yPos)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
	imgw, err := os.Create(address)
	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}
	err = jpeg.Encode(imgw, m, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		log.Fatalf("Ошибка Encode: %v", err)
	}
	defer imgw.Close()
	defer wg.Done()
}

func openSiteAndWatemark(link string) (image.Image, *os.File, error) {
	resp, err := http.Get(link)
	if err != nil {
		log.Fatalf("Ошибка открытия страницы в openSiteAndWatemark", err)
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Ошибка открытия страницы в openSiteAndWatemark", err)
		return nil, nil, err
	}
	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		err = fmt.Errorf("Ошибка открытия страницы в openSiteAndWatemark", err)
		return img, nil, err
	}
	wmb, err := os.Open(pagetype.Watermarks)
	if err != nil {
		err = fmt.Errorf("Ошибка открытия вотермарки в openSiteAndWatemark%v", err)
		return img, nil, err
	}
	return img, wmb, err
}

//приватная функция по изменению размера логотипа
func resizeImage(image, dimensions string) image.Image {
	width, height := parseCoordinates(dimensions, "x")
	src := openImage(image)
	return imaging.Fit(src, width, height, imaging.Lanczos)
}

//parseCoordinates приватная функция преобразование координат логотипа из строкового значения в int
func parseCoordinates(input, delimiter string) (int, int) {
	arr := strings.Split(input, delimiter)
	// convert a string to an int
	x, err := strconv.Atoi(arr[0])
	if err != nil {
		log.Fatalf("failed to parse x coordinate: %v", err)
	}
	y, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Fatalf("failed to parse y coordinate: %v", err)
	}
	return x, y
}

//openImage приватная функция открытия watermark
func openImage(name string) image.Image {
	src, err := imaging.Open(name)
	if err != nil {
		log.Fatalf("Ошибка отурытия воттермарки", err)
	}
	return src
}

//calcWaterMarkPosition приватная функция по вычислению координат расположения логотипа с отступом 20 пикселей
func calcWaterMarkPosition(bgDimensions, markDimensions image.Point, aspectRatio float64) (int, int) {
	bgX := bgDimensions.X
	bgY := bgDimensions.Y
	markX := markDimensions.X
	markY := markDimensions.Y
	padding := 20 * int(aspectRatio)
	return bgX - markX - padding, bgY - markY - padding
}
