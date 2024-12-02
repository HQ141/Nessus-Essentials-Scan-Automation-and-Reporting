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
	USERNAME, _username := os.LookupEnv("Nessus_Username")
	PASSWORD, _password := os.LookupEnv("Nessus_Password")
	URL, _url := os.LookupEnv("Nessus_Url")
	ScanName, _scanname := os.LookupEnv("Nessus_Scanname")
	if !_password || !_username || !_url || !_scanname {
		return models.Nessus{}, errors.New("environment variables not set")
	}

	config := models.Nessus{Username: USERNAME, Password: PASSWORD, Url: URL, ScanName: ScanName}
	return config, nil
}

// LoadConfig loads environment variables into the Config struct
func LoadDBConfig() (models.DB, error) {
	godotenv.Load(".env")
	USERNAME, _username := os.LookupEnv("DB_Username")
	PASSWORD, _password := os.LookupEnv("DB_Password")
	DATABASEURL, _databaseurl := os.LookupEnv("DATABASE_URL")
	if !_password || !_username || !_databaseurl {
		return models.DB{}, errors.New("environment variables not set")
	}

	config := models.DB{Username: USERNAME, Password: PASSWORD, Url: DATABASEURL}
	return config, nil
}
