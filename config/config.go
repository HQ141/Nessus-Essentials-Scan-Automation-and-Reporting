package config

import (
	"NessusEssAutomation/models"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds the configuration values
type Config struct {
	Username string
	Password string
	Url      string
	ScanName string
}

// LoadConfig loads environment variables into the Config struct
func LoadConfig() (models.Nessus, error) {
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
