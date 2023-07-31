package database

import (
	"os"
	"log"
	"context"
    
	"github.com/jackc/pgx/v5/pgxpool"
    "pokeverse/web-service/configuration"
)

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
    config := configuration.LoadEnv()
    connURL := configuration.ConnectURL(config)

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
func Query(pool *pgxpool.Pool, queryString string, args ...interface{}) ([]byte, error) {
    var data []byte

    rows, err := pool.Query(context.Background(), queryString, args...)
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
