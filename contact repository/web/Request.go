package web

import (
	"contactapp/errors"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"net/url"
)

// UnmarshalJSON parses data from request and return otherwise error return.
func UnmarshalJSON(request *http.Request, out interface{}) error {
	if request.Body == nil {
		fmt.Println("==============================err request.Body == nil==========================")
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}
	// fmt.Println("==============================err (request.Body)==========================")
	// fmt.Println(request.Body)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("==============================err ioutil.ReadAll==========================")
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	if len(body) == 0 {
		fmt.Println("==============================err len(body) == 0==========================")
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		fmt.Println("==============================err json.Unmarshal==========================")
		fmt.Println(body)
		fmt.Println(out)
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}

func ParseLimitOffset(query url.Values) (int, int) {
	limit, err := strconv.Atoi(query.Get("limit"))
	var iLimit, iOffset int = -1, -1
	if err == nil {
		iLimit = int(limit)
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err == nil {
		iOffset = int(offset)
	}
	return iLimit, iOffset
}