//Costum validator for json request
package costumValidation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var BookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}
