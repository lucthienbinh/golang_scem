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
	UseShortShip      bool   `json:"use_short_ship"`
	ShortShipID       uint   `json:"short_ship_id"`
	UseLongShip       bool   `json:"use_long_ship"`
	LongShipID        uint   `json:"long_ship_id"`
	TotalPrice        int64  `json:"total_price"`
	CreatedAt         int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         int64  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt
}

// OrderPay structure
type OrderPay struct {
	ID                 uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID            uint   `json:"order_id"`
	PayMethod          string `json:"pay_method"`
	PayServiceProvider string `json:"pay_service_provider"`
	PayStatus          bool   `json:"pay_status"`
	PayEmployeeID      uint   `json:"pay_employee_id"`
	TotalPrice         int64  `json:"total_price"`
	CreatedAt          int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// TransportType structure
type TransportType struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	Name                string `json:"name"`
	SameCity            bool   `json:"same_city"`
	BusStationFrom      string `json:"bus_station_from"`
	BusStationTo        string `json:"bus_station_to"`
	LongShipPrice       int64  `json:"long_ship_price"`
	ShortShipPricePerKm int64  `json:"short_ship_price_per_km"`
}

// OrderWorkflowData structure
type OrderWorkflowData struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	WorkflowKey         uint   `json:"workflow_key"`
	WorkflowInstanceKey uint   `json:"workflow_instance_key"`
	OrderID             uint   `json:"order_id"`
	CustomerReceiveID   uint   `json:"customer_receive_id"`
	PayMethod           string `json:"pay_method"`
	UseShortShip        bool   `json:"use_short_ship"`
	ShortShipID         uint   `json:"short_ship_id"`
	UseLongShip         bool   `json:"use_long_ship"`
	LongShipID          uint   `json:"long_ship_id"`
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

// OrderWorkflowCreate structure
type OrderWorkflowCreate struct {
	OrderID      uint   `json:"order_id" validate:"nonzero"`
	PayMethod    string `json:"pay_method" validate:"nonzero"`
	TotalPrice   int64  `json:"total_price" validate:"nonzero"`
	UseShortShip bool   `json:"use_short_ship"`
	UseLongShip  bool   `json:"use_long_ship"`
}
