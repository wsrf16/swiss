package totpkit

func GenerateGoogleCode(secret string) (string, error) {
	return GenerateCodeByPeriod(secret, 30, 6)
}

func CheckGoogleCode(secret string, targetCode string) bool {
	return CheckCode(secret, 30, 6, targetCode)
}
