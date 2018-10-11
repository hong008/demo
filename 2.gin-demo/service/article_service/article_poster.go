package article_service

import (
	"demo/2.gin-demo/pkg/file"
	"demo/2.gin-demo/pkg/qrcode"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 18:30 
    @File: article_poster.go    
*/

type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.QrCode
}

func NewArticlePoster(posterName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: posterName,
		Article:    article,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) CheckMergedImage(path string) bool {
	return !file.CheckNotExist(path + a.PosterName)
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		log.Fatalf("打开海报文件失败: [%v]", err)
		return nil, err
	}
	return f, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

func (a *ArticlePosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath()
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	if !a.CheckMergedImage(path) {
		mergedF, err := a.OpenMergedImage(path)
		if err != nil {
			fmt.Printf("a.CheckMergedImage fail: %v\n", err)
			return "", "", err
		}
		defer mergedF.Close()

		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			fmt.Printf("打开文件失败: [%v]\n", err)
			return "", "", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			fmt.Printf("file.MustOpen fail: %v\n", err)
			return "", "", err
		}
		defer qrF.Close()

		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			fmt.Printf("jpeg.Decode fail: %v\n", err)
			return "", "", err
		}

		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			fmt.Printf("err = %v\n", err)
			return "", "", err
		}
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		jpeg.Encode(mergedF, jpg, nil)
	}

	return fileName, path, nil
}
