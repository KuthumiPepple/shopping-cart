package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kuthumipepple/shopping-cart/database"
	"github.com/kuthumipepple/shopping-cart/models"
	"github.com/kuthumipepple/shopping-cart/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.OpenCollection("users")
var productCollection = database.OpenCollection("products")

func HashPassword(passwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), 12)
	if err != nil {
		log.Panic(err)
	}
	return string(hash)
}

func VerifyPassword(userPasswd, existingPasswd string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(existingPasswd), []byte(userPasswd))
	check := true
	var msg string
	if err != nil {
		msg = "invalid email or password"
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		newUser := models.User{}
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validator.New().Struct(newUser); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": newUser.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panic(err)
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user with this email already exists"})
			return
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": newUser.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panic(err)
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone numbeer is already in use"})
			return
		}

		password := HashPassword(*newUser.Password)

		newUser.Password = &password
		newUser.CreatedAt = time.Now().Local()
		newUser.UpdatedAt = time.Now().Local()
		newUser.ID = primitive.NewObjectID()
		newUser.UserID = newUser.ID.Hex()

		token, refreshToken, err := tokens.GenerateTokens(*newUser.Email, *newUser.FirstName, *newUser.LastName, newUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while generating tokens"})
			log.Panic(err)
		}

		newUser.Token = &token
		newUser.RefreshToken = &refreshToken
		newUser.UserCart = make([]models.UserProduct, 0)
		newUser.AddressDetails = make([]models.Address, 0)
		newUser.OrderStatus = make([]models.Order, 0)

		_, insertionErr := userCollection.InsertOne(ctx, newUser)
		if insertionErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully created the user!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := models.User{}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		foundUser := models.User{}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
			return
		}
		isPasswordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !isPasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, err := tokens.GenerateTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while generating tokens"})
			log.Panic(err)
		}
		tokens.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		c.JSON(http.StatusOK, gin.H{"JWT token": token})
	}
}

func AddProductAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		product.ProductID = primitive.NewObjectID()
		_, err := productCollection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "product not inserted"})
			return
		}

		c.JSON(http.StatusOK, "Successfully inserted new product!")

	}
}

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
