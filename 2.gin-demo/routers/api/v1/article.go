package v1

import (
	"demo/2.gin-demo/pkg/app"
	"demo/2.gin-demo/pkg/e"
	"demo/2.gin-demo/pkg/export"
	"demo/2.gin-demo/pkg/qrcode"
	"demo/2.gin-demo/pkg/setting"
	"demo/2.gin-demo/pkg/util"
	"demo/2.gin-demo/service/article_service"
	"demo/2.gin-demo/service/tag_service"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 13:59 
    @File: article.go    
*/

// @Summary 获取单个文章
// @Produce  json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":{"id":3,"created_on":1516937037,"modified_on":0,"tag_id":11,"tag":{"id":11,"created_on":1516851591,"modified_on":0,"name":"312321","created_by":"4555","modified_by":"","state":1},"content":"5555","created_by":"2412","modified_by":"","state":1},"msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	log.Printf("开始获取指定ID的文章...\n")
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	log.Println(fmt.Sprintf("id = %v", id))

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		log.Fatal(fmt.Sprintf("获取文章[%v]失败: [%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exists {
		log.Fatal(fmt.Sprintf("不存在指定ID[%v]的文章", id))
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		log.Fatal(fmt.Sprintf("从redis中获取指定ID[%v]的文章失败：[%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	log.Println(fmt.Sprintf("获取文章[%v]成功", id))
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// @Summary 获取多个文章
// @Produce  json
// @Param tag_id query int false "TagID"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":[{"id":3,"created_on":1516937037,"modified_on":0,"tag_id":11,"tag":{"id":11,"created_on":1516851591,"modified_on":0,"name":"312321","created_by":"4555","modified_by":"","state":1},"content":"5555","created_by":"2412","modified_by":"","state":1}],"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	tagId := -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(tagId).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		log.Fatal(fmt.Sprintf("统计包含标签[%v]的文章失败：[%v]", tagId, err))
		appG.Response(http.StatusOK, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		log.Fatal(fmt.Sprintf("从redis中获取所有文章失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary 新增文章
// @Produce  json
// @Param tag_id query int true "TagID"
// @Param title query string true "Title"
// @Param desc query string true "Desc"
// @Param content query string true "Content"
// @Param created_by query string true "CreatedBy"
// @Param cover_image_url query string true "CoverImageUrl"
// @Param state query int true "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	fmt.Println("增加文章...")
	log.Println("开始增加文章...")
	appG := app.Gin{
		C: c,
	}

	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	coverImageUrl := c.Query("cover_image_url")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Required(coverImageUrl, "cover_image_url").Message("封面不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: tagId}
	exists, err := tagService.ExistByID()
	if err != nil {
		log.Fatal(fmt.Sprintf("查看是否存在指定ID[%v]的标签时失败：[%v]", tagId, err))
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		log.Fatal(fmt.Sprintf("不存在指定ID[%v]的标签", tagId))
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         tagId,
		Title:         title,
		Desc:          desc,
		Content:       content,
		CoverImageUrl: coverImageUrl,
		State:         state,
		CreatedBy:     createdBy,
	}

	if err := articleService.Add(); err != nil {
		log.Fatal(fmt.Sprintf(""))
		appG.Response(http.StatusOK, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	log.Println(fmt.Sprintf("增加文章成功"))
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 修改文章
// @Produce  json
// @Param id param int true "ID"
// @Param tag_id query string false "TagID"
// @Param title query string false "Title"
// @Param desc query string false "Desc"
// @Param content query string false "Content"
// @Param modified_by query string true "ModifiedBy"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	log.Println("开始修改文章...")
	appG := app.Gin{
		C: c,
	}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	coverImageUrl := c.PostForm("cover_image_url")
	modifiedBy := c.PostForm("modified_by")

	var state = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")

	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.Required(title, "title").Message("标题不能为空")

	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.Required(title, "desc").Message("描述不能为空")

	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.Required(coverImageUrl, "cover_image_url").Message("封面不能为空")
	valid.MaxSize(coverImageUrl, 255, "cover_image_url").Message("封面最长为255字符")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		ID:            id,
		Title:         title,
		TagID:         tagId,
		Desc:          desc,
		Content:       content,
		CoverImageUrl: coverImageUrl,
		ModifiedBy:    modifiedBy,
	}

	exists, err := articleService.ExistByID()
	if err != nil {
		log.Fatal(fmt.Sprintf("查看是否存在指定ID[%v]的文章时失败：[%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exists {
		log.Fatal(fmt.Sprintf("不存在指定ID[%v]的文章", id))
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: tagId}
	exists, err = tagService.ExistByID()
	if err != nil {
		log.Fatal(fmt.Sprintf("判断数据库中是否存在指定ID[%v]的标签失败：[%v]", tagId, err))
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		log.Fatal(fmt.Sprintf("数据库中不存在指定ID[%v]的标签", tagId))
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		log.Fatal(fmt.Sprintf("修改标签失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	log.Println(fmt.Sprintf("修改文章[%v]成功", id))
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除文章
// @Produce  json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}

	exists, err := articleService.ExistByID()
	if err != nil {
		log.Fatal(fmt.Sprintf("判断数据库中是否存在指定ID[%v]的文章时失败：[%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		log.Fatal(fmt.Sprintf("数据库中不存在指定ID[%v]文章", id))
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		log.Fatal(fmt.Sprintf("删除指定ID[%v]的文章失败：[%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 导出文章
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /articles/export [post]
func ExportArticles(c *gin.Context) {
	appG := app.Gin{C: c}

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	articleService := article_service.Article{
		State: state,
	}

	fileName, err := articleService.Export()
	if err != nil {
		log.Fatal(fmt.Sprintf("导出文章时失败: [%v]", err))
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	log.Println(fmt.Sprintf("导出文章成功"))
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"export_url":      export.GetExcelFullUrl(fileName),
		"export_save_url": export.GetExcelPath() + fileName,
	})
}

// @Summary 导入文章
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /articles/import [post]
func ImportArticles(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatalf(fmt.Sprintf("导入文章时读取文件失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	articleService := article_service.Article{}
	err = articleService.Import(file)
	if err != nil {
		log.Fatal(fmt.Sprintf("导入文章时失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	log.Println("导入文章成功")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func GenerateArticlePoster(c *gin.Context) {
	appG := app.Gin{C: c}

	article := &article_service.Article{}
	qr := qrcode.NewQrCode(e.QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)
	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url": qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url":filePath+posterName,
	})
}
