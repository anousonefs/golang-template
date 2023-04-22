package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	PO_ADMIN     = "PO_ADMIN"
	VENDOR_ADMIN = "VENDOR_ADMIN"
	BRANCH_ADMIN = "BRANCH_ADMIN"
	BRANCH_STAFF = "BRANCH_STAFF"
)

const (
	REGISTER = "REGISTER"
	SENDING  = "SENDING"
	ARRIVE   = "ARRIVE"
	RECEIVE  = "RECEIVE"
)

const (
	DAY   = "DAY"
	MONTH = "MONTH"
	YEAR  = "YEAR"
)

// Config store all configurable of application.
// The values are read by Viper from a config file or environment variables.
type Config struct {
	// DBDriver is the database driver.
	DBDriver string
	// dbHost is the database host.
	dbHost string
	// dbPort is the database port.
	dbPort string
	// dbUser is the database user.
	dbUser string
	// password is the database password.
	password string
	// dbName is the database name.
	dbName   string
	AssetDir string
	BaseUrl  string

	// Port is the port of application.
	Port string

	// PasetoSecret is the secret key of paseto
	PasetoSecret []byte

	// OneSignalApiKey is the private key of onesignal.
	OneSignalApiKey string
	// OneSignalAppID is the public key of onesignal.
	OneSignalAppID string

	TwilioAccountID  string
	TwilioAuthToken  string
	TwilioServiceSID string
}

// DSNInfo is database connection for postgresql
func (c Config) DSNInfo() string {
	timeoutOption := fmt.Sprintf("-c statement_timeout=%d", 10*time.Minute/time.Millisecond)
	return fmt.Sprintf("user='%s' password='%s' host='%s' port=%s dbname='%s' sslmode=disable options='%s'", c.dbUser, c.password, c.dbHost, c.dbPort, c.dbName, timeoutOption)
}

// GetEnv returns the environment variable or default value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Init load config from file or environment variables.
func NewConfig() (config Config, err error) {
	config.DBDriver = GetEnv("DB_DRIVER", "postgres")
	config.dbHost = GetEnv("PGHOST", "127.0.0.1")
	config.dbPort = GetEnv("PGPORT", "5432")
	config.dbUser = os.Getenv("PGUSER")
	config.dbUser = "anousone"
	config.password = os.Getenv("PGSECRET")
	config.dbName = os.Getenv("PGDATABASE")
	config.dbName = "clean_architecture"
	config.BaseUrl = os.Getenv("BASE_URL")
	config.TwilioServiceSID = os.Getenv("TWILIO_SERVICE_SID")
	config.TwilioAuthToken = os.Getenv("TWILIO_AUTHTOKEN")
	config.TwilioAccountID = os.Getenv("TWILIO_ACCOUNT_ID")

	if config.TwilioAccountID == "" || config.TwilioAuthToken == "" || config.TwilioServiceSID == "" {
		return config, errors.New("Twilio config is empty")
	}

	if config.BaseUrl == "" {
		return config, errors.New("BASE_URL is empty")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	config.AssetDir = GetEnv("ASSET_DIR", homeDir)
	config.Port = GetEnv("PORT", "8080")

	config.PasetoSecret, err = hex.DecodeString(os.Getenv("PASETO_SECRET"))
	if len(config.PasetoSecret) != 32 {
		return config, err
	}

	config.OneSignalApiKey = os.Getenv("RESTAPIKEY")
	config.OneSignalAppID = os.Getenv("APPIDKEY")
	return
}

// Psql returns squirrel StatementBuilderType for PostgreSQL.
func Psql() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

// Mysql returns squirrel StatementBuilderType for MySQL.
func Mysql() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
}
