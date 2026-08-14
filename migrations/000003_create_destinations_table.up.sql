CREATE TABLE destinations (
    id BIGSERIAL PRIMARY KEY,
    country VARCHAR(100) NOT NULL,
    city VARCHAR(255) NOT NULL,
    population INTEGER,
    location GEOGRAPHY(POINT, 4326) NOT NULL
);