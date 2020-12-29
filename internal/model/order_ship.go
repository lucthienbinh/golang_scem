package model

// -------------------- Table in database --------------------

// OrderLongShip structure
type OrderLongShip struct {
	ID              uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID         uint   `json:"order_id"`
	CurrentLocation string `json:"current_location"`
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
	OLSQrCode string `json:"qls_qr_code"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderShortShip structure
type OrderShortShip struct {
	ID                  uint `gorm:"primary_key;<-:false" json:"id"`
	OrderID             uint `json:"order_id"`
	ShipperID           uint `json:"shipper_id"`
	CustomerReceiveID   uint `json:"customer_receive_id"`
	ShipperReceiveMoney bool `json:"shipper_receive_money"`
	// Message data in workflow - Start
	// Shipper Received
	ShipperReceived bool  `json:"shipper_received"`
	ReceivedTime    int64 `json:"received_time"`
	// Shipper Called
	ShipperCalled bool  `json:"shipper_called"`
	TimeConfirmed int64 `json:"time_confirmed"`
	CalledTime    int64 `json:"called_time"`
	// Shipper Received Money
	ShipperReceivedMoney bool  `json:"shipper_received_money"`
	ReceivedMoneyTime    int64 `json:"received_money_time"`
	// Shipper Shipped
	ShipperShipped bool  `json:"shipper_shipped"`
	ShippedTime    int64 `json:"shipped_time"`
	// Customer Receive Confirmed
	CusReceiveConfirmed     bool  `json:"cus_receive_confirmed"`
	CusReceiveConfirmedTime int64 `json:"cus_receive_confirmed_time"`
	// Shipper Confirmed
	ShipperConfirmed     string `json:"shipper_confirmed"`
	ShipperConfirmedTime int64  `json:"shipper_confirmed_time"`
	// Message data in workflow - End
	Canceled       bool   `json:"canceled"`
	CanceledReason string `json:"canceled_reason"`
	Finished       bool   `json:"finished"`
	OSSQrCode      string `json:"qss_qr_code"`
	UpdatedAt      int64  `gorm:"autoUpdateTime" json:"updated_at"`
}
