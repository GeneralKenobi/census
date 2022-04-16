package wrap

import (
	"github.com/GeneralKenobi/census/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

type requestDataParser struct {
	ginCtx *gin.Context
}

var _ RequestData = (*requestDataParser)(nil) // Interface guard

func (parser *requestDataParser) RequireIntPathParam(paramName string) (int, error) {
	param, err := parser.RequirePathParam(paramName)
	if err != nil {
		return 0, err
	}

	paramAsInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, api.StatusBadInput.WithMessage("path parameter %s has to be an integer", paramName)
	}

	return paramAsInt, nil
}

func (parser *requestDataParser) RequirePathParam(paramName string) (string, error) {
	param := parser.ginCtx.Param(paramName)
	if param == "" {
		return "", api.StatusBadInput.WithMessage("path parameter %s is required", paramName)
	}

	return param, nil
}

func (parser *requestDataParser) BindBody(target any) error {
	err := parser.ginCtx.ShouldBindJSON(target)
	if err != nil {
		return api.StatusBadInput.WithMessageAndCause(err, "malformed request body")
	}
	err = validateRequestBody(target)
	if err != nil {
		return err
	}

	return nil
}

// validateRequestBody validates a request body. If there were validation errors it converts them into a api.StatusBadInput error.
func validateRequestBody(toValidate any) error {
	err := validate.Struct(toValidate)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return api.StatusBadInput.WithMessageAndCause(err, "invalid request body")
	}

	validationMessages := make([]string, len(validationErrs))
	for i, validationErr := range validationErrs {
		validationMessages[i] = validationErr.Namespace() + ": " + validationErr.Tag()
	}
	return api.StatusBadInput.WithMessageAndCause(err, "invalid request body: %s", strings.Join(validationMessages, ", "))
}

// Use a single instance of Validate, it caches struct info.
var validate = validator.New()
