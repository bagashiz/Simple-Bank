package util

// constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	IDR = "IDR"
	// add more
)

// IsSupportedCurrency returns true if the currency is supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, IDR:
		return true
	}
	return false
}
