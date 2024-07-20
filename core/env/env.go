package env

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// These values will be replaced during the build process using build arguments
var (
	SHORT_COMMIT = ""
	LONG_COMMIT  = ""
	VERSION      = ""
	BUILD_TIME   = ""
)

const (
	DEV        = "DEV"
	PRODUCTION = "PRODUCTION"
	STAGING    = "STAGING"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	Port        string `env:"PORT"`
	Version     string `env:"VERSION"`
	Environment string `env:"ENVIRONMENT"`

	EnableStartupMessage bool `env:"ENABLE_STARTUP_MESSAGE"`
	EnablePrintRoutes    bool `env:"ENABLE_PRINT_ROUTES"`
	EnableStackTrace     bool `env:"ENABLE_STACK_TRACE"`

	SkipMigration bool `env:"SKIP_MIGRATION"`

	SendgridAPIKey   string `env:"SENDGRID_API_KEY"`
	SMTPURL          string `env:"SMTP_URL"`
	SMTPPort         string `env:"SMTP_PORT"`
	SMTPClientID     string `env:"SMTP_CLIENT_ID"`
	SMTPClientSecret string `env:"SMTP_CLIENT_SECRET"`
	EmailSender      string `env:"EMAIL_SENDER"`
	EmailSenderName  string `env:"EMAIL_SENDER_NAME"`

	JWTSecret string `env:"JWT_SECRET"`

	RedisURL string `env:"REDIS_URL"`
}

// Load Loads .env files and/or values from enviroment variables
//
// Note: It WILL NOT OVERRIDE an env variable that already exists - consider the .env file to set dev vars or sensible defaults
func Load(basePathOverwrite ...string) (Config, error) {

	// Set default values here
	config := Config{
		Environment:          DEV,
		EnableStartupMessage: true,
		EnablePrintRoutes:    false,
		EnableStackTrace:     false,
		SkipMigration:        false,
		EmailSender:          "no-reply@golangboilerplate.com",
		EmailSenderName:      "Golang Boilerplate (No Reply)",
	}

	environ := strings.ToUpper(Get("ENVIRONMENT", config.Environment))

	basePath := path.Join(".", "configs")
	if len(basePathOverwrite) > 0 {
		basePath = basePathOverwrite[0]
	}

	if environ == DEV {
		godotenv.Load(path.Join(basePath, ".env"))
	} else {
		filename := fmt.Sprintf(".%s.env", environ)
		godotenv.Load(path.Join(basePath, filename))
	}

	err := env.Parse(&config)

	// Overwrite since Version is setted by ldflags
	// E.g: go build -ldflags="-X 'project/core/env.VERSION=<desired_version>'" project/main.go
	config.Version = VERSION

	return config, err
}

// Get gets environment variable or the default value
func Get(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
