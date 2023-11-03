package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	UserId         string `json:"userId" validate:"required"`
	Token          string `json:"token"`
	RefreshToken   string `json:"refreshToken"`
	StandardClaims jwt.StandardClaims
}

func (s *SignedDetails) Valid() error {
	if time.Now().Unix() > s.StandardClaims.ExpiresAt {
		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}
	return nil
}

func GetAllToken(userId string) (token string, refreshToken string, tokenError error) {
	secret := os.Getenv("JWT_KEY")
	claims := &SignedDetails{
		UserId:         userId,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()},
	}

	refreshClaims := &SignedDetails{
		UserId:         userId,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix()},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	if err != nil {
		return "", "", err
	}

	refreshToken, refreshTokenError := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret))

	if refreshTokenError != nil {
		return "", "", refreshTokenError
	}

	return token, refreshToken, nil
}
