package models

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 11:01 
    @File: tag.go    
*/

type Tag struct {
	Model

	Name       string `json:"name"`        //标签名字
	CreatedBy  string `json:"created_by"`  //谁创建的标签
	ModifiedBy string `json:"modified_by"` //谁修改的
	State      int    `json:"state"`       //标签状态
}

//获取标签
func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var tags []Tag
	if pageSize > 0 && pageNum > 0 {
		db = db.Offset(pageNum).Limit(pageSize)
	}

	err := db.Where(maps).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("分页获取标签失败：[%v]", err)
		return nil, err
	}

	return tags, nil
}

//获取标签的count
func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		log.Fatalf("统计标签失败：[%v]", err)
		return 0, err
	}
	return count, nil
}

//增加tag
func AddTag(name string, state int, createBy string) error {
	tag := &Tag{
		Name:      name,
		State:     state,
		CreatedBy: createBy,
	}

	if err := db.Create(tag).Error; err != nil {
		log.Fatalf("创建新标签失败：[%v]", err)
		return err
	}
	return nil
}

//修改标签
func EditTag(id int, data interface{}) error {
	if err := db.Model(&Tag{}).Where("id = ?", id, 0).Update(data).Error; err != nil {
		log.Fatalf("修改指定ID[%v]的标签失败：[%v]", id, err)
		return err
	}
	return nil
}

//删除标签
func DeleteTag(id int) error {
	if err := db.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		log.Fatalf("删除指定ID[%v]的标签失败：[%v]", id, err)
		return err
	}
	return nil
}

//辅助方法
//查看是否已存在指定名字的标签
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("不存在名字为[%]的标签：[%v]", name, err)
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

//查看是否已经存在指定ID的标签
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("不存在指定ID[%v]的标签：[%v]", id, err)
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

//创建tag时增加创建时间
func (t *Tag) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("CreatedOn", time.Now().Unix())
}

//增加tag修改的时间
func (t *Tag) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("ModifiedOn", time.Now().Unix())
}

func CleanAllTag() (bool, error) {
	if err := db.Unscoped().Delete(&Tag{}).Error; err != nil {
		log.Fatalf("清除所有标签失败：[%v]", err)
		return false, err
	}
	return true, nil
}
