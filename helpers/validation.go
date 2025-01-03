package helpers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ValidateBody[T any](validate *validator.Validate, data T) (Response, error) {
	validationErrors := []ValidationError{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ValidationError

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
		return Response{
			Data:       nil,
			Message:    "Error Validation",
			Status:     http.StatusBadRequest,
			Validation: validationErrors,
		}, errs
	}

	return Response{
		Data:       data,
		Message:    "success",
		Status:     http.StatusOK,
		Validation: validationErrors,
	}, nil
}
