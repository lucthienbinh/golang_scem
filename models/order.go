package models

import (
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

// -------------------- Table in database --------------------

// OrderInfo structure
type OrderInfo struct {
	ID                uint   `gorm:"primary_key" json:"id"`
	Weight            int16  `json:"weight"`
	Volume            int16  `json:"volume"`
	Type              string `json:"type"`
	Image             string `json:"image"`
	CustomerSendID    uint   `json:"customer_send_id" `
	CustomerReceiveID uint   `json:"customer_receive_id"`
	TrasnportTypeID   uint   `json:"trasnport_type_id" binding:"required"`
	HasPackage        bool   `json:"has_package"`
	EmployeeID        uint   `json:"employee_id" binding:"required"`
	Receiver          string `json:"receiver" binding:"required"`
	Detail            string `json:"detail" binding:"required"`
	TotalPrice        int32  `json:"total_price" binding:"required"`
	Note              string `json:"note"`
	CreatedAt         int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderStatusJSON save data of zeebe client
type OrderStatusJSON struct {
	gorm.Model
	Name       string
	Attributes datatypes.JSON
}

// TransportType structure
type TransportType struct {
	ID              uint   `gorm:"primary_key" json:"id"`
	Name            string `json:"name"`
	RouteFixedPrice int32  `json:"fixed_price"`
	PricePerKm      int32  `json:"price_per_km"`
}

// -------------------- Struct use to covert data to json for handler --------------------

// OrderInfoDatabase structure to return data to frontend
type OrderInfoDatabase struct {
	ID                  uint   `json:"id"`
	Weight              int16  `json:"weight"`
	Volume              int16  `json:"volume"`
	Type                string `json:"type"`
	Image               string `json:"image"`
	HasPackage          bool   `json:"has_package"`
	CustomerSendName    string `json:"customer_send_name"`
	CustomerReceiveName string `json:"customer_receive_name"`
	TrasnportType       string `json:"trasnport_type"`
	EmployeeName        string `json:"employee_name"`
	Receiver            string `json:"receiver"`
	Detail              string `json:"detail"`
	TotalPrice          int32  `json:"total_price"`
	Note                string `json:"note"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}

// -------------------- Convert function --------------------
