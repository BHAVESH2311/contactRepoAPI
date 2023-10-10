package module

import (
	"contactapp/app"
	"contactapp/components/contact/controller"
	"contactapp/components/contact/service"
	"contactapp/repository"
)

func registerContactRoutes(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	contactService := service.NewContactService(appObj.DB, repository)

	contactController := controller.NewContactController(contactService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactController,
	})
}