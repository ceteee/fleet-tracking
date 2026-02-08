BEGIN;

DROP TABLE IF EXISTS vehicle_locations;

CREATE TABLE vehicle_locations (
    id BIGSERIAL PRIMARY KEY,
    vehicle_id VARCHAR(20) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    recorded_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_vehicle_locations_vehicle_time
ON vehicle_locations (vehicle_id, recorded_at DESC);

COMMIT;