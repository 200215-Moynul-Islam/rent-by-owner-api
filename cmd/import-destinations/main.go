package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const defaultCSVPath = "data/destinations.csv"

type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
}

type DestinationRecord struct {
	Country    string
	City       string
	Population any
	Latitude   float64
	Longitude  float64
}

func main() {
	config := loadConfig()
	csvPath := getCSVPath()

	db, err := openDatabase(config)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	log.Printf("Importing destinations from %s", csvPath)

	if err := importDestinations(db, csvPath); err != nil {
		log.Fatalf("destination import failed: %v", err)
	}

	log.Println("destination import completed successfully")
}

func loadConfig() Config {
	return Config{
		Host:     getEnv("POSTGRES_HOST", "postgres"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		Database: getEnv("POSTGRES_DB", "rent_by_owner"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}
}

func (c Config) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Database,
		c.User,
		c.Password,
		c.SSLMode,
	)
}

func getCSVPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	return defaultCSVPath
}

func openDatabase(config Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.dsn())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func importDestinations(db *sql.DB, csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("open CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if err := validateHeader(reader); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := createStagingTable(tx); err != nil {
		return err
	}

	stmt, err := prepareCopy(tx)
	if err != nil {
		return err
	}

	count, err := copyCSVRows(reader, stmt)
	if err != nil {
		stmt.Close()
		return err
	}

	if err := finishCopy(stmt); err != nil {
		return err
	}

	if err := insertDestinations(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	log.Printf("imported %d destinations", count)

	return nil
}

func validateHeader(reader *csv.Reader) error {
	reader.FieldsPerRecord = 5

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read CSV header: %w", err)
	}

	expected := []string{
		"Country",
		"City",
		"Population",
		"Latitude",
		"Longitude",
	}

	if len(header) != len(expected) {
		return fmt.Errorf("invalid CSV header")
	}

	for i := range expected {
		if strings.TrimSpace(header[i]) != expected[i] {
			return fmt.Errorf(
				"invalid CSV header: expected %q at column %d",
				expected[i],
				i+1,
			)
		}
	}

	return nil
}

func createStagingTable(tx *sql.Tx) error {
	const query = `
		CREATE TEMP TABLE destination_import (
			country TEXT NOT NULL,
			city TEXT NOT NULL,
			population INTEGER,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL
		)
	`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("create staging table: %w", err)
	}

	return nil
}

func prepareCopy(tx *sql.Tx) (*sql.Stmt, error) {
	const query = `
		COPY destination_import (
			country,
			city,
			population,
			latitude,
			longitude
		) FROM STDIN
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare COPY: %w", err)
	}

	return stmt, nil
}

func copyCSVRows(reader *csv.Reader, stmt *sql.Stmt) (int, error) {
	count := 0

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, fmt.Errorf(
				"read CSV row %d: %w",
				count+2,
				err,
			)
		}

		destination, err := parseDestination(record)
		if err != nil {
			return 0, fmt.Errorf(
				"parse CSV row %d: %w",
				count+2,
				err,
			)
		}

		if err := copyDestination(stmt, destination); err != nil {
			return 0, fmt.Errorf(
				"COPY row %d: %w",
				count+2,
				err,
			)
		}

		count++
	}

	return count, nil
}

func parseDestination(record []string) (DestinationRecord, error) {
	population, err := parsePopulation(record[2])
	if err != nil {
		return DestinationRecord{}, fmt.Errorf(
			"invalid population: %w",
			err,
		)
	}

	latitude, err := parseCoordinate(record[3], "latitude")
	if err != nil {
		return DestinationRecord{}, err
	}

	longitude, err := parseCoordinate(record[4], "longitude")
	if err != nil {
		return DestinationRecord{}, err
	}

	if err := validateCoordinates(latitude, longitude); err != nil {
		return DestinationRecord{}, err
	}

	return DestinationRecord{
		Country:    strings.TrimSpace(record[0]),
		City:       strings.TrimSpace(record[1]),
		Population: population,
		Latitude:   latitude,
		Longitude:  longitude,
	}, nil
}

func parsePopulation(value string) (any, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil, nil
	}

	population, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}

	if population < 0 {
		return nil, fmt.Errorf("population cannot be negative")
	}

	return population, nil
}

func parseCoordinate(value, name string) (float64, error) {
	coordinate, err := strconv.ParseFloat(
		strings.TrimSpace(value),
		64,
	)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, err)
	}

	return coordinate, nil
}

func validateCoordinates(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf(
			"latitude out of range: %f",
			latitude,
		)
	}

	if longitude < -180 || longitude > 180 {
		return fmt.Errorf(
			"longitude out of range: %f",
			longitude,
		)
	}

	return nil
}

func copyDestination(
	stmt *sql.Stmt,
	destination DestinationRecord,
) error {
	_, err := stmt.Exec(
		destination.Country,
		destination.City,
		destination.Population,
		destination.Latitude,
		destination.Longitude,
	)

	return err
}

func finishCopy(stmt *sql.Stmt) error {
	if _, err := stmt.Exec(); err != nil {
		stmt.Close()
		return fmt.Errorf("finish COPY: %w", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("close COPY: %w", err)
	}

	return nil
}

func insertDestinations(tx *sql.Tx) error {
	const query = `
		INSERT INTO destinations (
			country,
			city,
			population,
			location
		)
		SELECT
			country,
			city,
			population,
			ST_SetSRID(
				ST_MakePoint(longitude, latitude),
				4326
			)::geography
		FROM destination_import
	`

	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf(
			"insert destinations: %w",
			err,
		)
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))

	if value == "" {
		return fallback
	}

	return value
}