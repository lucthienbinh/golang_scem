package model

// -------------------- Table in database --------------------

// OrderLongShip structure
type OrderLongShip struct {
	ID              uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID         uint   `json:"order_id"`
	CurrentLocation string `json:"current_location"`
	PackageLoaded   bool   `gorm:"default:0" json:"package_loaded"`
	VehicleStarted  bool   `gorm:"default:0" json:"vehicle_started"`
	VehicleArrived  bool   `gorm:"default:0" json:"vehicle_arrived"`
	PackageUnloaded bool   `gorm:"default:0" json:"package_unloaded"`
	EmplLoadID      uint   `json:"empl_load_id"`
	EmplUnloadID    uint   `json:"empl_unload_id"`
	EmplDriverID    uint   `json:"empl_driver_id"`
	Finished        bool   `gorm:"default:0" json:"finished"`
}

// OrderShortShip structure
type OrderShortShip struct {
	ID                  uint   `gorm:"primary_key;<-:false" json:"id"`
	OrderID             uint   `json:"order_id"`
	EmplShipID          uint   `json:"empl_ship_id"`
	EmplShipReceived    bool   `json:"empl_ship_received"`
	EmplShipCalled      bool   `json:"empl_ship_called"`
	EmplShipShipped     bool   `json:"empl_ship_shipped"`
	TimeConfirmed       bool   `json:"time_confirmed"`
	CusReceiveConfirmed bool   `json:"cus_receive_confirmed"`
	ShipperConfirmed    string `json:"shipper_confirmed"`
	Finished            bool   `gorm:"default:0" json:"finished"`
}
