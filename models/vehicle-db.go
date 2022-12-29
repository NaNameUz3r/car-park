package models

//TODO: Candidate on refactoring...

import (
	// "github.com/jinzhu/gorm"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	"car-park/utils"
	"fmt"
	"os"
	"strings"
	"time"

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

	ManagerGetVehicleRoutesGeoJSON(vehicleID uint, notBefore, notAfter string) (string, error)
	ManagerGetVehicleRoutesGeopoints(VehicleID uint, notBefore, notAfter string) ([]GeoPoint, error)
	ManagerSaveGeoPoint(geoPoint GeoPoint) error

	GeoPointsByDates(notBefore, notAfter time.Time) []GeoPoint

	SaveRide(r Ride) error
	RideByID(id uint) Ride
	ManagerGetVehicleRides(vehicleID uint, notBefore string, notAfter string, inGeoJsons bool) ([]Ride, error)

	ReportsByType(reportType, reportPeriod, notBefore, notAfter string) []Report

	FindAllEnterprisesByIDs(enterprisesIDs pq.Int64Array) []Enterprise
	FindAllDrivers() []Driver

	LoadFixturesGeotracks(vehicleID uint)
}

type dbConn struct {
	connection *gorm.DB
}

//TODO: Make it parameterizable and secure.
var dbHost = "host=" + os.Getenv("DB_HOST")
var connectionString = dbHost + " port=5432 user=admin password=qwerty dbname=car_park_dev sslmode=disable TimeZone=UTC"

