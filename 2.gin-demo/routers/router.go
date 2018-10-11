package routers

import (
	"demo/2.gin-demo/middleware/jwt"
	"demo/2.gin-demo/pkg/export"
	"demo/2.gin-demo/pkg/qrcode"
	"demo/2.gin-demo/pkg/setting"
	"demo/2.gin-demo/pkg/upload"
	"demo/2.gin-demo/routers/api"
	"demo/2.gin-demo/routers/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 01:30 
    @File: router.go    
*/

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	//r.StaticFS("2.gin-demo/runtime/upload/images", http.Dir(upload.GetImageFullPath()))
	//r.StaticFS("2.gin-demo/runtime/export", http.Dir(export.GetExcelFullPath()))
	//r.StaticFS("2.gin-demo/runtime/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.StaticFS("runtime/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("runtime/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("runtime/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.GET("/auth", api.GetAuth)                                          //此处为获取token的方法
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) //此处为swagger路由
	r.POST("/upload", api.UploadImage)                                   //此处为博客封面上传路由
	r.POST("/tags/export", v1.ExportTag)                                 //导出标签
	r.POST("/tags/import", v1.ImportTag)                                 //导入标签
	r.POST("/articles/export", v1.ExportArticles)                        //导出文章
	r.POST("/articles/import", v1.ImportArticles)                        //导入文章

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT()) //此处添加自己实现的权限处理中间件，类似gin.Logger(), gin.Recovery()
	{
		//标签相关的请求
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//文章相关的请求
		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)

		//二维码生成
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)

	}

	return r
}
