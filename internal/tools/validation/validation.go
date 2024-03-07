package validation

import (
	"errors"
	"fmt"
	"lamoda_task/internal/tools/requests/goods"

	"github.com/go-playground/validator/v10"
)

func ValidateRequestStruct(req goods.Request) (string, error) {
	validate := validator.New()

	for _, dataVal := range req.Data {
		if err := validate.Struct(dataVal); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				for _, fieldError := range validationErrors {
					return dataVal.ID, fmt.Errorf("validation error in field %s: %s", fieldError.Field(), fieldError.Tag())
				}
			}

			return dataVal.ID, err
		}
	}

	return "", nil
}
