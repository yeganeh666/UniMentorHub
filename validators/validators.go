package validators

import (
	"Atrovan_Q1/services/models"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Init initialize the validator.
func Init() {
	validate = validator.New()
}

// Validate validates user input, based on the given model.
func Validate(model interface{}) (bool, models.Response) {
	err := validate.Struct(model)
	targets := []models.Target{}

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			targets = append(targets, models.Target{
				Name:        err.Field(),
				Description: err.Tag(),
			})
		}

		errResponse, err := json.Marshal(models.ErrorResponse{
			Code:      "BadArgument",
			Message:   "Your entered data does not meet the minimum requirements for assessment.",
			MessageFa: "داده‌های واردشده توسط شما حداقل الزامات ارزیابی را برآورده نمی‌کند.",
			Target:    targets,
		})

		if err != nil {
			return false, models.ServerDownResponse
		}
		return false, models.Response{StatusCode: http.StatusBadRequest, Body: errResponse}
	}
	return true, models.Response{}
}
