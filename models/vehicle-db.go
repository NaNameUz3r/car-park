package models

//TODO: Candidate on refactoring...

import (
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type VehicleDB interface {
	SaveVehicle(vehicle Vehicle) error
	UpdateVehicle(vehicle Vehicle) error
	DeleteVehicle(vehicle Vehicle) error
	FindAllVehicles(preload bool) []Vehicle
	VehicleByID(id uint) Vehicle
	FindAllCarModels() []CarModel
	CloseConnection()
	DestructiveReset() error

	SaveEnterprise(enterprise Enterprise) error
	SaveDriver(driver Driver) error

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

	db.AutoMigrate(&Enterprise{})
	db.AutoMigrate(&CarModel{}, &Vehicle{}, &Driver{})
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

func (db *dbConn) SaveVehicle(v Vehicle) error {
	return db.connection.Create(&v).Error
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

func (db *dbConn) SaveEnterprise(e Enterprise) error {
	return db.connection.Create(&e).Error
}

func (db *dbConn) SaveDriver(d Driver) error {
	return db.connection.Create(&d).Error
}

// DROPDATABASE!
func (db *dbConn) DestructiveReset() error {
	err := db.connection.Delete(&Vehicle{}, &CarModel{}).Error
	if err != nil {
		return err
	}

	//TODO: This need to be refactored
	db.connection.AutoMigrate(&Vehicle{})
	return nil
}
