package controllers

import "go.mongodb.org/mongo-driver/mongo"

type Application struct {
	productsCollection *mongo.Collection
	usersCollection    *mongo.Collection
}

func NewApplication(productsCollectionName, usersCollectionName *mongo.Collection) *Application {
	return &Application{
		productsCollection: productsCollectionName,
		usersCollection:    usersCollectionName,
	}
}
