package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

var timeLayout string = "2006-01-02 15:04:05"

func parseTime(s string) time.Time {
	t, err := time.Parse(timeLayout, s)
	if err != nil {
		fmt.Println("[WARNING] Time not parsed in fixtures")
	}
	return t
}

func LoadFixtureGeotracks(db *gorm.DB, vehicleID uint) {
	data := []GeoPoint{
		{
			TrackTime: parseTime("2015-09-15 03:48:07.235"),
			GeoX:      35.015021,
			GeoY:      32.519585,
			GeoZ:      136.1999969482422,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2015-09-15 03:48:24.734"),
			GeoX:      35.014954,
			GeoY:      32.519606,
			GeoZ:      126.5999984741211,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2015-09-15 03:48:25.660"),
			GeoX:      35.014871,
			GeoY:      32.519612,
			GeoZ:      123.0,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2015-09-15 03:48:26.819"),
			GeoX:      35.014824,
			GeoY:      32.519654,
			GeoZ:      120.5,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2015-09-15 03:48:27.828"),
			GeoX:      35.014776,
			GeoY:      32.519689,
			GeoZ:      118.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2015-09-15 03:48:29.720"),
			GeoX:      35.014704,
			GeoY:      32.519691,
			GeoZ:      119.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2015-09-15 03:48:30.669"),
			GeoX:      35.014657,
			GeoY:      32.519734,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
	}

	for _, d := range data {
		result := db.Create(&d)
		fmt.Print(result.Error, "AAAAAAAAAAAAAAAAAAAAA?")
	}
}
