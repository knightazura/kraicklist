package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Logger *Logger
	Writer http.ResponseWriter
}

type SuccessResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type FailedResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Code int `json:"error_code"`
	Data interface{} `json:"data"`
}

func InitResponse(logger *Logger, writer http.ResponseWriter) *Response {
	return &Response{
		Logger: logger,
		Writer: writer,
	}
}

func (r *Response) OkResponse(message string, data interface{}) {
	payload := SuccessResponse{
		Success: true,
		Message: message,
		Data: data,
	}

	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(r.Writer).Encode(payload)
	if err != nil {
		r.Logger.LogError("Failed to encode HTTP response message: %s", err.Error())
	}
}

func (r *Response) MethodNotAllowedResponse(message string) {
	r.commonErrorResponse(message, http.StatusMethodNotAllowed)
}

func (r *Response) BadRequestResponse(message string) {
	r.commonErrorResponse(message, http.StatusBadRequest)
}

func (r *Response) InternalServerErrorResponse(message string) {
	r.commonErrorResponse(message, http.StatusInternalServerError)
}

func (r *Response) commonErrorResponse(message string, code int) {
	r.Logger.LogError("%s", message)
	payload := createFailedResponse(message, code)

	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(code)
	err := json.NewEncoder(r.Writer).Encode(payload)
	if err != nil {
		r.Logger.LogError("Failed to encode HTTP response message: %s", err.Error())
	}
}

func createFailedResponse(message string, code int) FailedResponse {
	return FailedResponse{
		Success: false,
		Message: message,
		Code: code,
		Data: []string{},
	}
}