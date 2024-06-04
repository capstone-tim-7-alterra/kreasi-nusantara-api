package middlewares

import (
	msg "kreasi-nusantara-api/constants/message"
	http_util "kreasi-nusantara-api/utils/http"
	"kreasi-nusantara-api/utils/token"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := token.NewTokenUtil().GetClaims(c)
		if claims.Role != strings.ToLower("admin"){
			return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.UNAUTHORIZED)
		}

		return next(c)
	}
}

func IsSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := token.NewTokenUtil().GetClaims(c)
		if claims.Role != strings.ToLower("super_admin"){
			return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.UNAUTHORIZED)
		}

		return next(c)
	}
}

func IsAdminOrSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        claims := token.NewTokenUtil().GetClaims(c)
        role := strings.ToLower(claims.Role)
        if role != "admin" && role != "super_admin" {
            return http_util.HandleErrorResponse(c, http.StatusUnauthorized, msg.UNAUTHORIZED)
        }

        return next(c)
    }
}