//Costum validator for json request
package costumValidation

import (
	"time"

	"github.com/MikoBerries/SimpleBank/util"
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

//Currency check
var IsCurrency validator.Func = func(fl validator.FieldLevel) bool {
	//pick data and assert it to specific type string/int/ etc...
	currency, ok := fl.Field().Interface().(string)
	if ok && util.CheckCurrencySupport(currency) {
		//do some check logic
		// if {
		return true
		// }
	}
	return false
}
