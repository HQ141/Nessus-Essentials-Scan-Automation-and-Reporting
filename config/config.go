package config

import (
	"NessusEssAutomation/models"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables into the Config struct
func LoadNessusConfig() (models.Nessus, error) {
	godotenv.Load(".env")
	USERNAME, _username := os.LookupEnv("NESSUS_USERNAME")
	PASSWORD, _password := os.LookupEnv("NESSUS_PASSWORD")
	URL, _url := os.LookupEnv("NESSUS_URL")
	ScanName, _scanname := os.LookupEnv("NESSUS_SCANNAME")
	if !_password || !_username || !_url || !_scanname {
		return models.Nessus{}, errors.New("nessus environment variables not set")
	}

	config := models.Nessus{Username: USERNAME, Password: PASSWORD, Url: URL, ScanName: ScanName}
	return config, nil
}

// LoadConfig loads environment variables into the Config struct
func LoadSMTPConfig() (models.SMTP, error) {
	godotenv.Load(".env")
	USERNAME, _username := os.LookupEnv("SMTP_USERNAME")
	PASSWORD, _password := os.LookupEnv("SMTP_PASSWORD")
	URL, _url := os.LookupEnv("SMTP_URL")
	RECEIVERS, _receivers := os.LookupEnv("RECEIVERS")
	if !_password || !_username || !_url || !_receivers {
		return models.SMTP{}, errors.New("SMTP environment variables not set")
	}
	config := models.SMTP{Username: USERNAME, Password: PASSWORD, Url: URL, Recipients: RECEIVERS}
	return config, nil
}

// LoadConfig loads environment variables into the Config struct
func LoadDBConfig() (models.DB, error) {
	godotenv.Load(".env")
	USERNAME, _username := os.LookupEnv("POSTGRES_USER")
	PASSWORD, _password := os.LookupEnv("POSTGRES_PASSWORD")
	DATABASEURL, _databaseurl := os.LookupEnv("POSTGRES_URL")
	if !_password || !_username || !_databaseurl {
		return models.DB{}, errors.New("Postgress environment variables not set")
	}

	config := models.DB{Username: USERNAME, Password: PASSWORD, Url: DATABASEURL}
	return config, nil
}
