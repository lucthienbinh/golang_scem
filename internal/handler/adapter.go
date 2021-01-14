package handler

import (
	"os"

	"github.com/lucthienbinh/golang_scem/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
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

// ConnectSQLite to open connect with database
func ConnectSQLite() (err error) {
	// github.com/mattn/go-sqlite3
	dsn := os.Getenv("SQLITE_DSN")
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	return err
}

// MigrationDatabase when update database
func MigrationDatabase() (err error) {
	return db.AutoMigrate(
		&model.Employee{},
		&model.Customer{},
		&model.CustomerCredit{},
		&model.EmployeeType{},
		&model.DeliveryLocation{},
		&model.OrderInfo{},
		&model.TransportType{},
		&model.UserAuthenticate{},
		&model.UserFCMToken{},
		&model.OrderPay{},
		&model.LongShip{},
		&model.OrderLongShip{},
		&model.OrderShortShip{},
		&model.OrderWorkflowData{},
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
	if err := createCustomerCredit(); err != nil {
		return err
	}
	if err := createTransportType(); err != nil {
		return err
	}
	if err := createLongShip(); err != nil {
		return err
	}
	if err := createExampleOrder(); err != nil {
		return err
	}
	if err := createExampleOrderPay(); err != nil {
		return err
	}
	if err := createExampleOrderShortShip(); err != nil {
		return err
	}
	if err := createExampleOrder2(); err != nil {
		return err
	}
	if err := createCustomerFCMToken(); err != nil {
		return err
	}
	return
}

func deleteDatabase() (err error) {
	return db.Migrator().DropTable(
		&model.Employee{},
		&model.Customer{},
		&model.CustomerCredit{},
		&model.EmployeeType{},
		&model.DeliveryLocation{},
		&model.OrderInfo{},
		&model.TransportType{},
		&model.UserAuthenticate{},
		&model.UserFCMToken{},
		&model.OrderPay{},
		&model.LongShip{},
		&model.OrderLongShip{},
		&model.OrderShortShip{},
		&model.OrderWorkflowData{},
	)
}

func createEmployeeType() error {
	employeeType := &model.EmployeeType{Name: "Admin"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Input staff"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Delivery staff"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Load package staff"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	employeeType = &model.EmployeeType{Name: "Load package staff1"}
	if err := db.Create(employeeType).Error; err != nil {
		return err
	}
	return nil
}

func createDefaultEmployee() error {
	userAuth := &model.UserAuthenticate{Email: "admin@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee := &model.Employee{UserAuthID: userAuth.ID, Name: "Binh", Age: 40, Phone: 223334444, Gender: "male", Address: "12 Tran Hung Dao", IdentityCard: "17687t562765786", EmployeeTypeID: 1, Avatar: "image1.jpg"}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "inputstaff@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Huan", Age: 35, Phone: 448883333, Gender: "male", Address: "21 Huynh Thuc Khang", IdentityCard: "17687t562765786", EmployeeTypeID: 2, Avatar: "image2.jpg"}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "deliverystaff@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Tuan", Age: 37, Phone: 776664993, Gender: "male", Address: "21 Nhat Tao", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 4}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "loadpackagestaff@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Hung", Age: 47, Phone: 776334958, Gender: "male", Address: "84 Nguyen Trau", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg"}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}
	return nil
}

func createDefaultCustomer() error {
	userAuth := &model.UserAuthenticate{Email: "customer@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	customer := &model.Customer{UserAuthID: userAuth.ID, Name: "Customer One", Age: 18, Phone: 223334444, Gender: "male", Address: "12 Tran Hung Dao, Phuong 1, Quan 5"}
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	userAuth.CustomerID = customer.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "customer3@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	customer = &model.Customer{UserAuthID: userAuth.ID, Name: "Customer Three", Age: 18, Phone: 223334444, Gender: "female", Address: "13 Tran Hung Dao"}
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	userAuth.CustomerID = customer.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "customer2@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	customer = &model.Customer{UserAuthID: userAuth.ID, Name: "Customer Two", Age: 18, Phone: 223334444, Gender: "male", Address: "14 Tran Hung Dao"}
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	userAuth.CustomerID = customer.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	return nil
}

func createCustomerCredit() error {
	customerCredit := &model.CustomerCredit{CustomerID: 1, Phone: 223334444, ValidatePhone: true, AccountBalance: 30000000}
	if err := db.Create(customerCredit).Error; err != nil {
		return err
	}
	customerCredit = &model.CustomerCredit{CustomerID: 2, Phone: 223334444, ValidatePhone: true, AccountBalance: 30000000}
	if err := db.Create(customerCredit).Error; err != nil {
		return err
	}
	customerCredit = &model.CustomerCredit{CustomerID: 3, Phone: 223334444, ValidatePhone: true, AccountBalance: 30000000}
	if err := db.Create(customerCredit).Error; err != nil {
		return err
	}
	return nil
}

func createTransportType() error {
	transportType := &model.TransportType{SameCity: true, LocationOne: "HCM", ShortShipPricePerKm: 30000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{SameCity: true, LocationOne: "HCM", ShortShipPricePerKm: 26000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "HCM", LocationTwo: "DL", LongShipDuration: 86400, LongShipPrice: 100000, BusStationFrom: "231-233 Le Hong Phong", BusStationTo: "695-697 Quoc lo 20, Thi tran Lien Nghia, Huyen Duc Trong, Lam Dong", ShortShipPricePerKm: 30000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	transportType = &model.TransportType{LocationOne: "HCM", LocationTwo: "VT", LongShipDuration: 172800, LongShipPrice: 120000, BusStationFrom: "231-233 Le Hong Phong", BusStationTo: "192 Nam Ky Khoi Nghia, Phuong Thang Tam", ShortShipPricePerKm: 26000}
	if err := db.Create(transportType).Error; err != nil {
		return err
	}
	return nil
}

func createLongShip() error {
	longShip := &model.LongShip{TransportTypeID: 3, LicensePlate: "51A 435.22", EstimatedTimeOfDeparture: 1610599301, EstimatedTimeOfArrival: 1610999301, Finished: false}
	if err := db.Create(longShip).Error; err != nil {
		return err
	}
	longShip = &model.LongShip{TransportTypeID: 4, LicensePlate: "51B 425.82", EstimatedTimeOfDeparture: 1610099301, EstimatedTimeOfArrival: 1610399301, Finished: true}
	if err := db.Create(longShip).Error; err != nil {
		return err
	}
	return nil
}

func createDeliveryLocation() error {
	location := &model.DeliveryLocation{City: "HCM", District: "1"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "HCM", District: "2"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "HCM", District: "3"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "HCM", District: "Tan Binh"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "HCM", District: "Phu Nhuan"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "DL", District: "1"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "DL", District: "2"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "DL", District: "3"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "DL", District: "4"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "DL", District: "5"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrder() error {
	orderInfo := &model.OrderInfo{
		Weight: 2, Volume: 10, Type: "Normal", Image: "order1.png",
		CustomerSendID: 1, EmplCreateID: 2,
		Sender:           "269 Ngo Quyen, Quan 5, HCM",
		Receiver:         "38 Tran Hung Dao, Quan 1,HCM",
		Detail:           "May vi tinh ca nhan va ban phim may tinh",
		OrderShortShipID: 1, TransportTypeID: 1, ShortShipDistance: 20,
		TotalPrice: 200000, Note: "Giao hang vao buoi sang",
	}
	if err := db.Create(orderInfo).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderPay() error {
	orderPay := &model.OrderPay{
		OrderID: 1, PayStatus: true, TotalPrice: 200000, PayMethod: "cash",
	}
	if err := db.Create(orderPay).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderShortShip() error {
	orderShortShip := &model.OrderShortShip{
		OrderID: 1, ShipperID: 3, ShipperReceived: true,
		ShipperCalled: true, ShipperShipped: true,
		TimeConfirmed: 1608743635, CusReceiveConfirmed: true,
		Finished: true,
	}
	if err := db.Create(orderShortShip).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrder2() error {
	orderInfo := &model.OrderInfo{
		Weight: 3, Volume: 50, Type: "Special", Image: "order1.png",
		CustomerSendID: 1, EmplCreateID: 2,
		Sender:   "231-233 Le Hong Phong",
		Receiver: "74 Phan Chau Trinh, Quan 3, DL",
		Detail:   "May xay thit", UseLongShip: true, LongShipID: 1,
		OrderShortShipID: 1, TransportTypeID: 3, ShortShipDistance: 20,
		TotalPrice: 200000, Note: "Giao hang vao buoi trua",
	}
	if err := db.Create(orderInfo).Error; err != nil {
		return err
	}
	return nil
}

func createCustomerFCMToken() error {
	orderPay := &model.UserFCMToken{
		CustomerID: 1, UserAuthID: 5,
		Token: "f09fih-jQ9GhMz3riGfmJv:APA91bH7opzi3nvIeY1GLvJb0zZClx19ZztB5-6Bgg4jIsBi-9fnZWHpqYo1Za78W93VbdyiQureIFkck0MA6AaFik7LwQ2gIburmRCV2eR4ZBIp-YjQKRhIUHAYbu6YyQmfEJPsDgbn",
	}
	if err := db.Create(orderPay).Error; err != nil {
		return err
	}
	return nil
}
