package export

import "demo/2.gin-demo/pkg/setting"

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/27 16:49 
    @File: excel.go    
*/

func GetExcelFullUrl(name string) string {

	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}
