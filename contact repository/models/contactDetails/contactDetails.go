package contactdetails

import (
	"contactapp/models/contact"

	"github.com/jinzhu/gorm"
)

type ContactDetails struct {
	gorm.Model
	Contact      contact.Contact `gorm:"foreignkey:ContactRefer"`
	ContactRefer uint
	UserRefer   uint
	Type         string `json:"Type" gorm:"type:varchar(100)"`
	Value        string `json:"Value" gorm:"type:varchar(100)"`
}