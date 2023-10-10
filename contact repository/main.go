package main

import (
	"contactapp/app"
	"contactapp/components/log"
	module "contactapp/modules"
	"contactapp/repository"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	log := log.GetLogger()
	db := app.NewDBConnection(*log)

	if db == nil {
		log.Print("Db connection failed.")
	}
	defer func() {
		db.Close()
		log.Print("Db closed")
	}()
	var wg sync.WaitGroup
	var repo = repository.NewGormRepository()
	app := app.NewApp("Contact App", db, *log,
		&wg, repo)
	app.Init()

	module.RegisterModuleRoutes(app, repo)

	// Need to make sure app starts within 60 seconds of deployment so heroku is able to find port.
	go func() {
		err := app.StartServer()
		if err != nil {
			stopApp(app)
		}
	}()
	module.Configure(app)

	// Stop Server On System Call or Interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	stopApp(app)



	fmt.Print("sub barabar hai")
}
func stopApp(app *app.App) {
	app.Stop()
	app.WG.Wait()
	log.GetLogger().Print("App stopped.")
	os.Exit(0)
}