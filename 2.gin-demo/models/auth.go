package models

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/18 16:28 
    @File: auth.go    
*/

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	db.Select("id").Where(Auth{Username: username, Password: password}).Find(&auth)
	if auth.ID > 0 {
		return true
	}

	return false
}

