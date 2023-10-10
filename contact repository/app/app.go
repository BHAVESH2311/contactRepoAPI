package app

import (
	"contactapp/components/log"
	"contactapp/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	sync.Mutex
	Name string
	Router *mux.Router
	DB *gorm.DB
	Log log.Log

	Server *http.Server
	WG *sync.WaitGroup

	Repository repository.Repository
}

type ModuleConfig interface{
	TableMigration(wg *sync.WaitGroup)
}

type Controller interface{
	RegisterRoutes(*mux.Router)
}

func NewApp(name string, db *gorm.DB,log log.Log, wg *sync.WaitGroup, repository repository.Repository)*App{
	return &App{
		Name: name,
		DB: db,
		Log: log,
		WG : wg,
		Repository: repository,
	}
}

func NewDBConnection(log log.Log) *gorm.DB {
	url := fmt.Sprintf("root:bmbest..@tcp(localhost:3306)/contactapp?charset=utf8mb4&parseTime=true")

	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(3 * time.Minute)

	db.LogMode(true)
	
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci")
	db.BlockGlobalUpdate(true)
	return db
}


func (app *App)MigrateTables(configs []ModuleConfig){
	app.WG.Add(len(configs))
	for _,config := range configs{
		config.TableMigration(app.WG)
		app.WG.Done()
	}
	app.WG.Wait()
	app.Log.Print("End of Migration")
}

func (app *App)Init(){
	app.initializeRouter()
	app.initializeServer()
}

func (app *App) initializeRouter(){
	app.Log.Print("Initialising App Router")
	app.Router=mux.NewRouter().StrictSlash(true)
	app.Router = app.Router.PathPrefix("/api/v1/contactapp").Subrouter()

	app.Router.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode("Server is started")
	})
}

func (app *App) initializeServer() {
	headers:=handlers.AllowedHeaders([]string{
		"Content-Type","X-Total-Count","token",
	})
	methods:=handlers.AllowedMethods([]string{http.MethodPost,http.MethodGet,http.MethodDelete,http.MethodPut,http.MethodOptions})

	origin := handlers.AllowedOriginValidator(app.checkOrigin)

	app.Server=&http.Server{
		Addr: "0.0.0.0:4000",
		ReadTimeout: time.Second*60,
		WriteTimeout: time.Second*60,
		IdleTimeout: time.Second*60,
		Handler: handlers.CORS(headers,methods,origin)(app.Router),
	}
	app.Log.Print("Server started on Port 4000")
}

func (app *App) StartServer()error{
	app.Log.Print("Server Time : ",time.Now())
	app.Log.Print("Server Started on port : 4000" )

	if err:=app.Server.ListenAndServe();err!=nil{
		app.Log.Print(" ListenAndServe error occurred",err)
		return err
	}
	return nil
}

func (app *App) RegisterControllerRoutes(controllers []Controller) {
	app.Lock()
	defer app.Unlock()
	for _, controller := range controllers {
		controller.RegisterRoutes(app.Router.NewRoute().Subrouter())
	}
}

func (app *App) checkOrigin(origin string) bool{
	return true
}

func(app *App) Stop(){
	context,cancel:=context.WithTimeout(context.Background(),100*time.Millisecond)
	defer cancel()
	app.DB.Close()
	app.Log.Print("DB is Closed")

	err:=app.Server.Shutdown(context)
	if err!=nil{
		app.Log.Print("Failed to stop Server")
		return
	}
	app.Log.Print("Server Stopped Successfully")
}


