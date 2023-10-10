package module

import (
	"contactapp/app"
	"contactapp/components/contactDetails/controller"
	"contactapp/components/contactDetails/service"
	"contactapp/repository"
)

func registerContactDetailsRoutes(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	contactDetailService := service.NewContactDetailService(appObj.DB, repository)

	contactDetailController := controller.NewContactDetailController(contactDetailService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactDetailController,
	})
}