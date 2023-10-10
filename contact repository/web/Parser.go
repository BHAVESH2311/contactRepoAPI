package web

import (
	"net/http"
	"strconv"
	"strings"
	"contactapp/errors"
)

func LimitOffsetExtract(r *http.Request) (int, int, error) {

	pageStr := r.URL.Query().Get("page")
	if len(pageStr) == 0 {
		pageStr = "1"
	}

	limitStr := r.URL.Query().Get("limit") 
	if len(limitStr) == 0 {
		limitStr = "5"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return -1, -1, err
	}

	

	limit, err := strconv.Atoi(limitStr)
	
	if err != nil || limit < 1 {
		return -1, -1, err
	}

	if limit < 1 {
		return -1, -1, errors.NewValidationError("Limit cannot be less than 1.")
	}

	if limit > 100 {
		return -1, -1, errors.NewValidationError("Limit cannot be more than 100.")
	}

	// Calculate the offset
	offset := (page - 1) * limit

	return limit, offset, nil
}

func ParseIncludes(r *http.Request) []string {
	query := r.URL.Query()
	includes := query.Get("includes")

	if includes == "" {
		return []string{}
	}
	return strings.Split(includes, ",")
}