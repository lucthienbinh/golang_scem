package model

// -------------------- Table in database --------------------

// OrderLongShip structure
type OrderLongShip struct {
	ID              uint `gorm:"primary_key;<-:false" json:"id"`
	OrderID         uint `json:"order_id"`
	TrasnportTypeID uint `json:"trasnport_type_id"`
	LastLocationID  uint `json:"last_location_id"`
	EmplLoadID      uint `json:"empl_load_id"`
	EmplUnloadID    uint `json:"empl_unload_id"`
	EmplDriverID    uint `json:"empl_driveer_id"`
	Finished        bool `gorm:"default:0" json:"finished"`
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
