package models

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type GeoPoint struct {
	gorm.Model `json:"-"`
	TrackTime  time.Time `json:"track_time"`
	GeoX       float64   `json:"geo_x"`
	GeoY       float64   `json:"geo_y"`
	GeoZ       float64   `json:"geo_z"`
	VehicleID  uint      `json:"vehicle_id"`
}

func (g *GeoPoint) AfterFind(tx *gorm.DB) (err error) {
	validate = validator.New()
	loc, _ := time.LoadLocation("UTC")
	g.TrackTime = g.TrackTime.In(loc)

	var Vehicle Vehicle
	tx.First(&Vehicle, "id = ?", g.VehicleID)

	var enterprise Enterprise
	tx.First(&enterprise, "id = ?", Vehicle.EnterpriseID)
	errTimeZone := validate.Struct(enterprise)
	if errTimeZone == nil && enterprise.TimeZone != "" {
		locEnt, _ := time.LoadLocation(enterprise.TimeZone)
		g.TrackTime = g.TrackTime.In(locEnt)
	} else if errTimeZone != nil {

		log.Println("WARNING BROKEN TIMEZONE ", err)
	}

	return
}

type Ride struct {
	gorm.Model
	RideStart    time.Time  `json:"ride_start_time"`
	RideFinish   time.Time  `json:"ride_finish_time"`
	VehicleID    uint       `json:"vehicle_id"`
	GeoPoints    []GeoPoint `gorm:"-" json:"geo_points,omitempty"`
	GeoJSON      string     `gorm:"-" json:"geo_json,omitempty"`
	RideDistance float64    `json:"ride_distance"`
}

func (r *Ride) BeforeCreate(tx *gorm.DB) (err error) {
	loc, _ := time.LoadLocation("UTC")
	r.RideStart = r.RideStart.In(loc)
	r.RideFinish = r.RideFinish.In(loc)
	return
}

func (r *Ride) AfterFind(tx *gorm.DB) (err error) {
	validate = validator.New()
	loc, _ := time.LoadLocation("UTC")
	r.RideStart = r.RideStart.In(loc)
	r.RideFinish = r.RideFinish.In(loc)

	var Vehicle Vehicle
	tx.First(&Vehicle, "id = ?", r.VehicleID)

	var enterprise Enterprise
	tx.First(&enterprise, "id = ?", Vehicle.EnterpriseID)
	errTimeZone := validate.Struct(enterprise)
	if errTimeZone == nil && enterprise.TimeZone != "" {
		locEnt, _ := time.LoadLocation(enterprise.TimeZone)
		r.RideStart = r.RideStart.In(locEnt)
		r.RideFinish = r.RideFinish.In(locEnt)
	} else if errTimeZone != nil {

		log.Println("WARNING BROKEN TIMEZONE ", err)
	}

	return
}

type HumanReadRide struct {
	VehicleID uint `json:"vehicle_id"`
	RideID    uint `json:"ride_id"`

	RideStart    time.Time `json:"ride_start_time"`
	StartAddress string    `json:"start_address"`

	RideFinish    time.Time `json:"ride_finish_time"`
	FinishAddress string    `json:"finish_address"`

	RideDuration string `json:"ride_duration"`
}
