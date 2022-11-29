package models

import (
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	Description      string `json:"description" form:"description"`
	Price            uint   `gorm:"not null" json:"price" form:"price"`
	Mileage          uint   `json:"mileage" form:"mileage"`
	ManufacturedYear uint   `json:"manufactured" form:"manufactured"`

	Enterprise   Enterprise `gorm:"foreignKey:EnterpriseID" json:"-"`
	EnterpriseID uint       `json:"enterprise_id" form:"enterprise_id"`
	CarModel     CarModel   `json:"car_model" binding:"required"`
	CarModelID   uint       `json:"carmodel_id" form:"-"`
	Drivers      []Driver   `json:"-" gorm:"foreignKey:VehicleID"`
}

type CarModel struct {
	gorm.Model
	Brand            string  `json:"brand" form:"brand"`
	CarType          string  `json:"carmodel" form:"carmodel"`
	FuelTankCapacity uint    `json:"fueltankcapacity" form:"fueltankcapacity"`
	Seats            uint    `json:"seats" form:"seats"`
	EngineVolume     float64 `json:"enginevolume" form:"enginevolume"`
}

type VehicleService interface {
	SaveVehicle(Vehicle) (err error)
	UpdateVehicle(Vehicle) error
	DeleteVehicle(Vehicle) error
	FindAllVehicles(preload bool) []Vehicle
	FindAllCarModels() []CarModel
	VehicleByID(id uint) Vehicle
}

type vehicleService struct {
	vehicleDB VehicleDB
}

func NewVehicleService(vehicleDB VehicleDB) VehicleService {
	return &vehicleService{
		vehicleDB: vehicleDB,
	}
}

func (service *vehicleService) SaveVehicle(v Vehicle) (err error) {
	return service.vehicleDB.SaveVehicle(v)
}

func (service *vehicleService) UpdateVehicle(v Vehicle) error {
	return service.vehicleDB.UpdateVehicle(v)
}

func (service *vehicleService) DeleteVehicle(v Vehicle) error {
	return service.vehicleDB.DeleteVehicle(v)
}

func (service *vehicleService) FindAllVehicles(preload bool) []Vehicle {
	return service.vehicleDB.FindAllVehicles(preload)
}

func (service *vehicleService) FindAllCarModels() []CarModel {
	return service.vehicleDB.FindAllCarModels()
}

func (service *vehicleService) VehicleByID(id uint) Vehicle {
	return service.vehicleDB.VehicleByID(id)
}
