package controller

import (
	"contactapp/components/log"
	"contactapp/components/user/service"
	"contactapp/errors"
	"contactapp/auth"
	"contactapp/models/user"
	"contactapp/utils"
	"contactapp/web"
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)


type UserController struct {
	log log.Log
	service *service.UserService
}


func NewUserController(userService *service.UserService, log log.Log) *UserController {
	return &UserController{
		service: userService,
		log:     log,
	}
}

func (controller *UserController) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("/admin", controller.RegisterAdmin).Methods(http.MethodPost)

	userRouter.HandleFunc("/", auth.Protect(auth.IsAdmin(controller.GetAllUsers))).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", auth.Protect(controller.GetUser)).Methods(http.MethodGet)
	userRouter.HandleFunc("/", auth.Protect(auth.IsAdmin(controller.RegisterUser))).Methods(http.MethodPost)
	userRouter.HandleFunc("/{id}", auth.Protect(controller.UpdateUser)).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}", auth.Protect(controller.DeleteUser)).Methods(http.MethodDelete)

}


func (controller *UserController)RegisterAdmin(w http.ResponseWriter, r *http.Request){

	newAdmin:=user.User{}
	err:=web.UnmarshalJSON(r,&newAdmin)
	if(err!=nil){
		controller.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusBadRequest))
		return
	}

	newAdmin.IsAdmin=true
	newAdmin.Password,err = utils.HashPassword(newAdmin.Password)
	if err!=nil{
		controller.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusBadRequest))
		return
	}

	err=utils.UserValidator(newAdmin)
	if err!=nil{
		controller.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusForbidden))
		return
	}
	err=controller.service.CreateAdmin(&newAdmin)
	if err!=nil{
		controller.log.Print(err.Error())
		web.RespondError(w,err)
		return
	}

	web.RespondJSON(w,http.StatusCreated,newAdmin)


}


func(controller *UserController) RegisterUser(w http.ResponseWriter, r *http.Request){
	newUser := user.User{}

	err:=web.UnmarshalJSON(r,&newUser)
	if(err!=nil){
		controller.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusBadRequest))
		return
	}

	err=utils.UserValidator(newUser)
	if err!=nil{
		controller.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusForbidden))
		return
	}

	newUser.Password,err=utils.HashPassword(newUser.Password)
	if err!=nil{
		controller.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusBadRequest))
		return
	}
	err = controller.service.CreateUser(&newUser)

	if err!=nil{
		controller.log.Print(err.Error())
		web.RespondError(w,err)
		return
	}

	web.RespondJSON(w,http.StatusCreated,newUser)
}


func(c *UserController) GetUser(w http.ResponseWriter, r *http.Request){
	vars:=mux.Vars(r)
	uId,err:=strconv.Atoi(vars["id"])
	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusBadRequest))
		return
	}
	currentUser:=r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if currentUser.ID != uint64(uId) && !currentUser.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("user unauthorized")
		return
	}

	user,err := c.service.GetUserById(uint(uId))
	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,err)
		return
	}

	web.RespondJSON(w,http.StatusOK,user)


}


func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := web.LimitOffsetExtract(r)

	if err != nil {
		c.log.Print(err)
		web.RespondErrorMessage(w,http.StatusBadRequest,"cannot get all users")
		
		return
	}

	users, total, err := c.service.GetAllUsers(limit, offset)

	if err != nil {
		c.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSONWithXTotalCount(w, http.StatusOK, total, users)
}


func(c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request){

	vars:=mux.Vars(r)

	user:=user.User{}
	err:=web.UnmarshalJSON(r,&user)
	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,err)
		return
	}

	uId,err:=strconv.Atoi(vars["id"])

	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,err)
		return
	}

	user.ID=uint(uId)

	currentUser:=r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if currentUser.ID !=uint64(uId) && !currentUser.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized user")
		return
	}

	err=c.service.UpdateUser(&user)

	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,err)
		return
	}

	web.RespondJSON(w,http.StatusOK,user)

}


func (c *UserController) DeleteUser(w http.ResponseWriter,r *http.Request){
	vars:=mux.Vars(r)
	user:=user.User{}

	userId,err := strconv.Atoi(vars["id"])

	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,err)
		return
	}

	user.ID=uint(userId)
	currentUser:=r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if currentUser.ID!=uint64(userId) && !currentUser.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized Access")
		return
	}

	err=c.service.DeleteUser(&user)
	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,errors.NewHTTPError("User not Deleted",http.StatusBadRequest))
		return
	}

	web.RespondJSON(w,http.StatusOK,"User Deleted Successfully")

}


