package vehicle

import (
	"context"
	"fleet-management-system/internal/fleet/geofence"
	"fleet-management-system/internal/transport/rabbitmq"
	"log"
	"time"
)

type Service struct {
	repo              *Repository
	geofencePublisher *rabbitmq.Publisher
}

func NewService(repo *Repository, pub *rabbitmq.Publisher) *Service {
	return &Service{
		repo:              repo,
		geofencePublisher: pub,
	}
}

func (s *Service) RecordLocation(
	ctx context.Context,
	vehicleID string,
	lat float64,
	long float64,
	recordedAt time.Time,
) error {
	loc := Location{
		VehicleID:  vehicleID,
		Latitude:   lat,
		Longitude:  long,
		RecordedAt: recordedAt,
	}

	if err := s.repo.InsertLocation(ctx, loc); err != nil {
		return err
	}

	if geofence.IsInsideGeofence(lat, long) {
		event := geofence.GeofenceEvent{
			VehicleID: vehicleID,
			Event:     "geofence_entry",
			Timestamp: recordedAt.Unix(),
		}
		event.Location.Latitude = lat
		event.Location.Longitude = long

		if err := s.geofencePublisher.Publish(event); err != nil {
			log.Println("failed publish geofence event:", err)
		}
	}

	return nil
}

func (s *Service) GetLatestLocation(
	ctx context.Context,
	vehicleID string,
) (Location, error) {
	return s.repo.GetLatestLocation(ctx, vehicleID)
}

func (s *Service) GetLocationHistory(
	ctx context.Context,
	vehicleID string,
	start, end time.Time,
) ([]Location, error) {
	return s.repo.GetLocationHistory(ctx, vehicleID, start, end)
}
