package module

import (
	"contactapp/app"
	"contactapp/components/user/controller"
	"contactapp/components/user/service"
	"contactapp/repository"
)

func registerUserRoutes(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	userService := service.NewUserService(appObj.DB, repository)

	userController := controller.NewUserController(userService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		userController,
	})
}