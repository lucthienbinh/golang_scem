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

// GetGormInstance function
func GetGormInstance() *gorm.DB {
	return db
}

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
		&model.LongShipWorkflowData{},
		&model.OrderVoucher{},
		&model.CustomerNotification{},
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
	if err := createExampleOrderLongShip(); err != nil {
		return err
	}
	if err := createExampleOrderWorkflowData(); err != nil {
		return err
	}
	if err := createExampleOrder2(); err != nil {
		return err
	}
	if err := createOrderVoucher(); err != nil {
		return err
	}
	if err := createCustomerNotification(); err != nil {
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
		&model.LongShipWorkflowData{},
		&model.OrderVoucher{},
		&model.CustomerNotification{},
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
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Tuan", Age: 37, Phone: 776664993, Gender: "male", Address: "21 Nhat Tao", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 6}
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
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Hung", Age: 47, Phone: 776334958, Gender: "male", Address: "84 Nguyen Trau", IdentityCard: "17687t562765786", EmployeeTypeID: 4, Avatar: "image3.jpg"}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "deliverystaff2@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Hieu", Age: 37, Phone: 776664993, Gender: "male", Address: "21 Trung Son", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 11}
	if err := db.Create(employee).Error; err != nil {
		return err
	}
	userAuth.EmployeeID = employee.ID
	if err := db.Model(&userAuth).Updates(&userAuth).Error; err != nil {
		return err
	}

	userAuth = &model.UserAuthenticate{Email: "deliverystaff3@gmail.com", Password: "12345678"}
	if err := db.Create(userAuth).Error; err != nil {
		return err
	}
	employee = &model.Employee{UserAuthID: userAuth.ID, Name: "Thao", Age: 37, Phone: 776662293, Gender: "female", Address: "32 Xuan Son", IdentityCard: "17687t562774286", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 1}
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
	longShip := &model.LongShip{LSQrCode: "1611463545_4f163f5f0f9a621d.jpg", TransportTypeID: 3, LicensePlate: "51A 435.22", EstimatedTimeOfDeparture: 1610599301, EstimatedTimeOfArrival: 1610999301, Finished: true}
	if err := db.Create(longShip).Error; err != nil {
		return err
	}
	longShip = &model.LongShip{LSQrCode: "1611463464_d90bf90bbaa6cd39.jpg", TransportTypeID: 4, LicensePlate: "51B 425.82", EstimatedTimeOfDeparture: 1610099301, EstimatedTimeOfArrival: 1610399301, Finished: true}
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
	location = &model.DeliveryLocation{City: "VT", District: "1"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	location = &model.DeliveryLocation{City: "VT", District: "2"}
	if err := db.Create(location).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrder() error {
	orderInfo := &model.OrderInfo{
		Weight: 2, Volume: 10, Type: "Normal", Image: "box.jpg",
		CustomerSendID: 1, EmplCreateID: 2,
		Sender:           "Customer One - 269 Ngo Quyen, Quan 5, HCM - 5676765678",
		Receiver:         "Mai Thi Cuc - 38 Tran Hung Dao, Quan 1, HCM - 6765677867",
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

func createExampleOrderWorkflowData() error {
	orderPay := &model.OrderWorkflowData{
		OrderID: 1, WorkflowKey: "123abc", WorkflowInstanceKey: 321123321123, OrderPayID: 1, CustomerSendID: 1,
	}
	if err := db.Create(orderPay).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderLongShip() error {
	orderLongShip := &model.OrderLongShip{
		OrderID: 1, LongShipID: 1, CustomerSendID: 1,
	}
	if err := db.Create(orderLongShip).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrderShortShip() error {
	orderShortShip := &model.OrderShortShip{
		OrderID: 1, ShipperID: 3, CustomerSendID: 1, OSSQrCode: "1611465837_1e708c0f8e9a0213.jpg",
		Sender:   "Customer One - 269 Ngo Quyen, Quan 5, HCM - 5676765678",
		Receiver: "Mai Thi Cuc - 38 Tran Hung Dao, Quan 1, HCM - 6765677867",
	}
	if err := db.Create(orderShortShip).Error; err != nil {
		return err
	}
	return nil
}

func createExampleOrder2() error {
	orderInfo := &model.OrderInfo{
		Weight: 3, Volume: 50, Type: "Special", Image: "box.jpg",
		CustomerSendID: 1,
		Sender:         "Customer one - 231-233 Le Hong Phong - 6578678678",
		Receiver:       "Mac Thi Buoi - 74 Phan Chau Trinh, Quan 3, DL - 567865676",
		Detail:         "May xay thit", UseLongShip: true, LongShipID: 1,
		OrderShortShipID: 1, TransportTypeID: 3, ShortShipDistance: 20,
		TotalPrice: 200000, Note: "Giao hang vao buoi trua",
	}
	if err := db.Create(orderInfo).Error; err != nil {
		return err
	}
	return nil
}

func createOrderVoucher() error {
	orderVoucher := &model.OrderVoucher{
		Title:     "Khuyen mai ngay he 1",
		Content:   "Tan huong nhung ngay he sang khoang voi Move nice VN nha. Giam gia 30.000 tu hom nay den cuoi thang.",
		StartDate: 1611390576, EndDate: 1619390576, Discount: 30000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	orderVoucher = &model.OrderVoucher{
		Title:     "Khuyen mai ngay xuan 2",
		Content:   "Tan huong nhung ngay xuan mat me voi Move nice VN nha. Giam gia 50.000 tu hom nay den cuoi thang.",
		StartDate: 1615690576, EndDate: 1619390576, Discount: 50000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	orderVoucher = &model.OrderVoucher{
		Title:     "Khuyen mai ngay dong 3",
		Content:   "Tan huong nhung ngay xuan mat me voi Move nice VN nha. Giam gia 70.000 tu hom nay den cuoi thang.",
		StartDate: 1612390576, EndDate: 1620390576, Discount: 70000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	orderVoucher = &model.OrderVoucher{
		Title:     "Khuyen mai ngay phu nu 5",
		Content:   "Chuc mung ngay quoc te phu nu voi Move nice VN nha. Giam gia 100.000 tu hom nay den cuoi thang.",
		StartDate: 1612390576, EndDate: 1620390576, Discount: 100000,
	}
	if err := db.Create(orderVoucher).Error; err != nil {
		return err
	}
	return nil
}

func createCustomerNotification() error {
	customerNotification := &model.CustomerNotification{
		CustomerID: 1, Title: "Your order has started long ship trip",
		Content: "Order id: 1349014 Long ship id: 1268070",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	customerNotification = &model.CustomerNotification{
		CustomerID: 1, Title: "Your order has finished long ship trip",
		Content: "Order id: 1349014 Long ship id: 1268070",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	customerNotification = &model.CustomerNotification{
		CustomerID: 1, Title: "Shipper has called you",
		Content: "Your order id: 1349014 has been verified",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	customerNotification = &model.CustomerNotification{
		CustomerID: 1, Title: "Shipper has confirmed your package",
		Content: "Thanks for using our service. Finished order id: 1349014",
	}
	if err := db.Create(customerNotification).Error; err != nil {
		return err
	}
	return nil
}
