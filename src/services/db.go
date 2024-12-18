package services

import (
	"NessusEssAutomation/config"
	"NessusEssAutomation/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// dbaccess interacts with the database to insert or update vulnerability data.
// If the vulnerability is not already in the database, it inserts a new record.
// If the vulnerability exists, it updates the `vulnerability_last_observed_date`.
func dbaccess(vulnerability models.Vulnerability) (models.Vulnerability, error) {
	// Load database configuration from environment variables or config file
	conf, err := config.LoadDBConfig()
	if err != nil {
		log.Fatal("Error loading database configuration")
		return models.Vulnerability{}, err
	}

	// Establish a connection to the PostgreSQL database
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s", conf.Username, conf.Password, conf.Url))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return models.Vulnerability{}, err
	}

	// Check if the vulnerability already exists in the database
	present := bool(false)
	err = conn.QueryRow(context.Background(), fmt.Sprintf("SELECT EXISTS(select * from vulnerabilities where plugin_id='%s');", vulnerability.PluginID)).Scan(&present)
	if err != nil {
		log.Fatal("Failed to query the database for vulnerability existence")
		return models.Vulnerability{}, err
	}

	if !present {
		// Vulnerability does not exist, prepare to insert a new record
		query := `
		INSERT INTO vulnerabilities (
			plugin_id, asset_name, plugin_name, severity, cvss3_base_score,
			vulnerability_priority_rating, plugin_output, synopsis, description,
			solution, vulnerability_first_discovered_date, vulnerability_last_observed_date,
			plugin_publication_date, patch_publication_date, exploitability
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12,
			$13, $14, $15
		)`

		// Set the discovered and observed dates to the current date
		vulnerability.VulnerabilityFirstDiscoveredDate = time.Now().Local().Format("2006-01-02")
		vulnerability.VulnerabilityLastObservedDate = time.Now().Local().Format("2006-01-02")

		// Execute the query to insert the vulnerability
		_, err = conn.Exec(context.Background(), query,
			vulnerability.PluginID,
			vulnerability.AssetName,
			vulnerability.PluginName,
			vulnerability.Severity,
			vulnerability.CVSS3BaseScore,
			vulnerability.VulnerabilityPriorityRating,
			vulnerability.PluginOutput,
			vulnerability.Synopsis,
			vulnerability.Description,
			vulnerability.Solution,
			vulnerability.VulnerabilityFirstDiscoveredDate,
			vulnerability.VulnerabilityLastObservedDate,
			vulnerability.PluginPublicationDate,
			time.Now().Local().Format("2006-01-02"), // Patch publication date
			vulnerability.Exploitability,
		)
		if err != nil {
			log.Fatal("Failed to insert vulnerability record")
			return models.Vulnerability{}, err
		}
	} else {
		// Vulnerability exists, fetch the last observed and first discovered dates
		var dd1 pgtype.Date
		var dt2 pgtype.Date
		err = conn.QueryRow(context.Background(), "SELECT vulnerability_last_observed_date, vulnerability_first_discovered_date FROM vulnerabilities WHERE plugin_id=$1;", vulnerability.PluginID).Scan(&dd1, &dt2)
		if err != nil {
			log.Println("Failed to fetch existing vulnerability dates")
			return models.Vulnerability{}, err
		}

		// Update the vulnerability struct with fetched dates
		vulnerability.VulnerabilityLastObservedDate = dd1.Time.Format("2006-01-02")
		vulnerability.VulnerabilityFirstDiscoveredDate = dt2.Time.Format("2006-01-02")

		// Update the `vulnerability_last_observed_date` in the database
		query := `
		UPDATE vulnerabilities
		SET vulnerability_last_observed_date = $1
		WHERE plugin_id = $2`
		_, err = conn.Exec(context.Background(), query, time.Now().Local().Format("2006-01-02"), vulnerability.PluginID)
	}

	// Ensure the database connection is closed
	defer conn.Close(context.Background())

	// Check if there was an error during the process
	if err != nil {
		log.Println("Error processing vulnerability:", err)
		return models.Vulnerability{}, err
	}

	// Return the updated vulnerability struct
	return vulnerability, nil
}
