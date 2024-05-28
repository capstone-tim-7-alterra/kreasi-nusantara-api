package middlewares

import (
	msg "kreasi-nusantara-api/constants/message"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.MISSING_TOKEN)
		}
		if !strings.HasPrefix(authHeader, "Bearer") {
			return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.INVALID_AUTH_TYPE)
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer")

		token, err := token.NewTokenUtil().VerifyToken(tokenStr)
		if err != nil {
			return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.INVALID_TOKEN)
		}
		c.Set("user", token)

		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := token.NewTokenUtil().GetClaims(c)
		if !claims.IsAdmin {
			return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.UNAUTHORIZED)
		}

		return next(c)
	}
}

func IsSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := token.NewTokenUtil().GetClaims(c)
		if !claims.IsSuperAdmin {
			return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.UNAUTHORIZED)
		}

		return next(c)
	}
}