package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	PublicBaseURL string
	BarcodeDir    string
	BarcodeWidth  int
	BarcodeHeight int
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	DBTimeZone    string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		PublicBaseURL: strings.TrimRight(getEnv("PUBLIC_BASE_URL", "http://localhost:8080"), "/"),
		BarcodeDir:    getEnv("BARCODE_DIR", "images/barcodes"),
		BarcodeWidth:  getEnvInt("BARCODE_WIDTH", 640),
		BarcodeHeight: getEnvInt("BARCODE_HEIGHT", 220),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "barcode_db"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:    getEnv("DB_TIMEZONE", "Asia/Bangkok"),
	}

	return cfg, nil
}

func (c Config) DBDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode +
		" TimeZone=" + c.DBTimeZone
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return intValue
}
