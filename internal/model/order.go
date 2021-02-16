package model

import (
	"gorm.io/gorm"
)

// -------------------- Table in database --------------------

// OrderInfo structure
type OrderInfo struct {
	ID uint `gorm:"<-:create" json:"id"`
	// Package information
	Weight int64  `json:"weight"`
	Volume int64  `json:"volume"`
	Type   string `json:"type"`
	Image  string `json:"image"`
	// User information
	CustomerSendID    uint `json:"customer_send_id" validate:"nonzero"`
	CustomerReceiveID uint `json:"customer_receive_id"`
	EmplCreateID      uint `json:"empl_create_id"`
	// Delivery information
	Sender          string `json:"sender" validate:"nonzero"`
	Receiver        string `json:"receiver" validate:"nonzero"`
	TransportTypeID uint   `json:"transport_type_id" validate:"nonzero"`
	Detail          string `json:"detail" validate:"nonzero"`
	Note            string `json:"note"`
	TotalPrice      int64  `json:"total_price"`
	// Long ship and short ship
	UseLongShip       bool  `json:"use_long_ship"`
	LongShipID        uint  `json:"long_ship_id"`
	OrderLongShipID   uint  `json:"order_long_ship_id"`
	OrderShortShipID  uint  `json:"order_short_ship_id"`
	ShortShipDistance int64 `json:"short_ship_distance"`
	// Time information
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt
}

// TransportType structure
type TransportType struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	SameCity            bool   `json:"same_city"`
	LocationOne         string `json:"location_one" validate:"nonzero"`
	LocationTwo         string `json:"location_two"`
	BusStationFrom      string `json:"bus_station_from"`
	BusStationTo        string `json:"bus_station_to"`
	LongShipDuration    int64  `json:"long_ship_duration"`
	LongShipPrice       int64  `json:"long_ship_price"`
	ShortShipPricePerKm int64  `json:"short_ship_price_per_km" validate:"nonzero"`
}

// OrderPay structure
type OrderPay struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID             uint   `json:"order_id"`
	PayMethod           string `json:"pay_method"`
	PayStatus           bool   `json:"pay_status"`
	TotalPrice          int64  `json:"total_price"`
	FinishedStepOne     bool   `json:"finished_step_one"`
	FinishedStepTwo     bool   `json:"finished_step_two"`
	ShipperReceiveMoney bool   `json:"shipper_receive_money"`
	CreatedAt           int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderVoucher structure
type OrderVoucher struct {
	ID        uint   `gorm:"primary_key;<-:false" json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	StartDate int64  `json:"start_date"`
	EndDate   int64  `json:"end_date"`
	Discount  int64  `json:"discount"`
}

// -------------------- Struct uses to fetch data from database --------------------

// OrderInfoForPayment structure (create workflow)
type OrderInfoForPayment struct {
	ID                uint
	CustomerSendID    uint
	CustomerReceiveID uint
	UseLongShip       bool
	TotalPrice        int64
	LongShipID        uint
}

// OrderInfoForShipment structure (create short ship)
type OrderInfoForShipment struct {
	ID                uint
	TransportTypeID   uint
	UseLongShip       bool
	LongShipID        uint
	CustomerSendID    uint
	CustomerReceiveID uint
	Sender            string
	Receiver          string
}

// OrderInfoWithVoucher structure
type OrderInfoWithVoucher struct {
	CustomerSendID    uint   `json:"customer_send_id" validate:"nonzero"`
	Sender            string `json:"sender" validate:"nonzero"`
	Receiver          string `json:"receiver" validate:"nonzero"`
	TransportTypeID   uint   `json:"transport_type_id" validate:"nonzero"`
	Detail            string `json:"detail" validate:"nonzero"`
	Note              string `json:"note"`
	ShortShipDistance int64  `json:"short_ship_distance"`
	Image             string `json:"image"`
	OrderVoucherID    uint   `json:"order_voucher_id"`
}

// ConvertToBasicOrder function
func (ord *OrderInfoWithVoucher) ConvertToBasicOrder() (*OrderInfo, uint) {
	return &OrderInfo{
		CustomerSendID:    ord.CustomerSendID,
		Sender:            ord.Sender,
		Receiver:          ord.Receiver,
		TransportTypeID:   ord.TransportTypeID,
		Detail:            ord.Detail,
		Note:              ord.Note,
		ShortShipDistance: ord.ShortShipDistance,
		Image:             ord.Image,
	}, ord.OrderVoucherID
}
