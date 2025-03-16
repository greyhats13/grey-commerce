// Path: grey-user/internal/app/model/user.go

package model

import (
	"fmt"
	"time"
)

// Date is a custom type for handling date-only fields
type Date struct {
	time.Time
}

// ParseDate parses a string into a Date type
func ParseDate(dateStr string) (Date, error) {
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Date{}, err
	}
	return Date{parsedTime}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Date) UnmarshalJSON(b []byte) error {
	// Trim quotes
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("invalid date format")
	}
	s = s[1 : len(s)-1]
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (d Date) MarshalJSON() ([]byte, error) {
	formatted := d.Format("2006-01-02")
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}

// User represents a user in the system with validation tags.
type User struct {
	UserId         string          `json:"userId" dynamodbav:"userId"`
	ShopID         string          `json:"shopId" dynamodbav:"shopId" validate:"required"`
	Email          string          `json:"email" dynamodbav:"email" validate:"required,email"`
	Role           string          `json:"role" dynamodbav:"role" validate:"required"`
	Firstname      string          `json:"firstname" dynamodbav:"firstname" validate:"required"`
	Lastname       string          `json:"lastname" dynamodbav:"lastname" validate:"required"`
	Gender         string          `json:"gender" dynamodbav:"gender" validate:"required"`
	Birthdate      Date            `json:"birthdate" dynamodbav:"birthdate"`
	Addresses      []Address       `json:"addresses" dynamodbav:"addresses"`
	Phones         []Phone         `json:"phones" dynamodbav:"phones"`
	PaymentMethods []PaymentMethod `json:"paymentMethods" dynamodbav:"paymentMethods"`
	Image          Image           `json:"image" dynamodbav:"image"`
	CreatedAt      time.Time       `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt" dynamodbav:"updatedAt"`
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

// Helper functions to parse complex fields
func ParseAddresses(input []interface{}) ([]Address, error) {
	addresses := make([]Address, 0, len(input))
	for _, addr := range input {
		addrMap, ok := addr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid address format")
		}
		var address Address
		address.Type, _ = addrMap["type"].(string)
		address.Address, _ = addrMap["address"].(string)
		address.Subdistrict, _ = addrMap["subdistrict"].(string)
		address.District, _ = addrMap["district"].(string)
		address.City, _ = addrMap["city"].(string)
		address.Province, _ = addrMap["province"].(string)
		address.Country, _ = addrMap["country"].(string)
		if postalCode, ok := addrMap["postalCode"].(float64); ok {
			address.PostalCode = int(postalCode)
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func ParsePhones(input []interface{}) ([]Phone, error) {
	phones := make([]Phone, 0, len(input))
	for _, ph := range input {
		phMap, ok := ph.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid phone format")
		}
		var phone Phone
		phone.Type, _ = phMap["type"].(string)
		phone.Code, _ = phMap["code"].(string)
		phone.Number, _ = phMap["number"].(string)
		phones = append(phones, phone)
	}
	return phones, nil
}

func ParsePaymentMethods(input []interface{}) ([]PaymentMethod, error) {
	paymentMethods := make([]PaymentMethod, 0, len(input))
	for _, pm := range input {
		pmMap, ok := pm.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid payment method format")
		}
		var paymentMethod PaymentMethod
		paymentMethod.Type, _ = pmMap["type"].(string)
		paymentMethod.Name, _ = pmMap["name"].(string)
		paymentMethod.Number, _ = pmMap["number"].(string)
		if expiryMonth, ok := pmMap["expiryMonth"].(float64); ok {
			paymentMethod.ExpiryMonth = int(expiryMonth)
		}
		if expiryYear, ok := pmMap["expiryYear"].(float64); ok {
			paymentMethod.ExpiryYear = int(expiryYear)
		}
		paymentMethods = append(paymentMethods, paymentMethod)
	}
	return paymentMethods, nil
}

func ParseImage(input map[string]interface{}) (Image, error) {
	var img Image
	img.Name, _ = input["name"].(string)
	img.URL, _ = input["url"].(string)
	img.Desc, _ = input["desc"].(string)
	return img, nil
}
