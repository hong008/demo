package v1

import (
	"demo/2.gin-demo/pkg/app"
	"demo/2.gin-demo/pkg/e"
	"demo/2.gin-demo/pkg/export"
	log "demo/2.gin-demo/pkg/logging"
	"demo/2.gin-demo/pkg/setting"
	"demo/2.gin-demo/pkg/util"
	"demo/2.gin-demo/service/tag_service"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 10:50 
    @File: tag.go
*/

// @Summary 获取多个文章标签
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{"lists":[{"id":3,"created_on":1516849721,"modified_on":0,"name":"3333","created_by":"4555","modified_by":"","state":0}],"total":29},"msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageSize: setting.AppSetting.PageSize,
		PageNum:  util.GetPage(c),
	}

	tags, err := tagService.GetAll()
	if err != nil {
		log.Error(fmt.Sprintf("获取所有标签时失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	count, err := tagService.Count()
	if err != nil {
		log.Error(fmt.Sprintf("统计标签时失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	log.Info("获取所有标签成功")
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

// @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	appG := app.Gin{C: c}

	name := c.PostForm("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createBy := c.PostForm("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      name,
		CreatedBy: createBy,
		State:     state,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		log.Error(fmt.Sprintf("判断是否存在标签[%v]失败: [%v]", name, err))
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if exists {
		log.Error(fmt.Sprintf("不存在标签[%v]", name))
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		log.Error(fmt.Sprintf("增加标签[%v]失败：[%v]", name, err))
		appG.Response(http.StatusOK, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	log.Info(fmt.Sprintf("增加标签[%v]成功", name))
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 修改文章标签
// @Produce  json
// @Param id param int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	name := c.PostForm("name")
	modifiedBy := c.Query("modified_by")

	var state = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         id,
		State:      state,
		Name:       name,
		ModifiedBy: modifiedBy,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		log.Error(fmt.Sprintf("判断是否存在标签[%v]失败：[%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		log.Warn(fmt.Sprintf("修改标签失败， 不存在标签[%v]", id))
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		log.Error(fmt.Sprintf("修改标签失败：[%v] [%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	log.Info(fmt.Sprintf("修改标签[%v]成功", id))
}

// @Summary 删除文章标签
// @Produce  json
// @Param id param int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID: id,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		log.Error(fmt.Sprintf("删除标签时，根据ID[%v]判断是否存在失败：[%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		log.Warn(fmt.Sprintf("删除标签时， 不存在标签[%v]", id))
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Delete()
	if err != nil {
		log.Error(fmt.Sprintf("删除标签[%v]失败：[%v]", id, err))
		appG.Response(http.StatusOK, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	log.Info(fmt.Sprintf("删除标签[%v]成功", id))
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary 导出标签
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /tags/export [post]
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	fmt.Printf("name = %v  state = %v \n", name, state)
	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}
	filename, err := tagService.Export()
	if err != nil {
		log.Error(fmt.Sprintf("导出标签时失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	log.Info("导出标签成功")
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

// @Summary 导入标签
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Router /tags/import [post]
func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Error(fmt.Sprintf("导入标签时读取文件失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		log.Error(fmt.Sprintf("导入标签时失败：[%v]", err))
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	log.Info("导入标签成功")
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
