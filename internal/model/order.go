package model

import (
	"gorm.io/gorm"
)

// -------------------- Table in database --------------------

// OrderInfo structure
type OrderInfo struct {
	ID                uint   `gorm:"primary_key;<-:false" json:"id"`
	Weight            int16  `json:"weight"`
	Volume            int16  `json:"volume"`
	Type              string `json:"type"`
	Image             string `json:"image"`
	CustomerSendID    uint   `json:"customer_send_id" validate:"nonzero"`
	CustomerReceiveID uint   `json:"customer_receive_id"`
	EmplCreateID      uint   `json:"empl_create_id"`
	EmplShipID        uint   `json:"empl_ship_id"`
	OriginalSender    string `json:"original_sender" validate:"nonzero"`
	Sender            string `json:"sender" validate:"nonzero"`
	Receiver          string `json:"receiver" validate:"nonzero"`
	TransportTypeID   uint   `json:"transport_type_id" validate:"nonzero"`
	Detail            string `json:"detail" validate:"nonzero"`
	Note              string `json:"note"`
	UseLongShip       bool   `json:"use_long_ship"`
	UseShortShip      bool   `json:"use_short_ship"`
	LongShipID        uint   `json:"long_ship_id"`
	OrderLongShipID   uint   `json:"order_long_ship_id"`
	OrderShortShipID  uint   `json:"order_short_ship_id"`
	ShortShipDistance int64  `json:"short_ship_distance"`
	TotalPrice        int64  `json:"total_price"`
	CreatedAt         int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         int64  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt
}

// OrderInfoForPayment structure
type OrderInfoForPayment struct {
	ID                uint  `json:"id"`
	CustomerSendID    uint  `json:"customer_send_id" validate:"nonzero"`
	CustomerReceiveID uint  `json:"customer_receive_id"`
	UseShortShip      bool  `json:"use_short_ship"`
	UseLongShip       bool  `json:"use_long_ship"`
	TotalPrice        int64 `json:"total_price"`
}

// TransportType structure
type TransportType struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	Name                string `json:"name" validate:"nonzero"`
	SameCity            bool   `json:"same_city"`
	BusStationFrom      string `json:"bus_station_from"`
	BusStationTo        string `json:"bus_station_to"`
	LongShipDuration    int64  `json:"long_ship_duration"`
	LongShipPrice       int64  `json:"long_ship_price"`
	ShortShipPricePerKm int64  `json:"short_ship_price_per_km" validate:"nonzero"`
}

// OrderPay structure
type OrderPay struct {
	ID              uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID         uint   `json:"order_id"`
	PayMethod       string `json:"pay_method"`
	PayStatus       bool   `json:"pay_status"`
	TotalPrice      int64  `json:"total_price"`
	FinishedStepOne bool   `json:"finished_step_one"`
	FinishedStepTwo bool   `json:"finished_step_two"`
	ConfirmString   string `json:"confirm_string"`
	// If customer pay method is cash, we will need an employee to confirm
	PayEmployeeID uint `json:"pay_employee_id"`
	// If customer pay method is credit, we will need customer confirmation
	PayCustomerID uint  `json:"pay_customer_id"`
	CreatedAt     int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     int64 `gorm:"autoUpdateTime" json:"updated_at"`
}

// -------------------- Struct uses to fetch data for frontend --------------------

// OrderInfoFetchDB structure
type OrderInfoFetchDB struct {
	ID                  uint   `json:"id"`
	Weight              int16  `json:"weight"`
	Volume              int16  `json:"volume"`
	Type                string `json:"type"`
	Image               string `json:"image"`
	CustomerSendName    string `json:"customer_send_name"`
	CustomerReceiveName string `json:"customer_receive_name"`
	EmplCreateName      string `json:"empl_create_name"`
	EmplShipName        string `json:"empl_ship_name"`
	Receiver            string `json:"receiver"`
	TransportType       string `json:"transport_type"`
	Detail              string `json:"detail"`
	TotalPrice          int64  `json:"total_price"`
	Note                string `json:"note"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}
