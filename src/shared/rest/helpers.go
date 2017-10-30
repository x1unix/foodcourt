package rest

import "net/http"

// Create a new success rest
func Success(data interface{}) *JSONResponse {
	return &JSONResponse{
		Status: http.StatusOK,
		Content: data,
		ContentType: "application/json",
	}
}

func Echo(message string) *JSONResponse {
	var data JSONResponseBody = ErrorResponse{message}
	var ptrData = &data

	return Success(ptrData)
}

// Generate error rest
func ErrorFromString(errorMessage string, code int) *JSONResponse {
	if (code == 0) {
		code = http.StatusInternalServerError
	}

	var data JSONResponseBody = ErrorResponse{errorMessage}
	var ptrData = &data

	return &JSONResponse{
		Status: code,
		ContentType: "application/json",
		Content: ptrData,
	}
}

// Create a new error rest from app error
func Error(err error) *JSONResponse {
	return ErrorFromString(err.Error(), 500)
}