func NewVehicleDB() VehicleDB {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger:               logger.Default.LogMode(logger.Info),
		FullSaveAssociations: true,
		NowFunc: func() time.Time {
			utc, _ := time.LoadLocation("UTC")
			return time.Now().In(utc)
		},
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

func (db *dbConn) RideByID(id uint) Ride {
	var ride Ride
	ride.ID = id
	db.connection.Find(&ride)
	return ride
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
		db.connection.Preload("CarModel").Where("enterprise_id IN ?", array).Find(&vehicles)
	} else {
		queryBuilder.Model(&Vehicle{}).Where("enterprise_id IN ?", array).Select("id", "enterprise_id", "description", "price", "mileage", "manufactured_year", "car_model_id", "commissioning_date").Find(&vehicles)
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

func (db *dbConn) FindAllEnterprisesByIDs(enterprisesID pq.Int64Array) []Enterprise {
	var enterprises []Enterprise

	array := make([]int64, len(enterprisesID))
	copy(array, enterprisesID)
	db.connection.Where("id IN ?", array).Select("id", "enterprise_name", "headquarter_city", "time_zone").Find(&enterprises)
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

func (db *dbConn) ManagerGetVehicleRoutesGeoJSON(vehicleID uint, notBefore string, notAfter string) (string, error) {
	var geoJson string

	baseQuery := "SELECT ST_AsGeoJSON(ST_MakeLine((ARRAY(SELECT ST_MakePoint(g.geo_x, g.geo_y) FROM geo_points As g where vehicle_id = ? "

	if notBefore != "" {
		baseQuery = baseQuery + "AND track_time >= '" + notBefore + "' "
	}

	if notAfter != "" {
		baseQuery = baseQuery + "AND track_time <= '" + notAfter + "'"
	}

	baseQuery = baseQuery + "ORDER BY track_time))))"

	result := db.connection.Raw(
		baseQuery, vehicleID)

	if result.Error != nil {
		return "", result.Error
	} else {
		result.Scan(&geoJson)
		geoJsonStrip := strings.ReplaceAll(geoJson, "\"", "'")
		return geoJsonStrip, nil
	}
}

func (db *dbConn) ManagerGetVehicleRoutesGeopoints(VehicleID uint, notBefore string, notAfter string) ([]GeoPoint, error) {
	var geoPoints []GeoPoint
	var result *gorm.DB
	if notAfter != "" && notBefore != "" {
		result = db.connection.Order("track_time asc").Where("vehicle_id = ?", VehicleID).Where("track_time >= ?", notBefore).Where("track_time <= ?", notAfter).Find(&geoPoints)
		return geoPoints, result.Error
	}

	if notAfter != "" && notBefore == "" {
		result = db.connection.Order("track_time asc").Where("vehicle_id = ?", VehicleID).Where("track_time <= ?", notAfter).Find(&geoPoints)
		return geoPoints, result.Error
	}

	if notAfter == "" && notBefore != "" {
		result = db.connection.Order("track_time asc").Where("vehicle_id = ?", VehicleID).Where("track_time >= ?", notBefore).Find(&geoPoints)
		return geoPoints, result.Error
	}

	if notAfter == "" && notBefore == "" {
		result = db.connection.Order("track_time asc").Where("vehicle_id = ?", VehicleID).Find(&geoPoints)
		return geoPoints, result.Error
	}

	return geoPoints, result.Error
}

func (db *dbConn) ManagerSaveGeoPoint(geoPoint GeoPoint) error {
	return db.connection.Create(&geoPoint).Error
}

func (db *dbConn) GeoPointsByDates(notBefore, notAfter time.Time) []GeoPoint {
	var geoPoints []GeoPoint
	db.connection.Order("track_time asc").Where("track_time >= ?", notBefore).Where("track_time <= ?", notAfter).Find(&geoPoints)
	return geoPoints
}

func (db *dbConn) SaveRide(r Ride) error {
	return db.connection.Create(&r).Error
}
func (db *dbConn) ManagerGetVehicleRides(vehicleID uint, notBefore string, notAfter string, inGeoJsons bool) ([]Ride, error) {
	var rides []Ride
	var result *gorm.DB

	if notAfter != "" && notBefore != "" {
		result = db.connection.Order("ride_start asc").Where("vehicle_id = ?", vehicleID).Where("ride_start >= ?", notBefore).Where("ride_finish <= ?", notAfter).Find(&rides)
	} else if notAfter != "" && notBefore == "" {
		result = db.connection.Order("ride_start asc").Where("vehicle_id = ?", vehicleID).Where("ride_finish <= ?", notAfter).Find(&rides)
	} else if notAfter == "" && notBefore != "" {
		result = db.connection.Order("ride_start asc").Where("vehicle_id = ?", vehicleID).Where("ride_start >= ?", notBefore).Find(&rides)
	} else if notAfter == "" && notBefore == "" {
		result = db.connection.Order("ride_start asc").Where("vehicle_id = ?", vehicleID).Find(&rides)
	}

	var err2 error
	for i, ride := range rides {
		nb, na, err := utils.TimeStampsToUTCStrings(ride.RideStart, ride.RideFinish)
		if err != nil {
			return rides, err
		}
		if inGeoJsons {
			rides[i].GeoJSON, err2 = db.ManagerGetVehicleRoutesGeoJSON(ride.VehicleID, nb, na)
			if err != nil {
				return rides, err2
			}
		} else {
			rides[i].GeoPoints, err2 = db.ManagerGetVehicleRoutesGeopoints(ride.VehicleID, nb, na)
			if err != nil {
				return rides, err2
			}
		}
	}

	return rides, result.Error
}

//TODO: DELETE IT IS REDUNDANT
func (db *dbConn) ReportsByType(reportType, reportPeriod, notBefore, notAfter string) []Report {
	var reports []Report
	db.connection.Where("report_type = ?", reportType).Where("not_before >= ?", notBefore).Where("not_after <= ?", notAfter).Find(&reports)
	return reports
}

func (db *dbConn) LoadFixturesGeotracks(vehicleID uint) {
	LoadFixtureGeotracksRIDE1(db.connection, vehicleID)
	LoadFixtureGeotracksRIDE2(db.connection, vehicleID)
	LoadFixtureGeotracksRIDE3(db.connection, vehicleID)
}

func (db *dbConn) AutoMigrate() {
	db.connection.AutoMigrate(&Enterprise{}, &Manager{})
	db.connection.AutoMigrate(&CarModel{}, &Vehicle{}, &Driver{}, &GeoPoint{}, &Ride{})
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
