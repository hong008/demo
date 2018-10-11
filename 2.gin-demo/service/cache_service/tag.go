package cache_service

import (
	"demo/2.gin-demo/pkg/e"
	"strconv"
	"strings"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 18:01 
    @File: tag.go    
*/

type Tag struct {
	ID    int
	Name  string
	State int

	PageNum  int
	PageSize int
}

func (t *Tag) GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}

	if t.State > 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}

	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}

	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}
	return strings.Join(keys, "_")
}
