package model

import (
	"gorm.io/gorm"
)

// -------------------- Table in database --------------------

// UserAuthenticate structure for authentication ONLY
type UserAuthenticate struct {
	gorm.Model `json:"-"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	EmployeeID uint   `json:"emloyee_id"`
	CustomerID uint   `json:"customer_id"`
	Active     bool   `gorm:"default:1" json:"active"`
}

// UserFCMToken structure for Firebase Cloud Messaging
type UserFCMToken struct {
	ID         uint   `gorm:"primary_key;<-:false" json:"id"`
	UserAuthID uint   `json:"user_auth_id"`
	CustomerID uint   `json:"customer_id"`
	EmployeeID uint   `json:"employee_id"`
	Token      string `json:"token"`
}

// Customer structure
type Customer struct {
	ID         uint   `gorm:"primary_key;<-:false" json:"id"`
	UserAuthID uint   `gorm:"<-:create" json:"-"`
	Name       string `json:"name" validate:"nonzero"`
	Age        uint16 `json:"age" validate:"nonzero"`
	Phone      int64  `json:"phone" validate:"nonzero"`
	Gender     string `json:"gender" validate:"nonzero"`
	Address    string `json:"address" validate:"nonzero"`
	Point      int16  `json:"point"`
	DeletedAt  gorm.DeletedAt
}

// CustomerCredit structure
type CustomerCredit struct {
	ID             uint  `gorm:"primary_key;<-:false" json:"id"`
	CustomerID     uint  `json:"customer_id"`
	Phone          int64 `json:"phone"`
	ValidatePhone  bool  `json:"validate_phone"`
	AccountBalance int64 `json:"account_balance"`
	UpdatedAt      int64 `gorm:"autoUpdateTime" json:"updated_at"`
}

// CustomerNotification structure
type CustomerNotification struct {
	ID         uint   `gorm:"primary_key;<-:false" json:"id"`
	CustomerID uint   `json:"customer_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  int64  `gorm:"autoCreateTime" json:"created_at"`
}

// Employee structure
type Employee struct {
	ID                 uint   `gorm:"primary_key;<-:false" json:"id"`
	UserAuthID         uint   `gorm:"<-:create" json:"-"`
	Name               string `json:"name" validate:"nonzero"`
	Age                uint16 `json:"age" validate:"nonzero"`
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
	Name string `json:"name" validate:"nonzero"`
}

// DeliveryLocation structure
type DeliveryLocation struct {
	ID       uint   `gorm:"primary_key;<-:false" json:"id"`
	City     string `json:"city" validate:"nonzero"`
	District string `json:"district" validate:"nonzero"`
}

// -------------------- Struct uses to fetch data from database --------------------

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
	Age      uint16 `json:"age" validate:"nonzero"`
	Phone    int64  `json:"phone" validate:"nonzero"`
	Gender   string `json:"gender" validate:"nonzero"`
	Address  string `json:"address" validate:"nonzero"`
}

// EmployeeWithAuth structure
type EmployeeWithAuth struct {
	Email              string `json:"email" validate:"nonzero"`
	Password           string `json:"password" validate:"nonzero"`
	Name               string `json:"name" validate:"nonzero"`
	Age                uint16 `json:"age" validate:"nonzero"`
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
	Age                      uint16 `json:"age"`
	Phone                    int64  `json:"phone"`
	Gender                   string `json:"gender"`
	Address                  string `json:"address"`
	IdentityCard             string `json:"identity_card"`
	EmployeeTypeID           uint   `json:"employee_type_id"`
	EmployeeTypeName         string `json:"employee_type_name"`
	Avatar                   string `json:"avatar"`
	DeliveryLocationCity     string `json:"delivery_location_city"`
	DeliveryLocationDistrict string `json:"delivery_location_district"`
}

// EmployeeInfoForShortShip structure
type EmployeeInfoForShortShip struct {
	ID                 uint `json:"id"`
	EmployeeTypeID     uint `json:"employee_type_id" validate:"nonzero"`
	DeliveryLocationID uint `json:"delivery_location_id"`
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
