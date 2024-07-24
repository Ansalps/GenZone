package helper

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

func validateNameOrInitials(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// Regex pattern for initials (e.g., "J. K.")
	initialsRegex := regexp.MustCompile(`^([A-Z]\. )*[A-Z]\.$`)
	// Regex pattern for a regular name (e.g., "Doe")
	nameRegex := regexp.MustCompile(`^[A-Za-z]+$`)
	return initialsRegex.MatchString(value) || nameRegex.MatchString(value)
}
func Validate(value interface{}) error {
	// var translator = map[string]string{
	// 	"Name_required":            "Please enter  Name",
	// 	"Password_required":        "Please enter  Password",
	// 	"ConfirmPassword_required": "Please enter  ConfirmPassword",
	// 	"Email_email":              "Please enter a valid email address",
	// 	"UserID_required":          "Please enter a valid user id",
	// 	"AddressID_required":       "Please enter a valid address id",
	// 	"AddressID_number":         "Please enter a numerical value for address id",
	// 	"UserID_number":            "Please enter a numerical value for user id",
	// 	"ProductID_number":         "Please enter a numerical value for product id",
	// }
	// validate the struct body
	validate := validator.New()
	validate.RegisterValidation("nameOrInitials", validateNameOrInitials)
	err := validate.Struct(value)
	if err != nil {
		// var errs []string
		for _, e := range err.(validator.ValidationErrors) {
			// translationKey := e.Field() + "_" + e.Tag()
			// errMsg := translator[translationKey]
			// if errMsg == "" {
			// 	errMsg = e.ActualTag()
			// }
			// errs = append(errs, errMsg)
			switch e.Tag() {
			case "required":
				return fmt.Errorf("%s is required", e.Field())
			case "email":
				return fmt.Errorf("%s is not a valid email address", e.Field())
			case "numeric":
				return fmt.Errorf("%s shouls contain only digits", e.Field())
			case "len":
				return fmt.Errorf("%s shouls have a length of %s", e.Field(), e.Param())
			case "min":
				return fmt.Errorf("%s shouls have a minimum length of %s", e.Field(), e.Param())
			case "excludesall":
				return fmt.Errorf("%s shouls not contain space", e.Field())
			case "nameOrInitials":
				return fmt.Errorf("%s should be either initials or a regular name", e.Field())
			default:
				return fmt.Errorf("validation error for field %s", e.Field())
			}
		}
	}
	//return errors.New(strings.Join(errs, ", "))
	return nil
}
