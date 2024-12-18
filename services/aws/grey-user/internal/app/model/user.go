package model

import "time"

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type Address struct {
	Type        string `json:"type"`
	Address     string `json:"address"`
	Subdistrict string `json:"subdistrict"`
	District    string `json:"district"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	PostalCode  string `json:"postalCode"`
}

type Image struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Phone struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type PaymentMethod struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Number string `json:"number"`
}

type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Privileges  []string `json:"privileges"`
}

type User struct {
	UUID           string          `json:"uuid"`
	ShopID         string          `json:"shopId"`
	Email          string          `json:"email"`
	Role           string          `json:"role"`
	Firstname      string          `json:"firstname"`
	Lastname       string          `json:"lastname"`
	Birthdate      time.Time       `json:"birthdate"`
	Gender         Gender          `json:"gender"`
	Addresses      []Address       `json:"addresses"`
	Phones         []Phone         `json:"phones"`
	Image          *Image          `json:"image,omitempty"`
	PaymentMethods []PaymentMethod `json:"paymentMethods,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}
