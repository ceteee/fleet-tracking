package vehicle

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertLocation(
	ctx context.Context,
	loc Location,
) error {
	query := `
		INSERT INTO vehicle_locations
			(vehicle_id, latitude, longitude, recorded_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		loc.VehicleID,
		loc.Latitude,
		loc.Longitude,
		loc.RecordedAt,
	)

	return err
}

func (r *Repository) GetLatestLocation(
	ctx context.Context,
	vehicleID string,
) (Location, error) {
	var loc Location

	query := `
		SELECT id, vehicle_id, latitude, longitude, recorded_at, created_at
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		vehicleID,
	).Scan(
		&loc.ID,
		&loc.VehicleID,
		&loc.Latitude,
		&loc.Longitude,
		&loc.RecordedAt,
		&loc.CreatedAt,
	)

	return loc, err
}

func (r *Repository) GetLocationHistory(
	ctx context.Context,
	vehicleID string,
	start, end time.Time,
) ([]Location, error) {
	query := `
		SELECT id, vehicle_id, latitude, longitude, recorded_at, created_at	
		FROM vehicle_locations
		WHERE vehicle_id = $1
			AND recorded_at BETWEEN $2 AND $3
		ORDER BY recorded_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, vehicleID, start, end)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(
			&loc.ID,
			&loc.VehicleID,
			&loc.Latitude,
			&loc.Longitude,
			&loc.RecordedAt,
			&loc.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, loc)
	}

	return result, nil
}
