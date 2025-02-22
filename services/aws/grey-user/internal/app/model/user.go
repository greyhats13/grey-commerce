// Path: grey-user/internal/app/model/user.go

package model

import (
	"time"
)

// User represents a user in the system with validation tags.
type User struct {
	UUID          string          `json:"uuid" dynamodbav:"uuid"`
	ShopID        string          `json:"shopID" dynamodbav:"shopID" validate:"required"`
	Email         string          `json:"email" dynamodbav:"email" validate:"required,email"`
	Role          string          `json:"role" dynamodbav:"role" validate:"required"`
	Firstname     string          `json:"firstname" dynamodbav:"firstname" validate:"required"`
	Lastname      string          `json:"lastname" dynamodbav:"lastname" validate:"required"`
	Gender        string          `json:"gender" dynamodbav:"gender" validate:"required"`
	Birthdate     time.Time       `json:"birthdate" dynamodbav:"birthdate" validate:"datetime=2006-01-02"`
	Addresses     []Address       `json:"addresses" dynamodbav:"addresses"`
	Phones        []Phone         `json:"phones" dynamodbav:"phones"`
	PaymentMethod []PaymentMethod `json:"paymentMethod" dynamodbav:"paymentMethod"`
	Image         Image           `json:"image" dynamodbav:"image"`
	CreatedAt     time.Time       `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt" dynamodbav:"updatedAt"`
}

// Address represents a user's address with validation tags.
type Address struct {
	Type        string `json:"type" dynamodbav:"type" validate:"required"`
	Address     string `json:"address" dynamodbav:"address" validate:"required"`
	Subdistrict string `json:"subdistrict" dynamodbav:"subdistrict" validate:"required"`
	District    string `json:"district" dynamodbav:"district" validate:"required"`
	City        string `json:"city" dynamodbav:"city" validate:"required"`
	Province    string `json:"province" dynamodbav:"province" validate:"required"`
	Country     string `json:"country" dynamodbav:"country" validate:"required"`
	PostalCode  int    `json:"postalCode" dynamodbav:"postalCode" validate:"required"`
}

// Phone represents a user's phone with validation tags.
type Phone struct {
	Type   string `json:"type" dynamodbav:"type" validate:"required"`
	Code   string `json:"code" dynamodbav:"countryCode" validate:"required"`
	Number string `json:"number" dynamodbav:"number" validate:"required"`
}

// PaymentMethod represents a user's payment method with validation tags.
type PaymentMethod struct {
	Type        string `json:"type" dynamodbav:"type" validate:"required"`
	Number      string `json:"cardNumber" dynamodbav:"number" validate:"required"`
	ExpiryMonth int    `json:"expiryMonth" dynamodbav:"expiryMonth" validate:"omitempty,min=1,max=12"`
	ExpiryYear  int    `json:"expiryYear" dynamodbav:"expiryYear" validate:"omitempty,min=2023"`
}

// Image represents a user's image with validation tags.
type Image struct {
	Name string `json:"name" dynamodbav:"name" validate:"required"`
	URL  string `json:"url" dynamodbav:"url" validate:"required,url"`
	Desc string `json:"desc" dynamodbav:"desc"`
}
