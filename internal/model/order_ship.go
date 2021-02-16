package model

// -------------------- Table in database --------------------

// OrderWorkflowData structure
type OrderWorkflowData struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	WorkflowKey         string `json:"workflow_key"`
	WorkflowInstanceKey uint   `json:"workflow_instance_key"`
	// Mapping ID for data
	OrderID          uint `json:"order_id"`
	OrderPayID       uint `json:"order_pay_id"`
	LongShipID       uint `json:"long_ship_id"`
	OrderLongShipID  uint `json:"order_long_ship_id"`
	OrderShortShipID uint `json:"order_short_ship_id"`
	// Variable use for Zeebe gateway
	ShipperReceiveMoney bool `json:"shipper_receive_money"`
	UseLongShip         bool `json:"use_long_ship"`
	CustomerSendID      uint `json:"customer_send_id"`
	CustomerReceiveID   uint `json:"customer_receive_id"`
}

// LongShipWorkflowData structure
type LongShipWorkflowData struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	WorkflowKey         string `json:"workflow_key"`
	WorkflowInstanceKey uint   `json:"workflow_instance_key"`
	LongShipID          uint   `json:"long_ship_id"`
}

// LongShip structure
type LongShip struct {
	ID                       uint   `gorm:"<-:create" json:"id"`
	TransportTypeID          uint   `json:"transport_type_id" validate:"nonzero"`
	LicensePlate             string `json:"license_plate" validate:"nonzero"`
	EstimatedTimeOfDeparture int64  `json:"estimated_time_of_departure" validate:"nonzero"`
	EstimatedTimeOfArrival   int64  `json:"estimated_time_of_arrival" validate:"nonzero"`
	CurrentLocation          string `json:"current_location"`
	// Message data in workflow - Start
	// Package Loaded
	PackageLoaded bool  `json:"package_loaded"`
	EmplLoadID    uint  `json:"empl_load_id"`
	LoadedTime    int64 `json:"loaded_time"`
	// Vehicle Started
	VehicleStarted bool  `json:"vehicle_started"`
	EmplDriver1ID  uint  `json:"empl_driver_1_id"`
	StartedTime    int64 `json:"started_time"`
	// Vehicle Arrived
	VehicleArrived bool  `json:"vehicle_arrived"`
	EmplDriver2ID  uint  `json:"empl_driver_2_id"`
	ArrivedTime    int64 `json:"arrived_time"`
	// Package Unloaded
	PackageUnloaded bool  `json:"package_unloaded"`
	EmplUnloadID    uint  `json:"empl_unload_id"`
	UnloadedTime    int64 `json:"unloaded_time"`
	// Message data in workflow - End
	Finished  bool   `json:"finished"`
	LSQrCode  string `json:"ls_qr_code"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderLongShip structure
type OrderLongShip struct {
	ID                uint `gorm:"primary_key;<-:false" json:"id"`
	OrderID           uint `json:"order_id"`
	LongShipID        uint `json:"long_ship_id"`
	CustomerSendID    uint `json:"customer_send_id"`
	CustomerReceiveID uint `json:"customer_receive_id"`
}

// OrderShortShip structure
type OrderShortShip struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID             uint   `json:"order_id" validate:"nonzero"`
	ShipperID           uint   `json:"shipper_id" validate:"nonzero"`
	CustomerSendID      uint   `json:"customer_send_id"`
	CustomerReceiveID   uint   `json:"customer_receive_id"`
	Sender              string `json:"sender"`
	Receiver            string `json:"receiver"`
	ShipperReceiveMoney bool   `json:"shipper_receive_money"`
	// Message data in workflow - Start
	// Shipper Called
	ShipperCalled bool `json:"shipper_called"`
	// Shipper Received Money
	ShipperReceivedMoney bool  `json:"shipper_received_money"`
	ReceivedMoneyTime    int64 `json:"received_money_time"`
	// Shipper Shipped
	ShipperShipped bool  `json:"shipper_shipped"`
	ShippedTime    int64 `json:"shipped_time"`
	// Shipper Confirmed
	ShipperConfirmed     string `json:"shipper_confirmed"`
	ShipperConfirmedTime int64  `json:"shipper_confirmed_time"`
	// Message data in workflow - End
	Canceled       bool   `json:"canceled"`
	CanceledReason string `json:"canceled_reason"`
	Finished       bool   `json:"finished"`
	OSSQrCode      string `json:"oss_qr_code"`
	CreatedAt      int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      int64  `gorm:"autoUpdateTime" json:"updated_at"`
}
