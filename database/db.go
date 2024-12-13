package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/davidalvarez305/yd_cocktails/constants"
	_ "github.com/lib/pq"
)

type connection struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

var DB *sql.DB

func Connect() (*sql.DB, error) {
	conn := connection{
		host:     constants.PostgresHost,
		port:     constants.PostgresPort,
		user:     constants.PostgresUser,
		password: constants.PostgresPassword,
		dbName:   constants.PostgresDBName,
	}

	connectionString := connToString(conn)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	MaxOpenConnectionsStr := constants.MaxOpenConnections
	MaxIdleConnectionsStr := constants.MaxIdleConnections
	MaxConnectionLifetimeStr := constants.MaxConnectionLifetime

	MaxOpenConnections, err := strconv.Atoi(MaxOpenConnectionsStr)
	if err != nil {
		fmt.Printf("Error parsing MAX_OPEN_CONNECTIONS: %v\n", err)
		return nil, fmt.Errorf("error parsing max connections: %v", err)
	}

	MaxIdleConnections, err := strconv.Atoi(MaxIdleConnectionsStr)
	if err != nil {
		fmt.Printf("Error parsing MAX_IDLE_CONNECTIONS: %v\n", err)
		return nil, fmt.Errorf("error parsing idle connections: %v", err)
	}

	// Assuming the connection lifetime is in seconds, parse it to int and convert to time.Duration
	MaxConnectionLifetimeSeconds, err := strconv.Atoi(MaxConnectionLifetimeStr)
	if err != nil {
		fmt.Printf("Error parsing MAX_CONN_LIFETIME: %v\n", err)
		return nil, fmt.Errorf("error parsing max connection lfetime: %v", err)
	}
	MaxConnectionLifetime := time.Duration(MaxConnectionLifetimeSeconds) * time.Second

	// Set connection pool parameters
	db.SetMaxOpenConns(MaxOpenConnections)
	db.SetMaxIdleConns(MaxIdleConnections)
	db.SetConnMaxLifetime(MaxConnectionLifetime)

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to verify connection: %v", err)
	}

	DB = db
	return db, nil
}

func connToString(info connection) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		info.user, info.password, info.host, info.port, info.dbName)
}
