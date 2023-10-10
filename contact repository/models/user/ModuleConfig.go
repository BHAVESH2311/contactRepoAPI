package user

import (
	"contactapp/components/log"
	"sync"

	"github.com/jinzhu/gorm"
)

type ModuleConfig struct {
	DB *gorm.DB
}

func NewUserModuleConfig(db *gorm.DB) *ModuleConfig {

	return &ModuleConfig{
		DB:db,
	}
}

func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {
	var models []interface{}=[]interface{}{
		&User{},
	}

	for _,model:=range models{
		if err:=config.DB.AutoMigrate(model).Error;err!=nil{
			log.GetLogger().Print("Auto Migration ==> %s",err)
		}
	}

	log.GetLogger().Print("Test Module Configured")
}
