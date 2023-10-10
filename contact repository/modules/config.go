package module

import (
	"contactapp/app"
	"contactapp/models/contact"
	contactDetail "contactapp/models/contactDetails"
	"contactapp/models/user"
)

func Configure(appObj *app.App) {
	userModule := user.NewUserModuleConfig(appObj.DB)
	contactModule := contact.NewContactModuleConfig(appObj.DB)
	contactInfoModule := contactDetail.NewContactDetailModuleConfig(appObj.DB)

	appObj.MigrateTables([]app.ModuleConfig{userModule, contactModule, contactInfoModule})
}