CREATE INDEX idx_destinations_location_gist
ON destinations
USING GIST (location);