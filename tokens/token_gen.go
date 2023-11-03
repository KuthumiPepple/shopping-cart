package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kuthumipepple/shopping-cart/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserID    string
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
var userCollection = database.OpenCollection("users")

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

func UpdateAllTokens(signedToken, signedRefreshToken, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: time.Now().Local()})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	filter := bson.M{"user_id": userID}
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	if err != nil {
		log.Panic(err)
	}
}
