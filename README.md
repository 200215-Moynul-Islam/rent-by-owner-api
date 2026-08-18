# Rent by Owner API

A Go/Beego v2 REST API for searching rental destinations, backed by PostgreSQL + PostGIS for geospatial queries and `pg_trgm` for fuzzy text search.

## Tech Stack

- **Go 1.26** with **Beego v2** (ORM, routing, controllers)
- **PostgreSQL 17** + **PostGIS 3.5** for geography queries
- **pg_trgm** for trigram-based fuzzy matching
- **golang-migrate** for schema migrations
- Docker Compose for local orchestration

## Features

- **Search** — fuzzy match on city/country name (ILIKE + trigram similarity), re-ranked in-app by exact/prefix/partial match and population
- **Autocomplete** — same matching, DB-side ranked and limited to top 5 results
- **Nearby** — destinations within a radius (km) of a lat/lon point, sorted by distance
- **Within bounds** — destinations inside a lat/lon bounding box
- **CSV import tool** — bulk-loads destinations via `COPY` into a staging table, then upserts into `destinations` with PostGIS geography points

## Getting Started

```bash
git clone https://github.com/200215-Moynul-Islam/rent-by-owner-api.git
cd rent-by-owner-api
docker compose up --build
```

This starts:

- `postgres` — PostGIS-enabled database (health-checked before dependents start)
- `migrate` — runs all migrations in `./migrations` on startup
- `api` — the Beego app, listening on `:8080`

### Local (non-Docker) setup

1. Copy `conf/app.conf.sample` to `conf/app.conf` and adjust the `[database]` section.
2. Run migrations against your Postgres instance (see `./migrations`).
3. Run the app with the [bee CLI](https://github.com/beego/bee):

   ```bash
   go install github.com/beego/bee/v2@latest
   bee run
   ```

   `bee run` watches for file changes and rebuilds/restarts automatically — useful while developing. For a plain one-off run, `go run main.go` also works.

### Importing destination data

Place a CSV at `data/destinations.csv` with header `Country,City,Population,Latitude,Longitude` (`Population` may be blank), then:

```bash
go run ./cmd/import-destinations [path/to/file.csv]
```

Configurable via `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_SSLMODE` env vars (defaults match `docker-compose.yml`).

## API Reference

Base path: `/api/v1`

| Method | Endpoint                      | Query Params                     | Description                                    |
| ------ | ----------------------------- | -------------------------------- | ---------------------------------------------- |
| GET    | `/health`                     | —                                | Service health check                           |
| GET    | `/destinations/search`        | `q` (required)                   | Fuzzy search by city/country                   |
| GET    | `/destinations/autocomplete`  | `q` (required)                   | Top 5 ranked suggestions                       |
| GET    | `/destinations/nearby`        | `lat`, `lon`, `radius` (km)      | Destinations within radius, sorted by distance |
| GET    | `/destinations/within-bounds` | `north`, `south`, `east`, `west` | Destinations inside a bounding box             |

All responses follow:

```json
{
  "success": true,
  "message": "Destinations retrieved successfully.",
  "data": [ ... ]
}
```

## Project Structure

```
controllers/    # HTTP handlers
services/       # Business logic
repositories/   # Data access (raw SQL via Beego ORM)
dtos/           # Response shapes
utils/          # Response helpers, validation, in-app ranking
database/       # DB init/connection
routers/        # Route registration
migrations/     # golang-migrate SQL migrations
cmd/import-destinations/  # standalone CSV import tool
```
