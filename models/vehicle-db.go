package models

//TODO: Candidate on refactoring...

import (
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"

	"github.com/lib/pq"
	"golang.org/x/exp/slices"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type VehicleDB interface {
	SaveVehicle(vehicle Vehicle) (err error)
	UpdateVehicle(vehicle Vehicle) error
	DeleteVehicle(vehicle Vehicle) error
	FindAllVehicles(preload bool) []Vehicle
	VehicleByID(id uint) Vehicle
	FindAllCarModels() []CarModel
	CloseConnection()
	AutoMigrate()
	DestructiveReset() error

	SaveEnterprise(enterprise Enterprise) error
	SaveDriver(driver Driver) error
	SaveManager(manager Manager) error
	FindAllManagers() []Manager
	ManagerByID(id uint) Manager
	ManagerByCreds(username, password string) Manager
	ManagerFindAllVehicles(accessibleEnterprises pq.Int64Array, pagination Pagination, preload bool) []Vehicle
	ManagerFindAllDrivers(accessibleEnterprises pq.Int64Array) []Driver

	FindAllEnterprises() []Enterprise
	FindAllDrivers() []Driver
}

type dbConn struct {
	connection *gorm.DB
}

//TODO: Make it parameterizable and secure.
var connectionString = "host=localhost port=5432 user=admin password=qwerty dbname=car_park_dev sslmode=disable"

func NewVehicleDB() VehicleDB {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger:               logger.Default.LogMode(logger.Info),
		FullSaveAssociations: true,
	})
	if err != nil {
		panic("Failed to connect to database")
	}
	// db.LogMode(true)

	return &dbConn{
		connection: db,
	}
}

func (db *dbConn) CloseConnection() {
	sqlDb, err := db.connection.DB()
	if err != nil {
		panic("Failed to close database connection (fetching interface)")
	}

	err = fmt.Errorf(sqlDb.Close().Error())

	if err != nil {
		panic("Failed to close database connection")
	}
}

func (db *dbConn) SaveVehicle(v Vehicle) (err error) {
	result := db.connection.Create(&v)
	return result.Error
}

func (db *dbConn) UpdateVehicle(v Vehicle) error {
	return db.connection.Save(&v).Error
}

func (db *dbConn) DeleteVehicle(v Vehicle) error {
	var findV Vehicle
	findV.ID = v.ID
	db.connection.Preload("CarModel").Find(&findV)
	fmt.Println(findV)
	var carM CarModel
	carM.ID = findV.CarModel.ID
	fmt.Println(carM)
	err := db.connection.Delete(&carM).Error
	if err != nil {
		return err
	}
	return db.connection.Delete(&findV).Error
}

func (db *dbConn) VehicleByID(id uint) Vehicle {
	var vehicle Vehicle
	vehicle.ID = id
	db.connection.Preload("CarModel").Find(&vehicle)
	return vehicle

}

func (db *dbConn) FindAllVehicles(preload bool) []Vehicle {
	var vehicles []Vehicle
	var vehicle Vehicle
	vehicle.ID = 1
	if preload {
		db.connection.Preload("CarModel").Find(&vehicles)
	} else {
		db.connection.Preload("Drivers").Select("id", "enterprise_id", "description", "price", "mileage", "manufactured_year", "car_model_id").Find(&vehicles)
		// db.connection.Model(&vehicle).Association("Drivers").Find(&vehicle)
	}
	return vehicles
}

func (db *dbConn) ManagerFindAllVehicles(accessibleEnterprises pq.Int64Array, pagination Pagination, preload bool) []Vehicle {

	array := make([]int64, len(accessibleEnterprises))
	copy(array, accessibleEnterprises)

	var vehicles []Vehicle
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := db.connection.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	if preload {
		db.connection.Preload("CarModel").Find(&vehicles)
	} else {
		queryBuilder.Model(&Vehicle{}).Where("enterprise_id IN ?", array).Select("id", "enterprise_id", "description", "price", "mileage", "manufactured_year", "car_model_id").Find(&vehicles)
		// db.connection.Preload("Drivers").Where("enterprise_id IN ?", array).Select("id", "enterprise_id", "description", "price", "mileage", "manufactured_year", "car_model_id").Find(&vehicles)
		// db.connection.Model(&vehicle).Association("Drivers").Find(&vehicle)
	}
	return vehicles
}

func (db *dbConn) FindAllCarModels() []CarModel {
	var carModels []CarModel
	db.connection.Find(&carModels)
	return carModels
}

func (db *dbConn) FindAllEnterprises() []Enterprise {
	var enterprises []Enterprise
	db.connection.Select("id", "enterprise_name", "headquarter_city").Find(&enterprises)
	return enterprises
}

func (db *dbConn) FindAllDrivers() []Driver {
	var drivers []Driver
	db.connection.Find(&drivers)
	return drivers
}

func (db *dbConn) ManagerFindAllDrivers(accessibleEnterprises pq.Int64Array) []Driver {

	array := make([]int64, len(accessibleEnterprises))
	copy(array, accessibleEnterprises)
	var vehicles []Vehicle

	var drivers []Driver
	db.connection.Preload("Drivers").Where("enterprise_id IN ?", array).Select("id", "enterprise_id", "description", "price", "mileage", "manufactured_year", "car_model_id").Find(&vehicles)
	for _, v := range vehicles {
		if slices.Contains(array, int64(v.EnterpriseID)) {
			for _, d := range v.Drivers {
				drivers = append(drivers, d)
			}
		}
	}
	return drivers
}

func (db *dbConn) FindAllManagers() []Manager {
	var managers []Manager
	db.connection.Find(&managers)
	return managers
}

func (db *dbConn) ManagerByID(id uint) Manager {
	var manager Manager
	manager.ID = id
	db.connection.Find(&manager)
	return manager

}

func (db *dbConn) ManagerByCreds(username, password string) Manager {
	var manager Manager
	// manager.Login = username
	// manager.Password = password

	db.connection.Where("login = ?", username).Where("password = ?", password).First(&manager)
	fmt.Println("MANAGER FROM DB: ", manager)
	return manager

}

func (db *dbConn) SaveEnterprise(e Enterprise) error {
	return db.connection.Create(&e).Error
}

func (db *dbConn) SaveDriver(d Driver) error {
	return db.connection.Create(&d).Error
}

func (db *dbConn) SaveManager(m Manager) error {
	return db.connection.Create(&m).Error
}

func (db *dbConn) AutoMigrate() {
	db.connection.AutoMigrate(&Enterprise{}, &Manager{})
	db.connection.AutoMigrate(&CarModel{}, &Vehicle{}, &Driver{})
}

// DROPDATABASE!
func (db *dbConn) DestructiveReset() error {
	err := db.connection.Delete(&Vehicle{}, &CarModel{}, &Enterprise{}, &Manager{}, &Driver{}).Error
	if err != nil {
		return err
	}

	//TODO: This need to be refactored
	// db.AutoMigrate()
	return nil
}
