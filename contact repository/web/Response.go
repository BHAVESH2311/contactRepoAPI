package web

import(
	"contactapp/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(response))
}
// RespondError check error type and Write to ResponseWriter.
func RespondError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *errors.ValidationError:
		RespondJSON(w, http.StatusBadRequest, err)
	case *errors.HTTPError:
		httpError := err.(*errors.HTTPError)
		RespondJSON(w, httpError.HTTPStatus, httpError.ErrorKey)
	default:
		RespondErrorMessage(w, http.StatusInternalServerError, err.Error())
	}
}
// RespondErrorMessage make error response with payload.
func RespondErrorMessage(w http.ResponseWriter, code int, msg string) {
	RespondJSON(w, code, map[string]string{"error": msg})
}
// RespondJSONWithXTotalCount Make response with json format and add X-Total-Count header.
func RespondJSONWithXTotalCount(w http.ResponseWriter, code int, count int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	SetNewHeader(w, "X-Total-Count", strconv.Itoa(count))
	w.WriteHeader(code)
	w.Write([]byte(response))
}
// 	SetNewHeader(w,"total","10") will set header "total" : "10"
func SetNewHeader(w http.ResponseWriter, headerName, value string) {
	w.Header().Add("Access-Control-Expose-Headers", headerName)
	w.Header().Set(headerName, value)
}
