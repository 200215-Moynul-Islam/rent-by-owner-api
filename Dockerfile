# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rent-by-owner-api .
RUN go build -o import-destinations ./cmd/import-destinations


# Runtime stage
FROM alpine:3.22

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/rent-by-owner-api .
COPY --from=builder --chown=appuser:appgroup /app/import-destinations .
COPY --from=builder --chown=appuser:appgroup /app/conf ./conf

USER appuser

EXPOSE 8080

CMD ["./rent-by-owner-api"]