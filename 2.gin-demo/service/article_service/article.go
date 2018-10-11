package article_service

import (
	"demo/2.gin-demo/models"
	"demo/2.gin-demo/pkg/export"
	"demo/2.gin-demo/pkg/gredis"
	"demo/2.gin-demo/service/cache_service"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Unknwon/com"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"strconv"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 18:30 
    @File: article.go    
*/

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	if err := models.AddArticle(article); err != nil {
		log.Fatalf("增加文章失败：[%v]", err)
		return err
	}
	return nil
}

func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"modified_by":     a.ModifiedBy,
	})
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			log.Fatalf("从redis中获取文章[%v]失败：[%v]", a.ID, err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			log.Printf("从redis中获取文章[%v]成功", a.ID)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		log.Fatalf("获取指定ID的文章失败：[%v]", err)
		return nil, err
	}
	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]models.Article, error) {
	var (
		articles, cacheArticles []models.Article
	)

	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}

	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if len(data) == 0 || err != nil {
			log.Fatalf("从redis中获取所有文章失败：[%v]", err)
		} else {
			fmt.Printf("从redis中找到了数据..")
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		log.Fatalf("获取文章列表失败：[%v]", err)
		return nil, err
	}
	if len(articles) > 0 {
		gredis.Set(key, articles, 3600)
	}
	return articles, nil
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

//导出
func (a *Article) Export() (string, error) {
	//先取出所有文章
	articles, err := a.GetAll()
	if err != nil {
		log.Fatalf("导出文章列表时出错: [%v]", err)
		return "", err
	}
	log.Println(fmt.Sprintf("articles = %v", articles))
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("文章信息")
	if err != nil {
		log.Fatalf("导出文章时，新建Excel标签失败: [%v]", err)
		return "", err
	}

	row := sheet.AddRow()
	title := []string{"ID", "文章名", "标签", "描述", "内容", "封面", "创建人", "创建时间", "修改人", "修改时间", "状态"}
	var cell *xlsx.Cell
	for _, t := range title {
		cell = row.AddCell()
		cell.Value = t
	}
	//获取每篇文章的信息，然后在Excel中插入一行数据
	for _, a := range articles {
		if a.ID <= 0 {
			continue
		}
		values := []string{
			strconv.Itoa(a.ID),
			a.Title,
			strconv.Itoa(a.Tag.ID),
			a.Desc,
			a.Content,
			a.CoverImageUrl,
			a.CreatedBy,
			strconv.Itoa(a.CreatedOn),
			a.ModifiedBy,
			strconv.Itoa(a.ModifiedOn),
			strconv.Itoa(a.State),
		}
		row = sheet.AddRow()
		for _, v := range values {
			cell = row.AddCell()
			cell.Value = v
		}
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	filename := "article_" + currentTime + ".xlsx"

	fullPath := export.GetExcelFullPath() + filename
	//保存对应的Excel文件到对应的目录
	err = file.Save(fullPath)
	if err != nil {
		log.Fatalf(fmt.Sprintf("保存Excel时失败：[%v]", err))
		return "", err
	}
	return filename, nil
}

//导入
func (a *Article) Import(r io.Reader) error {
	excel, err := excelize.OpenReader(r)
	if err != nil {
		log.Fatalf("导入文章失败: [%v]", err)
		return err
	}

	rows := excel.GetRows("文章信息")
	for index, value := range rows {
		if index == 0 {
			continue
		}
		data := make(map[string]interface{})
		for i, v := range value {
			if rows[0][i] == "标签" {
				data["tag_id"] = com.StrTo(v).MustInt()
			}
			if rows[0][i] == "文章名" {
				data["title"] = v
			}
			if rows[0][i] == "描述" {
				data["desc"] = v
			}
			if rows[0][i] == "内容" {
				data["content"] = v
			}
			if rows[0][i] == "创建人" {
				data["created_by"] = v
			}
			if rows[0][i] == "状态" {
				data["state"] = com.StrTo(v).MustInt()
			}
			if rows[0][i] == "封面" {
				data["cover_image_url"] = v
			}
		}
		fmt.Printf("data = %v\n", data)
		models.AddArticle(data)
	}

	return nil
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	//maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}

	if a.TagID > 0 {
		maps["tag_id"] = a.TagID
	}
	return maps
}
