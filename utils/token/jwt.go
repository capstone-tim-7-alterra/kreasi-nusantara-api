package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TokenUtil interface {
	GenerateToken(id uuid.UUID, role string) (string, error)
	GetClaims(c echo.Context) *JWTClaim
}

type tokenUtil struct{}

func NewTokenUtil() *tokenUtil {
	return &tokenUtil{}
}

func (*tokenUtil) GenerateToken(id uuid.UUID, role string) (string, error) {
	claims := JWTClaim{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := unsignedToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (*tokenUtil) GetClaims(c echo.Context) *JWTClaim {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaim)
	return claims
}
