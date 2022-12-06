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
