package response

import (
	"fmt"
	"log/slog"
	"net/http"
	"restrocheck/pkg/logger/sl"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "ОК"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors, log *slog.Logger, w http.ResponseWriter, r *http.Request, status int, err error, message string) {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "phone":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid phone", err.Field()))
		case "salary":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid salary", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	log.Error(message, sl.Err(err))
	w.WriteHeader(status)
	render.JSON(w, r, Error(strings.Join(errMsgs, ", ")))
}

func RespondWithError(log *slog.Logger, w http.ResponseWriter, r *http.Request, status int, err error, message string) {
	log.Error(message, sl.Err(err))
	w.WriteHeader(status)
	render.JSON(w, r, Error(message))
}