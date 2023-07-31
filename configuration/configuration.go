package configuration

import (
	"os"
	"log"
	"fmt"
	"path/filepath"
    
	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
    Port     string
	Database string
    Username string
    Password string
	MaxConnections int
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
func LoadEnv() Config {
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
func ConnectURL(databaseConfig Config) string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        databaseConfig.Username,
        databaseConfig.Password,
        databaseConfig.Host,
        databaseConfig.Port,
        databaseConfig.Database,
    )
}
