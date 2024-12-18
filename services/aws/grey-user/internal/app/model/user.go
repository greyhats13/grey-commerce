//Path: grey-user/internal/app/model/user.go

package model

import "time"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type Address struct {
	Type        string `json:"type" dynamodbav:"type"`
	Address     string `json:"address" dynamodbav:"address"`
	Subdistrict string `json:"subdistrict" dynamodbav:"subdistrict"`
	District    string `json:"district" dynamodbav:"district"`
	City        string `json:"city" dynamodbav:"city"`
	Province    string `json:"province" dynamodbav:"province"`
	Country     string `json:"country" dynamodbav:"country"`
	PostalCode  string `json:"postalCode" dynamodbav:"postalCode"`
}

type Image struct {
	Name string `json:"name" dynamodbav:"name"`
	URL  string `json:"url" dynamodbav:"url"`
}

type Phone struct {
	Type   string `json:"type" dynamodbav:"type"`
	Number string `json:"number" dynamodbav:"number"`
}

type PaymentMethod struct {
	Type   string `json:"type" dynamodbav:"type"`
	Name   string `json:"name" dynamodbav:"name"`
	Number string `json:"number" dynamodbav:"number"`
}

type User struct {
	UUID           string          `json:"uuid" dynamodbav:"uuid"`
	ShopID         string          `json:"shopId" dynamodbav:"shopId"`
	Email          string          `json:"email" dynamodbav:"email"`
	Role           string          `json:"role" dynamodbav:"role"`
	Firstname      string          `json:"firstname" dynamodbav:"firstname"`
	Lastname       string          `json:"lastname" dynamodbav:"lastname"`
	Birthdate      time.Time       `json:"birthdate" dynamodbav:"birthdate"`
	Gender         Gender          `json:"gender" dynamodbav:"gender"`
	Addresses      []Address       `json:"addresses" dynamodbav:"addresses"`
	Phones         []Phone         `json:"phones" dynamodbav:"phones"`
	Image          *Image          `json:"image,omitempty" dynamodbav:"image,omitempty"`
	PaymentMethods []PaymentMethod `json:"paymentMethods,omitempty" dynamodbav:"paymentMethods,omitempty"`
	CreatedAt      time.Time       `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt" dynamodbav:"updatedAt"`
}
