package rest

import (
	"net/http"
	"encoding/json"
)

type JSONResponseBody interface{}

// JSON rest structure wrapper over native rest writer
type JSONResponse struct {
	Body string
	Content interface{}
	Status int
	Headers map[string]string
	ContentType string
}

// Set HTTP Status
func (r *JSONResponse) SetStatus(httpCode int) *JSONResponse {
	r.Status = httpCode
	return r
}

// Set HTTP headers
func (r *JSONResponse) SetHeader(header string, value string) *JSONResponse {
	r.Headers[header] = value

	return r
}

// Set content
func (r *JSONResponse) SetContent(data interface{}) *JSONResponse {
	r.Content = data
	return r
}

// Export rest content
func (r *JSONResponse) Write(responseWriter *http.ResponseWriter) *http.ResponseWriter {
	w := *responseWriter

	data, err := json.Marshal(r.Content)

	w.Header().Add("Content-Type", "application/json")

	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		data, _ = json.Marshal(ErrorResponse{err.Error()})
	} else {
		w.WriteHeader(r.Status)

		// Apply headers
		for k, v := range r.Headers {
			w.Header().Add(k, v)
		}
	}

	w.Write(data);

	return &w;
}
