package controller

import (
	"contactapp/components/auth/service"
	"contactapp/components/log"
	 "contactapp/auth"
	"contactapp/models/user"
	"contactapp/web"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// ContactController gives access to CRUD operations for entity
type AuthController struct {
	log     log.Log
	service *service.AuthService
}

func NewAuthController(AuthService *service.AuthService, log log.Log) *AuthController {
	return &AuthController{
		service: AuthService,
		log:     log,
	}
}

func (controller *AuthController) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/auth").Subrouter()

	userRouter.HandleFunc("/login", controller.Login).Methods(http.MethodPost)
	userRouter.HandleFunc("/register", controller.Register).Methods(http.MethodPost)

}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var enteredInfo struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&enteredInfo)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	var user *user.User = &user.User{
		Email: enteredInfo.Email,
	}

	err = controller.service.Login(user)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User Not Found"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(enteredInfo.Password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Credentials"))
		return
	}

	var claims = &auth.Claims{
		ID:       uint64(user.ID),
		FullName: user.FullName,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}
	token, err := auth.Sign(*claims)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "auth",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 10),
		Secure:  true,
	})

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Logged In Successfully"))
}

func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user *user.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	user.Password = string(hashedPassword)

	err = controller.service.Register(user)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User Not Found Credentials"))
		return
	}

	var claims = &auth.Claims{
		ID:       uint64(user.ID),
		Email:    user.Email,
		FullName: user.FullName,
		IsAdmin:  false,
	}

	// Get Token
	token, err := auth.Sign(*claims)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "auth",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 10),
		Secure:  true,
	})
	web.RespondJSON(w, http.StatusCreated, "User Registered Successfully")
}