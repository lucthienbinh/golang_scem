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
	HasPackage        bool   `json:"has_package"`
	CustomerSendID    uint   `json:"customer_send_id" validate:"nonzero"`
	CustomerReceiveID uint   `json:"customer_receive_id"`
	EmplCreateID      uint   `json:"empl_create_id"`
	EmplShipID        uint   `json:"empl_ship_id"`
	Receiver          string `json:"receiver" validate:"nonzero"`
	TrasnportTypeID   uint   `json:"trasnport_type_id" validate:"nonzero"`
	Detail            string `json:"detail" validate:"nonzero"`
	TotalPrice        int64  `json:"total_price" validate:"nonzero"`
	Note              string `json:"note"`
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

// OrderShip structure
type OrderShip struct {
	ID           uint `gorm:"primary_key;<-:false" json:"id"`
	OrderID      uint `json:"order_id"`
	UseShortShip bool `json:"use_short_ship"`
	ShortShipID  uint `json:"short_ship_id"`
	UseLongShip  bool `json:"use_long_ship"`
	LongShipID   uint `json:"long_ship_id"`
}

// TransportType structure
type TransportType struct {
	ID              uint   `gorm:"primary_key;<-:false" json:"id"`
	Name            string `json:"name"`
	SameCity        bool   `json:"same_city"`
	RouteFixedPrice int64  `json:"fixed_price"`
	PricePerKm      int64  `json:"price_per_km"`
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
	HasPackage          bool   `json:"has_package"`
	CustomerSendName    string `json:"customer_send_name"`
	CustomerReceiveName string `json:"customer_receive_name"`
	EmplCreateName      string `json:"empl_create_name"`
	EmplShipName        string `json:"empl_ship_name"`
	Receiver            string `json:"receiver"`
	TrasnportType       string `json:"trasnport_type"`
	Detail              string `json:"detail"`
	TotalPrice          int64  `json:"total_price"`
	Note                string `json:"note"`
	CreatedAt           int64  `json:"created_at"`
	UpdatedAt           int64  `json:"updated_at"`
}
