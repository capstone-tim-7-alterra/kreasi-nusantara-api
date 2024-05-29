package error

import (
	"errors"
	"kreasi-nusantara-api/constants/message"
)

var (
	// Password
	ErrFailedHashingPassword = errors.New(message.FAILED_HASHING_PASSWORD)
	ErrPasswordMismatch      = errors.New(message.PASSWORD_MISMATCH)

	// External Service
	ErrExternalServiceError = errors.New(message.EXTERNAL_SERVICE_ERROR)

	// Forbidden
	ErrForbiddenResource = errors.New(message.FORBIDDEN_RESOURCE)

	// Token
	ErrFailedGenerateToken = errors.New(message.FAILED_GENERATE_TOKEN)

	// DuplicateKey 
	ErrDuplicateKey = errors.New(message.DUPLICATE_KEY)
)
