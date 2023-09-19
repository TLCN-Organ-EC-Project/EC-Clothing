package util

const (
	S        = "S"
	M        = "M"
	L        = "L"
	XL       = "XL"
	XXL      = "XXL"
	OVERSIZE = "OVERSIZE"
)

// IsSupportedCurrency return true if the currency is supported
func IsSupportedSize(size string) bool {
	switch size {
	case S, M, L, XL, XXL, OVERSIZE:
		return true
	}
	return false
}
