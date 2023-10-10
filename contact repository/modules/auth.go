package module

import (
	"contactapp/app"
	"contactapp/components/auth/controller"
	"contactapp/components/auth/service"

	"contactapp/repository"
)

func registerAuthRoute(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	authService := service.NewAuthService(appObj.DB, repository)

	authController := controller.NewAuthController(authService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		authController,
	})
}