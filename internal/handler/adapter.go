package handler

import (
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
		DSN:                  "user=postgres password=postgres dbname=scem_database port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	return err
}

// ConnectMySQL to open connect with database
func ConnectMySQL() (err error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(localhost:3308)/scem_database?charset=utf8mb4&parseTime=True&loc=Local"
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
		// &model.OrderStatusJSON{},
		&model.TransportType{},
		&model.UserAuthenticate{},
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
	return
}

func deleteDatabase() (err error) {
	return db.Migrator().DropTable(
		&model.Employee{},
		&model.Customer{},
		&model.EmployeeType{},
		&model.DeliveryLocation{},
		&model.OrderInfo{},
		// &model.OrderStatusJSON{},
		&model.TransportType{},
		&model.UserAuthenticate{},
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
	return
}

func createDefaultEmployee() (err error) {
	userAuth := model.UserAuthenticate{Email: "admin@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	employee := model.Employee{UserAuthID: userAuth.ID, Name: "Binh", Age: 18, Phone: 223334444, Gender: "male", Address: "12 Tran Hung Dao", IdentityCard: "17687t562765786", EmployeeTypeID: 1, Avatar: "image1.jpg"}
	if err := db.Create(&employee).Error; err != nil {
		return err
	}
	userAuth = model.UserAuthenticate{Email: "inputstaff@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	employee = model.Employee{UserAuthID: userAuth.ID, Name: "Hoa", Age: 18, Phone: 448883333, Gender: "male", Address: "21 Huynh Thuc Khang", IdentityCard: "17687t562765786", EmployeeTypeID: 2, Avatar: "image2.jpg"}
	if err := db.Create(&employee).Error; err != nil {
		return err
	}
	userAuth = model.UserAuthenticate{Email: "deliverystaff@gmail.com", Password: "12345678"}
	if err := db.Create(&userAuth).Error; err != nil {
		return err
	}
	employee = model.Employee{UserAuthID: userAuth.ID, Name: "Tuan", Age: 18, Phone: 776664993, Gender: "male", Address: "21 Nhat Tao", IdentityCard: "17687t562765786", EmployeeTypeID: 3, Avatar: "image3.jpg", DeliveryLocationID: 4}
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
	transportType := model.TransportType{Name: "HCM", RouteFixedPrice: 0, PricePerKm: 30000}
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
		EmployeeID: 3, Receiver: "253 Tran Hung Dao, 23114321412",
		Detail: "May vi tinh ca nhan", TotalPrice: 200000, Note: "Giao hang vao buoi sang",
	}
	if err := db.Create(&orderInfo).Error; err != nil {
		return err
	}
	return
}
