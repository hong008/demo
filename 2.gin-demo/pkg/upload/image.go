package upload

import (
	"demo/2.gin-demo/pkg/file"
	"demo/2.gin-demo/pkg/logging"
	"demo/2.gin-demo/pkg/setting"
	"demo/2.gin-demo/pkg/util"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 15:11 
    @File: image.go    
*/

func GetImageFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

//获取图片文件名，并返回MD5加密后的文件名
func GetImageName(name string) string {
	ext := file.GetExt(name) //获取图片文件扩展名
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

//获取图片文件在项目中的相对路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

//获取图片文件的完整路径，
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

//校验图片格式是否符合要求
func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			logging.Info(fmt.Sprintf("校验图片格式成功: ext = [%v]", ext))
			return true
		}
	}
	return false
}

//校验封面文件大小
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		logging.Error(fmt.Sprintf("校验博客封面大小时出错: [%v]", err))
		return false
	}
	isOk := size <= setting.AppSetting.ImageMaxSize
	logging.Info(fmt.Sprintf("校验图片大小时：size = [%v]  ImageMaxSize = [%v]", size, setting.AppSetting.ImageMaxSize))
	return isOk
}

//判断是否存在文件
func CheckImage(src string) error {
	logging.Info(fmt.Sprintf("in CheckImage...src = [%v]", src))
	dir, err := os.Getwd()
	if err != nil {
		logging.Error(fmt.Sprintf("校验封面时出错：%v", err))
		return fmt.Errorf("os.Getwd err: %v", err)
	}
	logging.Info("in CheckImage...dir = [%v]", dir)
	logging.Info(fmt.Sprintf("%v", dir+"/"+src))
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		logging.Error(fmt.Sprintf("检查图片目录时失败：%v", err))
		return fmt.Errorf("file.IsNotExistMdDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm {
		logging.Error(fmt.Sprintf("校验图片权限失败"))
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	return nil
}
