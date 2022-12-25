package models

import (
	"time"
)

type Report struct {
	ReportType string                    `json:"report_type"`
	VehicleID  uint                      `json:"vehicle_id"`
	TimePeriod string                    `json:"time_repiod"`
	NotBefore  time.Time                 `json:"not_before"`
	NotAfter   time.Time                 `json:"not_after"`
	Results    map[time.Time]interface{} `gorm:"-" json:"results"`
}

// type MilageResult struct {
// 	SegmentStartTime time.Time
// 	Milage           float64
// }
