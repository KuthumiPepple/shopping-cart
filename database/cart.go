package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/kuthumipepple/shopping-cart/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrProductNotFound        = errors.New("product not found in database")
	ErrFailedToDecodeProducts = errors.New("cannot decode products into slice")
	ErrInvalidUserID          = errors.New("user is not valid")
	ErrFailedToUpdateUser     = errors.New("cannot add product to cart")
	ErrFailedToRemoveItem     = errors.New("cannot remove item from cart")
)

func AddProductToCart(ctx context.Context, productsCollection, usersCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	resultSet, err := productsCollection.Find(ctx, bson.M{"_id": productID})
	if err != nil {
		log.Println(err)
		return ErrProductNotFound

	}
	var productCart []models.UserProduct
	err = resultSet.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrFailedToDecodeProducts
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserID
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "usercart", Value: bson.D{
				{Key: "$each", Value: productCart},
			}},
		}},
	}

	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return ErrFailedToUpdateUser
	}
	return nil
}

func RemoveItemFromCart(ctx context.Context, usersCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserID
	}

	update := bson.M{
		"$pull": bson.M{
			"usercart": bson.M{
				"_id": productID,
			},
		},
	}

	_, err = usersCollection.UpdateMany(
		ctx,
		bson.M{"_id": id},
		update,
	)
	if err != nil {
		return ErrFailedToRemoveItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, productsCollection, usersCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserID
	}

	var productDetails models.UserProduct
	err = productsCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
		return ErrProductNotFound
	}

	var orderDetails models.Order
	orderDetails.OrderID = primitive.NewObjectID()
	orderDetails.OrderedAt = time.Now().Local()
	orderDetails.OrderCart = append(orderDetails.OrderCart, productDetails)
	orderDetails.PaymentMethod.CashOnDelivery = true
	orderDetails.Price = productDetails.Price

	filter := bson.M{"_id": id}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "orders", Value: orderDetails},
		}},
	}
	_, err = usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
