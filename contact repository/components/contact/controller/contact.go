package controller

import (
	"contactapp/auth"
	"contactapp/components/contact/service"
	"contactapp/components/log"
	"contactapp/errors"
	"contactapp/models/contact"
	"contactapp/web"
	"contactapp/models/user"
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)

type ContactController struct {
	log     log.Log
	service *service.ContactService
}

func NewContactController(contactService *service.ContactService,
	log log.Log) *ContactController {
	return &ContactController{
		service: contactService,
		log:     log,
	}
}
func (controller *ContactController) RegisterRoutes(router *mux.Router) {
	contactRouter := router.PathPrefix("/user/{userId}/contact").Subrouter()
	contactRouter.HandleFunc("/", auth.Protect(controller.CreateContact)).Methods(http.MethodPost)
	contactRouter.HandleFunc("/", auth.Protect(controller.GetAllContacts)).Methods(http.MethodGet)
	contactRouter.HandleFunc("/{id}", auth.Protect(controller.GetContact)).Methods(http.MethodGet)
	contactRouter.HandleFunc("/{id}", auth.Protect(controller.UpdateContact)).Methods(http.MethodPut)
	contactRouter.HandleFunc("/{id}", auth.Protect(controller.DeleteContact)).Methods(http.MethodDelete)
}


func (c *ContactController)CreateContact(w http.ResponseWriter, r *http.Request){
	newContact:=contact.Contact{}
	vars:=mux.Vars(r)
	err:=web.UnmarshalJSON(r,&newContact)
	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusBadRequest))
		return
	}
	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,errors.NewHTTPError(err.Error(),http.StatusBadRequest))
		return
	}

	userId,err:=strconv.Atoi(vars["userId"])

	if err != nil {
		c.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	newContact.UserId=uint(userId)

	contactOwner := user.User{}

	contactOwner.ID = uint(userId)

	userClaim:=r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID!=uint64(contactOwner.ID) && !userClaim.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized Access")
		return
	}
	err=c.service.CreateContact(&newContact)
	if err!=nil{
		c.log.Print(err)
		web.RespondErrorMessage(w,http.StatusBadRequest,"Cannot create Contact")
		return
	}

	web.RespondJSON(w,http.StatusOK,"Contact Created Successfully")

}

func (c *ContactController)GetAllContacts(w http.ResponseWriter, r *http.Request){

	limit,offset,err:=web.LimitOffsetExtract(r)

	vars:=mux.Vars(r)

	if err!=nil{
		c.log.Print(err)
		web.RespondErrorMessage(w,http.StatusBadRequest,"cannot get all contacts")
		return
	}


	userId,err:=strconv.Atoi(vars["userId"])

	if err!=nil{
		c.log.Print(err)
		return
	}

	userClaim :=r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID!=uint64(userId) && !userClaim.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized access")
		return
	}


	contacts,total,err:=c.service.GetAllContacts(uint(userId),limit,offset)

	if err!=nil{
		c.log.Print(err)
		web.RespondErrorMessage(w,http.StatusUnauthorized,"unauthorized access")
		return
	}

	web.RespondJSONWithXTotalCount(w,http.StatusOK,total,contacts)


}


func (c *ContactController) GetContact(w http.ResponseWriter, r *http.Request){

	vars:=mux.Vars(r)
	contactId,err:=strconv.Atoi(vars["id"])

	if err!=nil{
		c.log.Print(err)
		web.RespondErrorMessage(w,http.StatusBadRequest,"Invalid Contact Id")
		return
	}


	userId,err:=strconv.Atoi(vars["userId"])


	if err!=nil{
		c.log.Print(err)
		web.RespondErrorMessage(w,http.StatusBadRequest,"Invalid Contact Id")
		return
	}

	userClaim:=r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID !=uint64(userId) && !userClaim.IsAdmin{

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized access")
		return
	}

	includes := web.ParseIncludes(r)

	contact ,err:=c.service.GetContact(uint(userId),contactId,includes)

	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,err)
		return
	}
	web.RespondJSON(w,http.StatusOK,contact)
}


func (controller *ContactController) UpdateContact(w http.ResponseWriter, r *http.Request) {
	contactToUpdate := contact.Contact{}

	vars:=mux.Vars(r)

	// Unmarshal JSON.
	err := web.UnmarshalJSON(r, &contactToUpdate)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactId,err:=strconv.Atoi(vars["id"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactToUpdate.ID = uint(contactId)

	userId,err:=strconv.Atoi(vars["userId"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactToUpdate.UserId = uint(userId)

	userClaim := r.Context().Value(auth.UserIDKey).(*auth.Claims)
	if userClaim.ID != uint64(contactToUpdate.UserId) && !userClaim.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not Authorized to access")
		return
	}

	err = controller.service.UpdateContact(&contactToUpdate)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, contactToUpdate)
}


func (c *ContactController) DeleteContact(w http.ResponseWriter, r *http.Request){

	vars:=mux.Vars(r)
	contactToDelete := contact.Contact{}

	contactId,err:=strconv.Atoi(vars["id"])
	if err != nil {
		c.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactToDelete.ID=uint(contactId)

	userId,err:=strconv.Atoi(vars["userId"])

	if err != nil {
		c.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactToDelete.UserId = uint(userId)

	userClaim:=r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID!=uint64(contactToDelete.UserId) && !userClaim.IsAdmin{
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Unauthorized Access")
		return
	}

	err=c.service.DeleteContact(&contactToDelete)

	if err!=nil{
		c.log.Print(err)
		web.RespondError(w,err)
		return
	}

	web.RespondJSON(w,http.StatusOK,"Deleted User Successfully")

}










