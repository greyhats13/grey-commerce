// Path: grey-user/internal/app/model/user.go

package model

import (
	"time"
)

// User untuk DynamoDB, menambahkan tag `dynamodbav:"..."`.
type User struct {
	UUID          string          `json:"uuid" dynamodbav:"uuid"`
	ShopID        string          `json:"shopID" dynamodbav:"shopID" validate:"required"`
	Email         string          `json:"email" dynamodbav:"email" validate:"required,email"`
	Role          string          `json:"role" dynamodbav:"role" validate:"required"`
	Firstname     string          `json:"firstname" dynamodbav:"firstname" validate:"required"`
	Lastname      string          `json:"lastname" dynamodbav:"lastname" validate:"required"`
	Gender        string          `json:"gender" dynamodbav:"gender" validate:"required"`
	Birthdate     time.Time       `json:"birthdate" dynamodbav:"birthdate" validate:"required"`
	Addresses     []Address       `json:"addresses" dynamodbav:"addresses" validate:"required,dive"`
	Phones        []Phone         `json:"phones" dynamodbav:"phones" validate:"required,dive"`
	PaymentMethod []PaymentMethod `json:"paymentMethod" dynamodbav:"paymentMethod" validate:"required,dive"`
	Image         Image           `json:"image" dynamodbav:"image" validate:"required"`
	CreatedAt     time.Time       `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt" dynamodbav:"updatedAt"`
}

// Contoh struct detail lain
type Address struct {
	Label    string `json:"label" dynamodbav:"label" validate:"required"`
	Street   string `json:"street" dynamodbav:"street" validate:"required"`
	City     string `json:"city" dynamodbav:"city" validate:"required"`
	Postcode string `json:"postcode" dynamodbav:"postcode" validate:"required"`
	Country  string `json:"country" dynamodbav:"country" validate:"required"`
}

type Phone struct {
	CountryCode string `json:"countryCode" dynamodbav:"countryCode" validate:"required"`
	Number      string `json:"number" dynamodbav:"number" validate:"required"`
}

type PaymentMethod struct {
	Type        string `json:"type" dynamodbav:"type" validate:"required"`
	CardNumber  string `json:"cardNumber" dynamodbav:"cardNumber" validate:"required"`
	ExpiryMonth int    `json:"expiryMonth" dynamodbav:"expiryMonth" validate:"required"`
	ExpiryYear  int    `json:"expiryYear" dynamodbav:"expiryYear" validate:"required"`
}

type Image struct {
	URL  string `json:"url" dynamodbav:"url" validate:"required,url"`
	Desc string `json:"desc" dynamodbav:"desc"`
}
