// Path: grey-user/internal/app/model/user.go

package model

import (
	"time"
)

// User represents a user in the system with validation tags.
type User struct {
	UserId         string          `json:"userId" dynamodbav:"userId"`
	ShopID         string          `json:"shopId" dynamodbav:"shopId" validate:"required"`
	Email          string          `json:"email" dynamodbav:"email" validate:"required,email"`
	Role           string          `json:"role" dynamodbav:"role" validate:"required"`
	Firstname      string          `json:"firstname" dynamodbav:"firstname" validate:"required"`
	Lastname       string          `json:"lastname" dynamodbav:"lastname" validate:"required"`
	Gender         string          `json:"gender" dynamodbav:"gender" validate:"required"`
	Birthdate      time.Time       `json:"birthdate" dynamodbav:"birthdate"`
	Addresses      []Address       `json:"addresses" dynamodbav:"addresses"`
	Phones         []Phone         `json:"phones" dynamodbav:"phones"`
	PaymentMethods []PaymentMethod `json:"paymentMethods" dynamodbav:"paymentMethods"`
	Image          Image           `json:"image" dynamodbav:"image"`
	CreatedAt      time.Time       `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt" dynamodbav:"updatedAt"`
}

// Address represents a user's address with validation tags.
// If you want each subfield optional too, remove "required" or replace with "omitempty"
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

// Phone represents a user's phone. Remove "required" if you want them optional
type Phone struct {
	Type   string `json:"type" dynamodbav:"type" validate:"required"`
	Code   string `json:"code" dynamodbav:"code" validate:"required"`
	Number string `json:"number" dynamodbav:"number" validate:"required"`
}

// PaymentMethod represents a user's payment method
type PaymentMethod struct {
	Type        string `json:"type" dynamodbav:"type" validate:"required"`
	Name        string `json:"name" dynamodbav:"name" validate:"required"`
	Number      string `json:"number" dynamodbav:"number" validate:"required"`
	ExpiryMonth int    `json:"expiryMonth" dynamodbav:"expiryMonth" validate:"omitempty,min=1,max=12"`
	ExpiryYear  int    `json:"expiryYear" dynamodbav:"expiryYear" validate:"omitempty,min=2023"`
}

// Image represents a user's image
type Image struct {
	Name string `json:"name" dynamodbav:"name" validate:"required"`
	URL  string `json:"url" dynamodbav:"url" validate:"required,url"`
	Desc string `json:"desc" dynamodbav:"desc"`
}
