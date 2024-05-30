package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TokenUtil interface {
	GenerateToken(username string, isAdmin bool, isSuperAdmin bool) (string, error)
	VerifyToken(token string) (*JWTClaim, error)
	GetClaims(c echo.Context) *JWTClaim
}

type tokenUtil struct{}

func NewTokenUtil() *tokenUtil {
	return &tokenUtil{}
}

func (*tokenUtil) GenerateToken(username string, isAdmin bool, isSuperAdmin bool) (string, error) {
	claims := JWTClaim{
		Username:     username,
		IsAdmin:      isAdmin,
		IsSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := unsignedToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (*tokenUtil) VerifyToken(token string) (*JWTClaim, error) {
	claims := &JWTClaim{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (*tokenUtil) GetClaims(c echo.Context) *JWTClaim {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaim)
	return claims
}