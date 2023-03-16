package util

const (
	usd    = "usd"
	euro   = "euro"
	birr   = "birr"
	shilng = "shilng"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case usd, euro, birr, shilng:
		return true
	}
	return false
}
