package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Enterprise struct {
	gorm.Model
	EnterpriseName  string    `json:"enterprise_name"`
	HeadquarterCity string    `json:"headquarter_city"`
	Vehicles        []Vehicle `gorm:"foreignKey:EnterpriseID" json:"-"`
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

	FindAllEnterprises() []Enterprise
	FindAllDrivers() []Driver
	FindAllManagers() []Manager

	ManagerByID(id uint) Manager
	ManagerFindAllVehicles(accessibleEnterprises pq.Int64Array, preload bool) []Vehicle
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

func (service *enterpriseSerivce) SaveEnterprise(e Enterprise) error {
	return service.vehicleDB.SaveEnterprise(e)
}

func (service *enterpriseSerivce) SaveDriver(d Driver) error {
	return service.vehicleDB.SaveDriver(d)
}

func (service *enterpriseSerivce) SaveManager(m Manager) error {
	return service.vehicleDB.SaveManager(m)
}

func (service *enterpriseSerivce) FindAllEnterprises() []Enterprise {
	return service.vehicleDB.FindAllEnterprises()
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

func (service *enterpriseSerivce) ManagerFindAllVehicles(accessibleEnterprises pq.Int64Array, preload bool) []Vehicle {
	return service.vehicleDB.ManagerFindAllVehicles(accessibleEnterprises, preload)
}

func (service *enterpriseSerivce) ManagerFindAllDrivers(accessibleEnterprises pq.Int64Array) []Driver {
	return service.vehicleDB.ManagerFindAllDrivers(accessibleEnterprises)
}
