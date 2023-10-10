package module

 

import (

    "contactapp/app"

    "contactapp/repository"

)

 

func RegisterModuleRoutes(app *app.App, repository repository.Repository) {

    log := app.Log

    log.Print("============RegisterModuleRoutes.go==============")

    app.WG.Add(4)

    go registerUserRoutes(app, repository)

    go registerContactRoutes(app, repository)

    go registerContactDetailsRoutes(app, repository)

    go registerAuthRoute(app, repository)

 

    app.WG.Wait()

}