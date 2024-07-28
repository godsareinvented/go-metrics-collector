package handler

import (
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

func ProcessValidationError(error error) (string, int) {
	errors := error.(validator.ValidationErrors)
	if strings.Contains(errors[0].Field(), "Name") {
		return "Metric name not passed on or incorrect", http.StatusNotFound
	}
	return "incorrect metric data", http.StatusBadRequest
}
