package contact

import (
	"contactapp/models/user"

	"github.com/jinzhu/gorm"
)

type Contact struct {
	gorm.Model
	User     user.User `gorm:"foreignkey:UserId"` // use UserRefer as foreign key
	UserId   uint      
	FullName string    `json:"FullName" gorm:"type:varchar(100)"`
}