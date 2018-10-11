package qrcode

import (
	"demo/2.gin-demo/pkg/file"
	"demo/2.gin-demo/pkg/setting"
	"demo/2.gin-demo/pkg/util"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
	"log"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/29 10:23 
    @File: qrcode.go    
*/

type QrCode struct {
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = ".jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Ext:    EXT_JPG,
		Level:  level,
		Mode:   mode,
	}
}

func GetQrCodePath() string {
	return setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetQrCodePath() + name
}

func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	if file.CheckNotExist(src) {
		return false
	}
	return true
}

func (q *QrCode) Encode(path string) (string, string, error) {
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	src := path + name

	if file.CheckNotExist(src) {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			log.Fatalf("生成二维码->encode时失败: [%v]", err)
			return "", "", err
		}
		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			log.Fatalf("生成二维码->scale时失败: [%v]", err)
			return "", "", err
		}
		f, err := file.MustOpen(name, path)
		if err != nil {
			log.Fatalf("打开文件失败：[%v]", err)
			return "", "", err
		}
		defer f.Close()

		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}
	return name, path, nil
}
