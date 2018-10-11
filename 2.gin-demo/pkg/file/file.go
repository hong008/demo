package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 11:52 
    @File: file.go    
*/

//get 文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	return len(content), err
}

//get文件后缀名
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

//检查文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

//check文件权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

//检查文件是否存在，如果不存在则新建目录
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			log.Printf("创建目录时失败： [%v]", err)
			return err
		}
	}
	return nil
}

//新建目录
func MkDir(src string) error {
	return os.Mkdir(src, os.ModePerm)
}

//打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	//file, err := os.Create(name)
	return os.OpenFile(name, flag, perm)
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}
	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to OpenFile: %v", err)
	}
	return f, nil
}
