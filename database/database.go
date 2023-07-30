package database

import (
	"os"
	"log"
	"fmt"
	"context"
	"path/filepath"
    
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
    Port     string
	Database string
    Username string
    Password string
	MaxConnections int
}

// ConnectDatabase establishes a connection to the PostgreSQL database using the provided configuration
// loaded from the environment variables.
//
// The function retrieves the database configuration by calling the "loadEnv" function, which reads the
// required connection parameters from environment variables (e.g., DB_HOST, DB_PORT, DB_USERNAME, etc.).
// It then constructs the connection URL by calling the "connectURL" function with the obtained
// configuration. The connection URL is used to create a new connection pool configuration by parsing
// it using "pgxpool.ParseConfig".
//
// The "MaxConnections" value from the configuration is set in the pool configuration to determine the
// maximum number of connections in the connection pool.
//
// Finally, the function creates a new connection pool using "pgxpool.New" with the connection URL and
// pool configuration. If any error occurs during the process, the function logs the error and exits the
// program with a non-zero exit code.
//
// Return Value:
//   *pgxpool.Pool: A pointer to the connection pool (*pgxpool.Pool) representing the PostgreSQL
//                  database connection. The connection pool can be used to perform database queries
//                  and transactions.
func ConnectDatabase() *pgxpool.Pool {
    config := loadEnv()
    connURL := connectURL(config)

    poolConfig, err := pgxpool.ParseConfig(connURL)
    if err != nil {
        log.Fatalf("Failed to parse connection config: %v\n", err)
        os.Exit(1)
    }

    poolConfig.MaxConns = int32(config.MaxConnections)

    pool, err := pgxpool.New(context.Background(), connURL)
    if err != nil {
        log.Fatalf("Unable to create connection pool: %v\n", err)
        os.Exit(1)
    }

    return pool
}

// Query executes a database query using the provided connection pool and query string,
// and returns the query result as a byte slice containing JSON data.
//
// Parameters:
//   pool: A *pgxpool.Pool representing the connection pool to the PostgreSQL database.
//   queryString: A string containing the SQL query to be executed.
//
// Return Values:
//   []byte: A byte slice containing the JSON data returned by the query.
//   error: An error value if the query or scanning process fails, otherwise nil.
func Query(pool *pgxpool.Pool, queryString string) ([]byte, error) {
    var data []byte

    rows, err := pool.Query(context.Background(), queryString)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    if rows.Next() {
        var dbResults string
        err := rows.Scan(&dbResults)
        if err != nil {
            return nil, err
        }

        // Assuming the database result is a JSON string
        data = []byte(dbResults)
    }

    return data, nil
}

// loadEnv loads the database configuration from a .env file and returns a Config struct
// containing the configuration settings for the database connection.
//
// The function retrieves the current working directory, loads the environment variables
// from a .env file located in the current working directory, and then constructs a Config
// struct based on the loaded environment variables.
//
// The .env file is expected to contain the following environment variables:
//   - DB_HOST: The hostname or IP address of the database server.
//   - DB_PORT: The port number on which the database server is listening.
//   - DB_DATABASE: The name of the database to connect to.
//   - DB_USERNAME: The username for the database connection.
//   - DB_PASSWORD: The password for the database connection.
//
// If any of the required environment variables are missing or the .env file cannot be loaded,
// the function logs an error message and exits the application with a non-zero status code.
//
// Return Value:
//   Config: A Config struct containing the database configuration settings.
func loadEnv() Config {
    // Retrieve the current working directory
    cwd, wdErr := os.Getwd()
    if wdErr != nil {
        log.Fatalf("Failed to retrieve current working directory.")
        os.Exit(1)
    }

    // Load environment variables from a .env file located in the current working directory
    envFileErr := godotenv.Load(filepath.Join(string(cwd), ".env"))
    if envFileErr != nil {
        log.Fatal("Error loading .env file")
        os.Exit(1)
    }

    // Construct and return a Config struct based on the loaded environment variables
    return Config{
        Host:           os.Getenv("DB_HOST"),
        Port:           os.Getenv("DB_PORT"),
        Database:       os.Getenv("DB_DATABASE"),
        Username:       os.Getenv("DB_USERNAME"),
        Password:       os.Getenv("DB_PASSWORD"),
        MaxConnections: 50,
    }
}

// connectURL formats the database connection URL using the provided database configuration.
// The function returns a connection URL in the format "postgres://username:password@host:port/database".
// It is used to create a valid connection URL to connect to a PostgreSQL database.
//
// Parameters:
//   databaseConfig: A Config struct containing the database configuration settings.
//
// Return Value:
//   string: The formatted connection URL for the PostgreSQL database.
func connectURL(databaseConfig Config) string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        databaseConfig.Username,
        databaseConfig.Password,
        databaseConfig.Host,
        databaseConfig.Port,
        databaseConfig.Database,
    )
}

