package handler

import (
	"os"

	"github.com/lucthienbinh/golang_scem/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// ConnectPostgres to open connect with database
func ConnectPostgres() (err error) {
	// https://github.com/go-gorm/postgres
	// https://stackoverflow.com/questions/50085286/postgresql-fatal-ident-authentication-failed-for-user
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("POSTGRES_DSN"),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	return err
}

// ConnectMySQL to open connect with database
func ConnectMySQL() (err error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := os.Getenv("MYSQL_DSN")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

// MigrationDatabase when update database
func MigrationDatabase() (err error) {
	return db.AutoMigrate(
		&model.Employee{},
		&model.Customer{},
		&model.EmployeeType{},
		&model.DeliveryLocation{},
		&model.OrderInfo{},
		&model.TransportType{},
		&model.UserAuthenticate{},
		&model.EmployeeFCMToken{},
		&model.CustomerFCMToken{},
		&model.OrderPay{},
		&model.OrderShip{},
		&model.OrderLongShip{},
		&model.OrderShortShip{},
	)
}

// RefreshDatabase remove all table and create new data
func RefreshDatabase() (err error) {
	if err := deleteDatabase(); err != nil {
		return err
	}
	if err := MigrationDatabase(); err != nil {
		return err
	}
	if err := createDeliveryLocation(); err != nil {
		return err
	}
	if err := createEmployeeType(); err != nil {
		return err
	}
	if err := createDefaultEmployee(); err != nil {
		return err
	}
	if err := createDefaultCustomer(); err != nil {
		return err
	}
	if err := createTransportType(); err != nil {
		return err
	}
	if err := createExampleOrder(); err != nil {
		return err
	}
	if err := createExampleOrderPay(); err != nil {
		return err
	}
	if err := createExampleOrderShip(); err != nil {
		return err
	}
	if err := createExampleOrderShortShip(); err != nil {
		return err
	}
	return
}

func deleteDatabase() (err error) {
	return db.Migrator().DropTable(
		&model.Employee{},
		&model.Customer{},
		&model.EmployeeType{},
		&model.DeliveryLocation{},
		&model.OrderInfo{},
		&model.TransportType{},
		&model.UserAuthenticate{},
		&model.EmployeeFCMToken{},
		&model.CustomerFCMToken{},
		&model.OrderPay{},
		&model.OrderShip{},
		&model.OrderLongShip{},
		&model.OrderShortShip{},
	)
}

func createEmployeeType() (err error) {
	employeeType := model.EmployeeType{Name: "Admin"}
	if err := db.Create(&employeeType).Error; err != nil {
		return err
	}
	employeeType = model.EmployeeType{Name: "Input staff"}
	if err := db.Create(&employeeType).Error; err != nil {
		return err
	}
	employeeType = model.EmployeeType{Name: "Delivery staff"}
	if err := db.Create(&employeeType).Error; err != nil {
		return err
	}
	employeeType = model.EmployeeType{Name: "Load package staff"}
	if err := db.Create(&employeeType).Error; err != nil {
		return err
	}
	return
}

func createDefaultEmployee() (err error) {
	userAuth := model.UserAuthenticate{Email: "admin@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	employee := model.Employee{UserAuthID: userAuth.ID, Name: "Binh", Age: 40, Phone: 223334444, Gender: "male", Address: "12 Tran Hung Dao", IdentityCard: "17687t562765786", EmployeeTypeID: 1, Avatar: "image1.jpg"}
	if err := db.Create(&employee).Error; err != nil {
		return err
	}
	userAuth = model.UserAuthenticate{Email: "inputstaff@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	employee = model.Employee{UserAuthID: userAuth.ID, Name: "Huan", Age: 35, Phone: 448883333, Gender: "male", Address: "21 Huynh Thuc Khang", IdentityCard: "17687t562765786", EmployeeTypeID: 2, Avatar: "image2.jpg"}
	if err := db.Create(&employee).Error; err != nil {
		return err
	}
	userAuth = model.UserAuthenticate{Email: "deliverystaff@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	employee = model.Employee{UserAuthID: userAuth.ID, Name: "Tuan", Age: 37, Phone: 776664993, Gender: "male", Address: "21 Nhat Tao", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 4}
	if err := db.Create(&employee).Error; err != nil {
		return err
	}
	userAuth = model.UserAuthenticate{Email: "loadpackagestaff@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	employee = model.Employee{UserAuthID: userAuth.ID, Name: "Hung", Age: 47, Phone: 776334958, Gender: "male", Address: "84 Nguyen Trau", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg"}
	if err := db.Create(&employee).Error; err != nil {
		return err
	}
	return
}

func createDefaultCustomer() (err error) {
	userAuth := model.UserAuthenticate{Email: "customer@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	customer := model.Customer{UserAuthID: userAuth.ID, Name: "Customer One", Age: 18, Phone: 223334444, Gender: "male", Address: "12 Tran Hung Dao, Phuong 1, Quan 5"}
	if err := db.Create(&customer).Error; err != nil {
		return err
	}

	userAuth = model.UserAuthenticate{Email: "customer3@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	customer = model.Customer{UserAuthID: userAuth.ID, Name: "Customer Three", Age: 18, Phone: 223334444, Gender: "female", Address: "13 Tran Hung Dao"}
	if err := db.Create(&customer).Error; err != nil {
		return err
	}

	userAuth = model.UserAuthenticate{Email: "customer2@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	customer = model.Customer{UserAuthID: userAuth.ID, Name: "Customer Two", Age: 18, Phone: 223334444, Gender: "male", Address: "14 Tran Hung Dao"}
	return db.Create(&customer).Error
}

func createTransportType() (err error) {
	transportType := model.TransportType{Name: "HCM", SameCity: true, RouteFixedPrice: 0, PricePerKm: 30000}
	if err := db.Create(&transportType).Error; err != nil {
		return err
	}
	transportType = model.TransportType{Name: "HCM-DL", RouteFixedPrice: 100000, PricePerKm: 40000}
	if err := db.Create(&transportType).Error; err != nil {
		return err
	}
	transportType = model.TransportType{Name: "HCM-VT", RouteFixedPrice: 120000, PricePerKm: 25000}
	if err := db.Create(&transportType).Error; err != nil {
		return err
	}
	transportType = model.TransportType{Name: "HCM-CT", RouteFixedPrice: 70000, PricePerKm: 30000}
	if err := db.Create(&transportType).Error; err != nil {
		return err
	}
	return
}

func createDeliveryLocation() (err error) {
	location := model.DeliveryLocation{City: "HCM", District: "1"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "HCM", District: "2"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "HCM", District: "3"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "HCM", District: "Tan Binh"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "HCM", District: "Phu Nhuan"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "DL", District: "1"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "DL", District: "2"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "DL", District: "3"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "DL", District: "4"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	location = model.DeliveryLocation{City: "DL", District: "5"}
	if err := db.Create(&location).Error; err != nil {
		return err
	}
	return
}

func createExampleOrder() (err error) {
	orderInfo := model.OrderInfo{
		Weight: 2, Volume: 10, Type: "Normal", Image: "order1.png",
		CustomerSendID: 1, TrasnportTypeID: 2, HasPackage: true,
		EmplShipID: 3, EmplCreateID: 2, Receiver: "253 Tran Hung Dao, Quan 1, SDT 23114321412",
		Detail:     "May vi tinh ca nhan va ban phim may tinh",
		TotalPrice: 200000, Note: "Giao hang vao buoi sang",
	}
	if err := db.Create(&orderInfo).Error; err != nil {
		return err
	}
	return
}

func createExampleOrderPay() (err error) {
	orderPay := model.OrderPay{
		OrderID: 1, PayStatus: true, PayEmployeeID: 2, TotalPrice: 200000, PayMethod: "cash",
	}
	if err := db.Create(&orderPay).Error; err != nil {
		return err
	}
	return
}

func createExampleOrderShip() (err error) {
	orderShip := model.OrderShip{
		OrderID: 1, UseShortShip: true, ShortShipID: 1,
	}
	if err := db.Create(&orderShip).Error; err != nil {
		return err
	}
	return
}

func createExampleOrderShortShip() (err error) {
	orderShortShip := model.OrderShortShip{
		OrderID: 1, ShipperID: 3, ShipperReceived: true,
		ShipperCalled: true, ShipperShipped: true,
		TimeConfirmed: 1608743635, CusReceiveConfirmed: true,
		Finished: true,
	}
	if err := db.Create(&orderShortShip).Error; err != nil {
		return err
	}
	return
}
