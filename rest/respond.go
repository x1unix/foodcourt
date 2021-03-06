package rest

import "net/http"

type RequestHandler func(w http.ResponseWriter, r *http.Request)

const msgOk = "OK"
const msgBadRequest = "Bad Request"
const msgNotFound = "Not Found"

// Create a new success rest
func Success(data interface{}) *JSONResponse {
	return &JSONResponse{
		Status:      http.StatusOK,
		Content:     data,
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
	if code == 0 {
		code = http.StatusInternalServerError
	}

	var data JSONResponseBody = ErrorResponse{errorMessage}
	var ptrData = &data

	return &JSONResponse{
		Status:      code,
		ContentType: "application/json",
		Content:     ptrData,
	}
}

func NotFound(w *http.ResponseWriter) {
	ErrorFromString(msgNotFound, 404).Write(w)
}

func BadRequest(w *http.ResponseWriter, err string) {

	if len(err) == 0 {
		err = msgBadRequest
	}

	ErrorFromString(err, 400).Write(w)
}

func Ok(w *http.ResponseWriter) {
	Echo(msgOk).Write(w)
}

// Create a new error rest from app error
func Error(err error) *JSONResponse {
	return ErrorFromString(err.Error(), 500)
}

func HttpError(err error, code int) *JSONResponse {
	return HttpErrorFromString(err.Error(), code)
}

func HttpErrorFromString(err string, code int) *JSONResponse {
	return ErrorFromString(err, code)
}
