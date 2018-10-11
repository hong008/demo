package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 14:15 
    @File: article.go    
*/

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
}

//gorm 回调函数
//创建时的时间
func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

//创建修改时的时间
func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}

//文章相关的接口

//查看是否存在指定ID的文章
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

//获取文章数量
func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		log.Fatalf("获取文章数量失败：[%v]", err)
		return 0, err
	}
	return count, nil
}

//获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]Article, error) {
	fmt.Println(fmt.Sprintf("num = %v size = %v maps = %v", pageNum, pageSize, maps))
	var articles []Article
	if pageNum > 0 && pageSize > 0 {
		db = db.Offset(pageNum).Limit(pageSize)
	}
	err := db.Preload("Tag").Where(maps).Find(&articles).Error
	//err := db.Where(maps).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("获取文章列表失败：[%v]", err)
		return nil, err
	}
	fmt.Printf("articles = %+v\n", articles)
	return articles, nil
}

//获取指定ID的文章
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("获取指定ID[%v]的文章失败: [%v]", id, err)
		return nil, err
	}
	return &article, nil
}

//增加文章
func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	}

	if err := db.Create(&article).Error; err != nil {
		log.Fatalf("新增文章失败：[%v]", err)
		return err
	}
	return nil
}

//修改文章
func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ?", id).Update(data).Error; err != nil {
		log.Fatalf("编辑文章[%v]失败：[%v]", id, err)
		return err
	}
	return nil
}

//删除文章
func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(Article{}).Error; err != nil {
		return err
	}
	return nil
}

func CleanAllArticle() bool {
	db.Unscoped().Delete(&Article{})
	return true
}
