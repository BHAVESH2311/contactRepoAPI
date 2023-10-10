package controller

import (
	auth "contactapp/auth"
	"contactapp/components/contactDetails/service"
	"contactapp/components/log"
	"contactapp/errors"
	contactDetails "contactapp/models/contactDetails"
	"contactapp/web"
	"encoding/json"
	"net/http"
	"contactapp/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

type ContactDetailController struct {
	log     log.Log
	service *service.ContactDetailService
}

func NewContactDetailController(contactService *service.ContactDetailService,
	log log.Log) *ContactDetailController {
	return &ContactDetailController{
		service: contactService,
		log:     log,
	}
}

func (controller *ContactDetailController) RegisterRoutes(router *mux.Router) {
	contactInfoRouter := router.PathPrefix("/user/{userId}/contact/{contactId}/contactinfo").Subrouter()

	contactInfoRouter.HandleFunc("/", auth.Protect(controller.CreateContactDetails)).Methods(http.MethodPost)
	contactInfoRouter.HandleFunc("/", auth.Protect(controller.GetAllContactDetails)).Methods(http.MethodGet)
	contactInfoRouter.HandleFunc("/{id}", auth.Protect(controller.GetContactDetails)).Methods(http.MethodGet)
	contactInfoRouter.HandleFunc("/{id}", auth.Protect(controller.UpdateContactDetails)).Methods(http.MethodPut)
	contactInfoRouter.HandleFunc("/{id}", auth.Protect(controller.DeleteContactDetails)).Methods(http.MethodDelete)

}

func (controller *ContactDetailController) CreateContactDetails(w http.ResponseWriter, r *http.Request) {
	newContactDetail := contactDetails.ContactDetails{}

	vars:=mux.Vars(r)

	err := web.UnmarshalJSON(r, &newContactDetail)

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	userId,err:=strconv.Atoi(vars["userId"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	middleware.CurrentUserCheker(w,r,userId)

	// userClaim := r.Context().Value(auth.UserIDKey).(*auth.Claims)

	// if userClaim.ID != uint64(userId) && !userClaim.IsAdmin {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode("Not Authorized to access")
	// 	return
	// }


	newContactDetail.UserRefer = uint(userId)

	

	contactId, err := strconv.Atoi(vars["contactId"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	newContactDetail.ContactRefer= uint(contactId)

	err = controller.service.CreateContactDetail(&newContactDetail)

	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	var response = map[string]interface{}{
		"message":     "Contact Info Created Successfully",
		"contactInfo": newContactDetail,
	}

	web.RespondJSON(w, http.StatusOK, response)
}

func (controller *ContactDetailController) GetAllContactDetails(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := web.LimitOffsetExtract(r)
	vars:=mux.Vars(r)

	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	userId,err:=strconv.Atoi(vars["userId"])


	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	userClaim := r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID != uint64(userId) && !userClaim.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not Authorized to Access")
		return
	}

	

	contactId, err := strconv.Atoi(vars["contactId"])

	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	includes := web.ParseIncludes(r)

	contacts, total, err := controller.service.GetAllContactDetails(uint(userId), uint(contactId), includes, limit, offset)

	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSONWithXTotalCount(w, http.StatusOK, total, contacts)
}

func (controller *ContactDetailController) GetContactDetails(w http.ResponseWriter, r *http.Request) {

	vars:=mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactId, err := strconv.Atoi(vars["contactId"])

	userClaim := r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID != uint64(userId) && !userClaim.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not Authorized to Access")
		return
	}

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactDetailId, err := strconv.Atoi(vars["id"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	includes := web.ParseIncludes(r)

	contactInfo, err := controller.service.GetContactDetails(uint(userId), uint(contactId), uint(contactDetailId), includes)

	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, contactInfo)
}

func (controller *ContactDetailController) UpdateContactDetails(w http.ResponseWriter, r *http.Request) {
	cdToUpdate := contactDetails.ContactDetails{}

	vars:=mux.Vars(r)

	err := web.UnmarshalJSON(r, &cdToUpdate)

	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	contactDetailId, err := strconv.Atoi(vars["id"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	cdToUpdate.ID = uint(contactDetailId)

	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	userClaim := r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID != uint64(userId )&& !userClaim.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not Authorized to Access")
		return
	}

	cdToUpdate.UserRefer = uint(userId)

	contactId, err := strconv.Atoi(vars["contactId"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	cdToUpdate.ContactRefer = uint(contactId)

	err = controller.service.UpdateContactDetails(&cdToUpdate)

	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	var response = map[string]interface{}{
		"message":     "Contact Info Updated Successfully",
		"contactInfo": cdToUpdate,
	}

	web.RespondJSON(w, http.StatusOK, response)
}

func (controller *ContactDetailController) DeleteContactDetails(w http.ResponseWriter, r *http.Request) {

	cdToDelete := contactDetails.ContactDetails{}
	vars:=mux.Vars(r)

	contactDetailId, err := strconv.Atoi(vars["id"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	cdToDelete.ID = uint(contactDetailId)

	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	cdToDelete.UserRefer = uint(userId)

	userClaim := r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID != uint64(userId) && !userClaim.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not Authorized to Access")
		return
	}

	contactId, err := strconv.Atoi(vars["contactId"])

	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}

	cdToDelete.ContactRefer = uint(contactId)

	err = controller.service.DeleteContactDetails(&cdToDelete)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Contact Info Deleted Successfully")
}