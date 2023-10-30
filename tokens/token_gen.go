package tokens

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserID    string
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

func GenerateTokens(email, firstName, lastName, userID string) (signedToken, signedRefreshToken string, err error) {
	claims := SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Duration(24) * time.Hour)),
		},
	}

	refreshClaims := SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Duration(168) * time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}
