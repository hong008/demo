package api

import (
	"demo/2.gin-demo/pkg/e"
	"demo/2.gin-demo/pkg/logging"
	"demo/2.gin-demo/pkg/upload"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 15:29 
    @File: upload.go    
*/

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image") //获取上传的图片
	if err != nil {
		logging.Warn(fmt.Sprintf("上传图片失败：[%v]", err))
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
		return
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(fmt.Sprintf("校验图片失败：%v", err))
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.Warn(fmt.Sprintf("保存图片失败: %v", err))
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
