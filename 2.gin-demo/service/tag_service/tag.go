package tag_service

import (
	"demo/2.gin-demo/models"
	"demo/2.gin-demo/pkg/export"
	"demo/2.gin-demo/pkg/gredis"
	"demo/2.gin-demo/service/cache_service"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"strconv"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/27 16:01 
    @File: tag.go    
*/

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}
	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)
	cache := cache_service.Tag{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			log.Fatalf("从redis获取tag失败：[%v]", err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		log.Fatalf("从数据库中获取标签失败：[%v]", err)
		return nil, err
	}
	if len(tags) > 0 {
		gredis.Set(key, tags, 3600)
	}
	return tags, nil

}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		log.Fatalf("导出标签出错：[%v]", err)
		return "", err
	}
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		log.Fatalf("导出标签到文件时出错：[%v]", err)
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	//currentTime := strconv.Itoa(int(time.Now().Unix()))
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	filename := "tag_" + currentTime + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	//保存对应的Excel文件到对应的目录
	err = file.Save(fullPath)
	if err != nil {
		log.Fatalf(fmt.Sprintf("保存Excel时失败：[%v]", err))
		return "", err
	}
	return filename, nil
}

//导入标签
func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		log.Fatal(fmt.Sprintf("导入时，初始化reader失败：[%v]", err))
		return err
	}
	rows := xlsx.GetRows("标签信息")
	fmt.Printf("rows = %v\n", rows)
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}

			models.AddTag(data[1], 1, data[2])
		}
	}
	return nil
}

//
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	//maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}
