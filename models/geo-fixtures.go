package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

var timeLayout string = "2006-01-02 15:06:05"

func parseTime(s string) time.Time {
	utcLoc, _ := time.LoadLocation("UTC")
	time, err := time.ParseInLocation(time.RFC3339, s, utcLoc)
	if err != nil {
		fmt.Println("[WARNING] Time not parsed in fixtures")
	}
	return time
}

func LoadFixtureGeotracksRIDE1(db *gorm.DB, vehicleID uint) {
	data := []GeoPoint{
		{
			TrackTime: parseTime("2022-09-16T05:48:07Z"),
			GeoX:      -83.01925,
			GeoY:      42.34551,
			GeoZ:      136.1999969482422,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T05:50:24Z"),
			GeoX:      -83.26635,
			GeoY:      42.24106,
			GeoZ:      126.5999984741211,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T05:55:25Z"),
			GeoX:      -83.03009,
			GeoY:      42.45232,
			GeoZ:      123.0,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:01:26Z"),
			GeoX:      -83.02666,
			GeoY:      42.45325,
			GeoZ:      120.5,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:05:27Z"),
			GeoX:      -83.02638,
			GeoY:      42.45705,
			GeoZ:      118.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:10:29Z"),
			GeoX:      -83.02656,
			GeoY:      42.46037,
			GeoZ:      119.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:20:30Z"),
			GeoX:      -83.02362,
			GeoY:      42.46316,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:24:30Z"),
			GeoX:      -83.27922,
			GeoY:      42.23849,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:30:30Z"),
			GeoX:      -83.02196,
			GeoY:      42.46322,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:35:30Z"),
			GeoX:      -83.02042,
			GeoY:      42.46322,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:35:30Z"),
			GeoX:      -83.01354,
			GeoY:      42.46338,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:40:30Z"),
			GeoX:      -83.00768,
			GeoY:      42.46343,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:41:30Z"),
			GeoX:      -83.00582,
			GeoY:      42.46791,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T06:58:30Z"),
			GeoX:      -83.00608,
			GeoY:      42.47541,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T07:01:30Z"),
			GeoX:      -83.00337,
			GeoY:      42.47797,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T07:06:30Z"),
			GeoX:      -82.99358,
			GeoY:      42.47807,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},

		{
			TrackTime: parseTime("2022-09-16T07:10:30Z"),
			GeoX:      -82.98736,
			GeoY:      42.4781,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-16T07:12:30Z"),
			GeoX:      -83.27586,
			GeoY:      42.23434,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
	}

	for _, d := range data {
		result := db.Create(&d)
		fmt.Print(result.Error, "AAAAAAAAAAAAAAAAAAAAA?")
	}
}

func LoadFixtureGeotracksRIDE2(db *gorm.DB, vehicleID uint) {
	data := []GeoPoint{
		{
			TrackTime: parseTime("2022-09-17T05:48:07Z"),
			GeoX:      -83.26389,
			GeoY:      42.24143,
			GeoZ:      136.1999969482422,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T05:50:24Z"),
			GeoX:      -83.26635,
			GeoY:      42.24106,
			GeoZ:      126.5999984741211,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T05:55:25Z"),
			GeoX:      -83.03009,
			GeoY:      42.45232,
			GeoZ:      123.0,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:01:26Z"),
			GeoX:      -83.02666,
			GeoY:      42.45325,
			GeoZ:      120.5,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:05:27Z"),
			GeoX:      -83.02638,
			GeoY:      42.45705,
			GeoZ:      118.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:10:29Z"),
			GeoX:      -83.02656,
			GeoY:      42.46037,
			GeoZ:      119.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:20:30Z"),
			GeoX:      -83.02362,
			GeoY:      42.46316,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:24:30Z"),
			GeoX:      -83.27922,
			GeoY:      42.23849,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:30:30Z"),
			GeoX:      -83.02196,
			GeoY:      42.46322,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:35:30Z"),
			GeoX:      -83.02042,
			GeoY:      42.46322,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:35:30Z"),
			GeoX:      -83.01354,
			GeoY:      42.46338,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:40:30Z"),
			GeoX:      -83.00768,
			GeoY:      42.46343,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T06:41:30Z"),
			GeoX:      -83.00582,
			GeoY:      42.46791,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-17T07:04:30Z"),
			GeoX:      -83.02428,
			GeoY:      42.34905,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
	}

	for _, d := range data {
		result := db.Create(&d)
		fmt.Print(result.Error, "AAAAAAAAAAAAAAAAAAAAA?")
	}
}

func LoadFixtureGeotracksRIDE3(db *gorm.DB, vehicleID uint) {
	data := []GeoPoint{
		{
			TrackTime: parseTime("2022-09-20T05:48:07Z"),
			GeoX:      -83.03078,
			GeoY:      42.44916,
			GeoZ:      136.1999969482422,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T05:50:24Z"),
			GeoX:      -83.02988,
			GeoY:      42.45074,
			GeoZ:      126.5999984741211,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T05:55:25Z"),
			GeoX:      -83.03009,
			GeoY:      42.45232,
			GeoZ:      123.0,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:01:26Z"),
			GeoX:      -83.02666,
			GeoY:      42.45325,
			GeoZ:      120.5,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:05:27Z"),
			GeoX:      -83.02638,
			GeoY:      42.45705,
			GeoZ:      118.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:10:29Z"),
			GeoX:      -83.02656,
			GeoY:      42.46037,
			GeoZ:      119.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:20:30Z"),
			GeoX:      -83.02362,
			GeoY:      42.46316,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:24:30Z"),
			GeoX:      -83.27922,
			GeoY:      42.23849,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:30:30Z"),
			GeoX:      -83.02196,
			GeoY:      42.46322,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:35:30Z"),
			GeoX:      -83.02042,
			GeoY:      42.46322,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:35:30Z"),
			GeoX:      -83.01354,
			GeoY:      42.46338,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:40:30Z"),
			GeoX:      -83.00768,
			GeoY:      42.46343,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:41:30Z"),
			GeoX:      -83.00582,
			GeoY:      42.46791,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T06:58:30Z"),
			GeoX:      -83.00608,
			GeoY:      42.47541,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:01:30Z"),
			GeoX:      -83.00337,
			GeoY:      42.47797,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:06:30Z"),
			GeoX:      -82.99358,
			GeoY:      42.47807,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},

		{
			TrackTime: parseTime("2022-09-20T07:10:30Z"),
			GeoX:      -82.98736,
			GeoY:      42.4781,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:12:30Z"),
			GeoX:      -82.98701,
			GeoY:      42.4879,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:20:30Z"),
			GeoX:      -82.95757,
			GeoY:      42.49379,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:22:30Z"),
			GeoX:      -82.93877,
			GeoY:      42.49405,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:26:30Z"),
			GeoX:      -82.91782,
			GeoY:      42.51017,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:33:30Z"),
			GeoX:      -82.85922,
			GeoY:      42.558,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:37:30Z"),
			GeoX:      -82.85613,
			GeoY:      42.57047,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:41:30Z"),
			GeoX:      -82.83338,
			GeoY:      42.57542,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:45:30Z"),
			GeoX:      -82.81343,
			GeoY:      42.58241,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:50:30Z"),
			GeoX:      -82.81365,
			GeoY:      42.59605,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T07:55:30Z"),
			GeoX:      -82.83245,
			GeoY:      42.59782,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T08:03:30Z"),
			GeoX:      -82.83248,
			GeoY:      42.59822,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T08:05:30Z"),
			GeoX:      -82.82231,
			GeoY:      42.59835,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T08:10:30Z"),
			GeoX:      -82.81836,
			GeoY:      42.59994,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
		{
			TrackTime: parseTime("2022-09-20T08:13:30Z"),
			GeoX:      -82.81414,
			GeoY:      42.60028,
			GeoZ:      120.9000015258789,
			VehicleID: vehicleID,
		},
	}

	for _, d := range data {
		result := db.Create(&d)
		fmt.Print(result.Error, "AAAAAAAAAAAAAAAAAAAAA?")
	}
}
