package models

import "github.com/jinzhu/gorm"

type Enterprise struct {
	gorm.Model
	EnterpriseName  string    `json:"enterprise_name"`
	HeadquarterCity string    `json:"headquarter_city"`
	Vehicles        []Vehicle `gorm:"foreignKey:EnterpriseID" json:"-"`
	Drivers         []Driver  `gorm:"foreignKey:EnterpriseID" json:"-"`
}

type Driver struct {
	gorm.Model
	Enterprise   Enterprise `gorm:"foreignKey:EnterpriseID" json:"-"`
	EnterpriseID uint       `json:"enterprise_id"`
	VehicleID    uint
	IsActive     bool    `gorm:"index:bool_constraint,unique, where:is_active" json:"is_active"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Salary       float64 `json:"salary"`
}

type EnterpriseService interface {
	SaveEnterprise(Enterprise) error
	SaveDriver(Driver) error

	FindAllEnterprises() []Enterprise
	FindAllDrivers() []Driver
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

func (service *enterpriseSerivce) FindAllEnterprises() []Enterprise {
	return service.vehicleDB.FindAllEnterprises()
}

func (service *enterpriseSerivce) FindAllDrivers() []Driver {
	return service.vehicleDB.FindAllDrivers()
}
