package model

// -------------------- Table in database --------------------

// OrderLongShip structure
type OrderLongShip struct {
	ID              uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID         uint   `json:"order_id"`
	CurrentLocation string `json:"current_location"`
	PackageLoaded   bool   `json:"package_loaded"`
	VehicleStarted  bool   `json:"vehicle_started"`
	VehicleArrived  bool   `json:"vehicle_arrived"`
	PackageUnloaded bool   `json:"package_unloaded"`
	EmplLoadID      uint   `json:"empl_load_id"`
	EmplUnloadID    uint   `json:"empl_unload_id"`
	EmplDriverID    uint   `json:"empl_driver_id"`
	Finished        bool   `json:"finished"`
	OLSQrCode       string `json:"qls_qr_code"`
	UpdatedAt       int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

// OrderShortShip structure
type OrderShortShip struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID             uint   `json:"order_id"`
	ShipperID           uint   `json:"shipper_id"`
	ShipperReceived     bool   `json:"shipper_received"`
	ShipperCalled       bool   `json:"shipper_called"`
	TimeConfirmed       int64  `json:"time_confirmed"`
	ShipperShipped      bool   `json:"shipper_shipped"`
	CusReceiveConfirmed bool   `json:"cus_receive_confirmed"`
	ShipperConfirmed    string `json:"shipper_confirmed"`
	Canceled            bool   `json:"canceled"`
	CanceleReason       string `json:"canceled_reason"`
	Finished            bool   `json:"finished"`
	OSSQrCode           string `json:"qss_qr_code"`
	UpdatedAt           int64  `gorm:"autoUpdateTime" json:"updated_at"`
}
