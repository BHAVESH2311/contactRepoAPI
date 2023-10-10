package middleware
import (

	"net/http"
	"contactapp/auth"
	"encoding/json"
)

func CurrentUserCheker(w http.ResponseWriter, r *http.Request,userId int){

	userClaim := r.Context().Value(auth.UserIDKey).(*auth.Claims)

	if userClaim.ID != uint64(userId) && !userClaim.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Not Authorized to Access")
		return
	}
}