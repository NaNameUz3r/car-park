package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Enterprise struct {
	gorm.Model
	EnterpriseName  string    `json:"enterprise_name"`
	HeadquarterCity string    `json:"headquarter_city"`
	Vehicles        []Vehicle `gorm:"foreignKey:EnterpriseID" json:"-"`
	TimeZone        string    `json:"iana_timezone" validate:"timezone"`
}

type Driver struct {
	gorm.Model

	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Salary    float64 `json:"salary"`

	VehicleID uint `json:"vehicle_id" gorm:"default:null"`
	IsActive  bool `json:"is_active"`
}

type Manager struct {
	gorm.Model
	FirstName             string        `json:"first_name"`
	LastName              string        `json:"last_name"`
	Login                 string        `json:"login"`
	Password              string        `json:"password"`
	AccessibleEnterprises pq.Int64Array `gorm:"type:integer[]" json:"accesible_enterprises"`
}

type EnterpriseService interface {
	SaveEnterprise(Enterprise) error
	SaveDriver(Driver) error
	SaveManager(Manager) error
	SaveVehicle(Vehicle) (err error)
	UpdateVehicle(Vehicle) error
	DeleteVehicle(Vehicle) error
	VehicleByID(id uint) Vehicle

	ManagerVehicleByID(id uint) Vehicle

	FindAllEnterprisesByIDs(enterpriseIDs pq.Int64Array) []Enterprise
	FindAllDrivers() []Driver
	FindAllManagers() []Manager

	ManagerByID(id uint) Manager
	ManagerFindAllVehicles(accessibleEnterprises pq.Int64Array, paginations Pagination, preload bool) []Vehicle
	ManagerGetVehicleRoutesGeoJSON(vehicleID uint, notBefore string, notAfter string) (string, error)
	ManagerGetVehicleRoutesGeopoints(vehicleID uint, notBefore string, notAfter string) ([]GeoPoint, error)
	ManagerSaveGeoPoint(geoPoint GeoPoint) error
	GeoPointsByDates(notBefore, notAfter time.Time) []GeoPoint

	SaveRide(r Ride) error
	RideByID(id uint) Ride
	ManagerGetVehicleRides(vehicleID uint, notBefore string, notAfter string, inGeoJsons bool) ([]Ride, error)

	ManagerFindAllDrivers(accessibleEnterprises pq.Int64Array) []Driver
	ManagerByCreds(username, password string) Manager
}

type enterpriseSerivce struct {
	vehicleDB VehicleDB
}

func NewEnterpriseSerivce(vehicleDB VehicleDB) EnterpriseService {
	return &enterpriseSerivce{
		vehicleDB: vehicleDB,
	}
}

func (service *enterpriseSerivce) ManagerVehicleByID(id uint) Vehicle {
	return service.vehicleDB.VehicleByID(id)
}

func (service *enterpriseSerivce) SaveEnterprise(e Enterprise) error {
	return service.vehicleDB.SaveEnterprise(e)
}

func (service *enterpriseSerivce) SaveDriver(d Driver) error {
	return service.vehicleDB.SaveDriver(d)
}
func (service *enterpriseSerivce) SaveRide(r Ride) error {
	return service.vehicleDB.SaveRide(r)
}

func (service *enterpriseSerivce) SaveManager(m Manager) error {
	return service.vehicleDB.SaveManager(m)
}

func (service *enterpriseSerivce) SaveVehicle(v Vehicle) (err error) {
	return service.vehicleDB.SaveVehicle(v)
}

func (service *enterpriseSerivce) UpdateVehicle(v Vehicle) error {
	return service.vehicleDB.UpdateVehicle(v)
}

func (service *enterpriseSerivce) DeleteVehicle(v Vehicle) error {
	return service.vehicleDB.DeleteVehicle(v)
}

func (service *enterpriseSerivce) VehicleByID(id uint) Vehicle {
	return service.vehicleDB.VehicleByID(id)
}

func (service *enterpriseSerivce) RideByID(id uint) Ride {
	return service.vehicleDB.RideByID(id)
}

func (service *enterpriseSerivce) FindAllEnterprisesByIDs(enterprisesIDs pq.Int64Array) []Enterprise {
	return service.vehicleDB.FindAllEnterprisesByIDs(enterprisesIDs)
}

func (service *enterpriseSerivce) FindAllDrivers() []Driver {
	return service.vehicleDB.FindAllDrivers()
}
func (service *enterpriseSerivce) FindAllManagers() []Manager {
	return service.vehicleDB.FindAllManagers()
}

func (service *enterpriseSerivce) ManagerByID(id uint) Manager {
	return service.vehicleDB.ManagerByID(id)
}

func (service *enterpriseSerivce) ManagerByCreds(username, password string) Manager {
	return service.vehicleDB.ManagerByCreds(username, password)
}

func (service *enterpriseSerivce) ManagerFindAllVehicles(accessibleEnterprises pq.Int64Array, pagination Pagination, preload bool) []Vehicle {
	return service.vehicleDB.ManagerFindAllVehicles(accessibleEnterprises, pagination, preload)
}

func (service *enterpriseSerivce) ManagerFindAllDrivers(accessibleEnterprises pq.Int64Array) []Driver {
	return service.vehicleDB.ManagerFindAllDrivers(accessibleEnterprises)
}

func (service *enterpriseSerivce) ManagerGetVehicleRoutesGeoJSON(vehicleID uint, notBefore string, notAfter string) (string, error) {
	return service.vehicleDB.ManagerGetVehicleRoutesGeoJSON(vehicleID, notBefore, notAfter)
}

func (service *enterpriseSerivce) ManagerGetVehicleRoutesGeopoints(vehicleID uint, notBefore string, notAfter string) ([]GeoPoint, error) {
	return service.vehicleDB.ManagerGetVehicleRoutesGeopoints(vehicleID, notBefore, notAfter)
}

func (service *enterpriseSerivce) ManagerSaveGeoPoint(geoPoint GeoPoint) error {
	return service.vehicleDB.ManagerSaveGeoPoint(geoPoint)
}

func (service *enterpriseSerivce) GeoPointsByDates(notBefore, notAfter time.Time) []GeoPoint {
	return service.vehicleDB.GeoPointsByDates(notBefore, notAfter)
}
func (service *enterpriseSerivce) ManagerGetVehicleRides(vehicleID uint, notBefore string, notAfter string, inGeoJsons bool) ([]Ride, error) {
	return service.vehicleDB.ManagerGetVehicleRides(vehicleID, notBefore, notAfter, inGeoJsons)
}
