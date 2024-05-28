package token

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	http_util "kreasi-nusantara-api/utils/http"
	msg "kreasi-nusantara-api/constants/message"
)

type JWTClaim struct {
	Username     string `json:"username"`
	IsAdmin      bool   `json:"is_admin"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	jwt.RegisteredClaims
}

func GetJWTConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTClaim)
		},
		ErrorHandler: jwtErrorHandler,
		SigningKey:   []byte(os.Getenv("JWT_KEY")),
	}
}

func jwtErrorHandler(c echo.Context, err error) error {
	code := http.StatusUnauthorized

	if errors.Is(err, echojwt.ErrJWTInvalid) {
		return http_util.HandleErrorResponse(
			c,
			code,
			msg.INVALID_TOKEN,
		)
	}

	return http_util.HandleErrorResponse(
		c,
		code,
		msg.UNAUTHORIZED,
	)
}