package vehicle

import "time"

type Location struct {
	ID         int64
	VehicleID  string
	Latitude   float64
	Longitude  float64
	RecordedAt time.Time
	CreatedAt  time.Time
}
