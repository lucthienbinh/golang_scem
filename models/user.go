package models

import (
	"github.com/jinzhu/gorm"
)

// UserAuthenticate structure for authentication
type UserAuthenticate struct {
	gorm.Model `json:"-"`
	Email      string `json:"email" binding:"email"`
	Password   string `json:"password"`
	Active     bool   `gorm:"default:1" json:"active"`
}

// -------------------- Table in database --------------------

// Customer structure
type Customer struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	UserAuthID uint   `json:"-"`
	Name       string `json:"name" binding:"required"`
	Age        uint8  `json:"age" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Gender     string `json:"gender" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Point      int16  `json:"point"`
}

// Employee structure
type Employee struct {
	ID                 uint   `gorm:"primary_key" json:"id"`
	UserAuthID         uint   `json:"user_auth_id" json:"-"`
	Name               string `json:"name" binding:"required"`
	Age                uint8  `json:"age" binding:"required"`
	Phone              string `json:"phone" binding:"required"`
	Gender             string `json:"gender" binding:"required"`
	Address            string `json:"address" binding:"required"`
	IdentityCard       string `json:"indentity_card" binding:"required"`
	EmployeeTypeID     uint   `json:"employee_type_id" binding:"required"`
	Avatar             string `json:"avatar" binding:"required"`
	DeliveryLocationID uint   `json:"delivery_location_id"`
}

// EmployeeType structure
type EmployeeType struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `binding:"required"`
}

// DeliveryLocation structure
type DeliveryLocation struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	City     string `json:"city"`
	District string `json:"district"`
}

// -------------------- Struct use to covert data to json for handler --------------------

// Login structure
type Login struct {
	Email    string `form:"email" json:"email" xml:"email"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// EmployeeBasicInfo structure
type EmployeeBasicInfo struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	EmployeeTypeID uint   `json:"employee_type_id"`
	Phone          string `json:"phone"`
	Address        string `json:"location"`
}

// CustomerBasicInfo structure
type CustomerBasicInfo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"location"`
}

// CustomerWithAuth structure
type CustomerWithAuth struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Age      uint8  `json:"age" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
	Address  string `json:"address" binding:"required"`
}

// EmployeeWithAuth structure
type EmployeeWithAuth struct {
	Email              string `json:"email" binding:"required"`
	Password           string `json:"password" binding:"required"`
	Name               string `json:"name" binding:"required"`
	Age                uint8  `json:"age" binding:"required"`
	Phone              string `json:"phone" binding:"required"`
	Gender             string `json:"gender" binding:"required"`
	Address            string `json:"address" binding:"required"`
	IdentityCard       string `json:"indentity_card" binding:"required"`
	EmployeeTypeID     uint   `json:"employee_type_id" binding:"required"`
	Avatar             string `json:"avatar" binding:"required"`
	DeliveryLocationID uint   `json:"delivery_location_id"`
}

// -------------------- Convert function --------------------

// ConvertToBasic to keep safe for sensitive info
func (c *Customer) ConvertToBasic() CustomerBasicInfo {
	return CustomerBasicInfo{
		ID:      c.ID,
		Name:    c.Name,
		Phone:   c.Phone,
		Address: c.Address,
	}
}

// ConvertCWAToNormal to fetch data from front end
func (c *CustomerWithAuth) ConvertCWAToNormal() (*Customer, *UserAuthenticate) {
	return &Customer{
			Name:    c.Name,
			Age:     c.Age,
			Phone:   c.Phone,
			Gender:  c.Gender,
			Address: c.Address,
		}, &UserAuthenticate{
			Email:    c.Email,
			Password: c.Password,
		}
}

// ConvertToBasic to keep safe for sensitive info
func (e *Employee) ConvertToBasic() EmployeeBasicInfo {
	return EmployeeBasicInfo{
		ID:             e.ID,
		Name:           e.Name,
		EmployeeTypeID: e.EmployeeTypeID,
		Phone:          e.Phone,
		Address:        e.Address,
	}
}

// ConvertEWAToNormal to fetch data from front end
func (c *EmployeeWithAuth) ConvertEWAToNormal() (*Employee, *UserAuthenticate) {
	return &Employee{
			Name:               c.Name,
			Age:                c.Age,
			Phone:              c.Phone,
			Gender:             c.Gender,
			Address:            c.Address,
			IdentityCard:       c.IdentityCard,
			Avatar:             c.Avatar,
			DeliveryLocationID: c.DeliveryLocationID,
		}, &UserAuthenticate{
			Email:    c.Email,
			Password: c.Password,
		}
}
