package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      *string            `json:"first_name" bson:"first_name" validate:"required,min=2,max=30"`
	LastName       *string            `json:"last_name" bson:"last_name" validate:"required,min=2,max=30"`
	Password       *string            `validate:"required,min=6"`
	Email          *string            `validate:"email,required"`
	Phone          *string            `validate:"required"`
	Token          *string
	RefreshToken   *string   `json:"refresh_token" bson:"refresh_token"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
	UserID         string    `json:"user_id" bson:"user_id"`
	UserCart       []UserProduct
	AddressDetails []Address `json:"addresses" bson:"addresses"`
	OrderStatus    []Order   `json:"orders" bson:"orders"`
}

type Product struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       *uint64
	Rating      *uint8
	Image       *string
}

type UserProduct struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json:"product_name" bson:"product_name"`
	Price       int
	Rating      *uint
	Image       *string
}

type Address struct {
	AddressID primitive.ObjectID `bson:"_id"`
	House     *string            `json:"house_name" bson:"house_name"`
	Street    *string            `json:"street_name" bson:"street_name"`
	City      *string            `json:"city_name" bson:"city_name"`
	ZipCode   *string            `json:"zip_code" bson:"zip_code"`
}

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id"`
	OrderCart     []UserProduct      `json:"order_list" bson:"order_list"`
	OrderedAt     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price         int                `json:"total_price" bson:"total_price"`
	Discount      *int
	PaymentMethod Payment `json:"payment_method" bson:"payment_method"`
}
type Payment struct {
	Digital        bool
	CashOnDelivery bool
}
