package auth

import (
	"contactapp/errors"
	"time"
	"context"
	"net/http"
	"encoding/json"

	//"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID       uint64
	FullName string
	Email    string
	IsAdmin  bool
	jwt.StandardClaims
}

type contextKey string

const UserIDKey contextKey = "userId"

var secretKeyJWT = []byte("SecretJwt")

func Sign(claims  Claims)(string,error){
	claims.StandardClaims.ExpiresAt=time.Now().Add(time.Minute *20).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedToken,err:=token.SignedString(secretKeyJWT)
	if err!=nil{
		return "",err
	}
	return signedToken,nil
}


func Verify(token string)(*Claims,error){

	var userClaims = &Claims{}

	tokenObj,err:=jwt.ParseWithClaims(token,userClaims,func(token *jwt.Token)(interface{},error){
		return secretKeyJWT,nil
	})
	if err!=nil{
			return nil,errors.NewValidationError("unauthorized access")
	}

	if !tokenObj.Valid {
		return nil, errors.NewValidationError("token invalid")
	}
	return userClaims,nil

}


func Protect(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
	
		token := r.Header.Values("auth")

		if len(token) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("No Token Present"))
			return
		}
		payload, err := Verify(token[0])

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid Token")
			return
		}

		con := context.WithValue(r.Context(), UserIDKey, payload)
		newReq := r.WithContext(con)

		next.ServeHTTP(w, newReq)
	})
}


func IsAdmin(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		user := r.Context().Value(UserIDKey).(*Claims)

		if !user.IsAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not Authorized to access this Route"))
			return
		}

		h.ServeHTTP(w, r)
	})
}




