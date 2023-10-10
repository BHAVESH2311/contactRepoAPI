package user

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	FullName string `json:"FullName" gorm:"type:varchar(100)"`
	IsAdmin bool    `json:"IsAdmin" gorm:"type:boolean" default:"false"`
	Email string `json:"Email" gorm:"type:varchar(100)"`
	Password string `json:"Password" type:"varchar(100)"`
}