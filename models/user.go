package models

import (
	"gorm.io/gorm"
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
	ID         uint   `gorm:"primary_key;<-:false" json:"id"`
	UserAuthID uint   `gorm:"<-:create" json:"-"`
	Name       string `json:"name" validate:"nonzero"`
	Age        uint8  `json:"age" validate:"nonzero"`
	Phone      int64  `json:"phone" validate:"nonzero"`
	Gender     string `json:"gender" validate:"nonzero"`
	Address    string `json:"address" validate:"nonzero"`
	Point      int16  `json:"point"`
	DeletedAt  gorm.DeletedAt
}

// Employee structure
type Employee struct {
	ID                 uint   `gorm:"primary_key;<-:false" json:"id"`
	UserAuthID         uint   `gorm:"<-:create" json:"-"`
	Name               string `json:"name" validate:"nonzero"`
	Age                uint8  `json:"age" validate:"nonzero"`
	Phone              int64  `json:"phone" validate:"nonzero"`
	Gender             string `json:"gender" validate:"nonzero"`
	Address            string `json:"address" validate:"nonzero"`
	IdentityCard       string `json:"identity_card" validate:"nonzero"`
	EmployeeTypeID     uint   `json:"employee_type_id" validate:"nonzero"`
	Avatar             string `json:"avatar" validate:"nonzero"`
	DeliveryLocationID uint   `json:"delivery_location_id"`
	DeletedAt          gorm.DeletedAt
}

// EmployeeType structure
type EmployeeType struct {
	ID   uint   `gorm:"primary_key;<-:false" json:"id"`
	Name string `validate:"nonzero" json:"name"`
}

// DeliveryLocation structure
type DeliveryLocation struct {
	ID       uint   `gorm:"primary_key;<-:false" json:"id"`
	City     string `json:"city" validate:"nonzero"`
	District string `json:"district" validate:"nonzero"`
}

// -------------------- Struct uses to fetch data for frontend --------------------

// Login structure
type Login struct {
	Email    string `form:"email" json:"email" xml:"email"  validate:"nonzero"`
	Password string `form:"password" json:"password" xml:"password" validate:"nonzero"`
}

// EmployeeBasicInfo structure
type EmployeeBasicInfo struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	EmployeeTypeID uint   `json:"employee_type_id"`
	Phone          int64  `json:"phone"`
	Address        string `json:"location"`
}

// CustomerBasicInfo structure
type CustomerBasicInfo struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Phone   int64  `json:"phone"`
	Address string `json:"location"`
}

// CustomerWithAuth structure
type CustomerWithAuth struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
	Name     string `json:"name" validate:"nonzero"`
	Age      uint8  `json:"age" validate:"nonzero"`
	Phone    int64  `json:"phone" validate:"nonzero"`
	Gender   string `json:"gender" validate:"nonzero"`
	Address  string `json:"address" validate:"nonzero"`
}

// EmployeeWithAuth structure
type EmployeeWithAuth struct {
	Email              string `json:"email" validate:"nonzero"`
	Password           string `json:"password" validate:"nonzero"`
	Name               string `json:"name" validate:"nonzero"`
	Age                uint8  `json:"age" validate:"nonzero"`
	Phone              int64  `json:"phone" validate:"nonzero"`
	Gender             string `json:"gender" validate:"nonzero"`
	Address            string `json:"address" validate:"nonzero"`
	IdentityCard       string `json:"identity_card" validate:"nonzero"`
	EmployeeTypeID     uint   `json:"employee_type_id" validate:"nonzero"`
	Avatar             string `json:"avatar" validate:"nonzero"`
	DeliveryLocationID uint   `json:"delivery_location_id"`
}

// EmployeeInfoFetchDB structure
type EmployeeInfoFetchDB struct {
	ID                       uint   `json:"id"`
	Name                     string `json:"name"`
	Age                      uint8  `json:"age"`
	Phone                    int64  `json:"phone"`
	Gender                   string `json:"gender"`
	Address                  string `json:"address"`
	IdentityCard             string `json:"identity_card"`
	EmployeeTypeName         string `json:"employee_type_name"`
	Avatar                   string `json:"avatar"`
	DeliveryLocationCity     string `json:"delivery_location_city"`
	DeliveryLocationDistrict string `json:"delivery_location_district"`
}

// SelectStuct for select options
type SelectStuct struct {
	Value uint   `json:"value"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

// -------------------- Convert function to keep safe sensitive info--------------------

// ConvertToBasic function
func (c *Customer) ConvertToBasic() CustomerBasicInfo {
	return CustomerBasicInfo{
		ID:      c.ID,
		Name:    c.Name,
		Phone:   c.Phone,
		Address: c.Address,
	}
}

// ConvertToBasic function
func (e *Employee) ConvertToBasic() EmployeeBasicInfo {
	return EmployeeBasicInfo{
		ID:             e.ID,
		Name:           e.Name,
		EmployeeTypeID: e.EmployeeTypeID,
		Phone:          e.Phone,
		Address:        e.Address,
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
			EmployeeTypeID:     c.EmployeeTypeID,
			DeliveryLocationID: c.DeliveryLocationID,
		}, &UserAuthenticate{
			Email:    c.Email,
			Password: c.Password,
		}
}
