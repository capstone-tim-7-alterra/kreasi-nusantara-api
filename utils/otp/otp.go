package otp

import "crypto/rand"

type OTPUtil interface {
	GenerateOTP(length int) (string, error)
}

type otpUtil struct {}

func NewOTPUtil() *otpUtil {
	return &otpUtil{}
}

func (o *otpUtil) GenerateOTP(length int) (string, error) {
	const charset = "0123456789"

	otp := make([]byte, length)
	_, err := rand.Read(otp)
	if err != nil {
		return "", err
	}
	for i := range otp {
		otp[i] = charset[otp[i]%byte(len(charset))]
	}
	return string(otp), nil
}