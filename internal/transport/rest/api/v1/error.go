package v1

import "net/http"

var (
	ErrorInvalidQueryParam = &Response{
		Code:  http.StatusBadRequest,
		Error: "Invalid query parameters",
	}

	ErrorInvalidRequestBody = &Response{
		Code:  http.StatusBadRequest,
		Error: "Invalid request body",
	}

	ErrorInvalidPathParam = &Response{
		Code:  http.StatusBadRequest,
		Error: "Invalid path parameters",
	}
)
