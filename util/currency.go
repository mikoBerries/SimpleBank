//Enum of currency that support in apps
package util

//add your new enum here
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	IDR = "IDR"
)

//CheckCurrencySupport consume string returning true when currency is supported
func CheckCurrencySupport(currency string) bool {
	switch currency {
	case USD, EUR, CAD, IDR:
		return true
	}
	return false
}
