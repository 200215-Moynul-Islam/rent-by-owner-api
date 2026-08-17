CREATE INDEX idx_destinations_city_trgm
ON destinations USING GIN (city gin_trgm_ops);

CREATE INDEX idx_destinations_country_trgm
ON destinations USING GIN (country gin_trgm_ops);