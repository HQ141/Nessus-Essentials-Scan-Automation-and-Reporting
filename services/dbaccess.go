package services

import (
	"NessusEssAutomation/config"
	"NessusEssAutomation/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func dbaccess(vulnerability models.Vulnerability) (models.Vulnerability, error) {
	// urlExample := "postgres://username:password@url/database_name"
	conf, err := config.LoadDBConfig()
	if err != nil {
		log.Fatal("Error Loading variables")
		return models.Vulnerability{}, err
	}
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s", conf.Username, conf.Password, conf.Url))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	present := bool(false)
	err = conn.QueryRow(context.Background(), fmt.Sprintf("SELECT EXISTS(select * from vulnerabilities where plugin_id='%s');", vulnerability.PluginID)).Scan(&present)
	if err != nil {
		log.Fatal("failed to query database")
		return models.Vulnerability{}, err
	}
	if !present {
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
		vulnerability.VulnerabilityFirstDiscoveredDate = time.Now().Local().Format("2006-01-02")
		vulnerability.VulnerabilityLastObservedDate = time.Now().Local().Format("2006-01-02")
		// Execute the query using the vulnerability values
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
			time.Now().Local().Format("2006-01-02"), //vulnerability.PatchPublicationDate,
			vulnerability.Exploitability,
		)
		if err != nil {
			log.Fatal("cant insert record")
			return models.Vulnerability{}, err
		}
	} else {
		// err = conn.QueryRow(context.Background(), fmt.Sprintf("select vulnerability_last_observed_date,vulnerability_first_discovered_date from vulnerabilities where plugin_id='%s';", vulnerability.PluginID)).Scan(&vulnerability.VulnerabilityLastObservedDate, &vulnerability.VulnerabilityFirstDiscoveredDate)
		var dd1 pgtype.Date
		var dt2 pgtype.Date
		err = conn.QueryRow(context.Background(), "select vulnerability_last_observed_date,vulnerability_first_discovered_date from vulnerabilities where plugin_id=$1;", vulnerability.PluginID).Scan(&dd1, &dt2)
		if err != nil {
			return models.Vulnerability{}, err
		}
		vulnerability.VulnerabilityLastObservedDate = dd1.Time.Format("2006-01-02")
		vulnerability.VulnerabilityFirstDiscoveredDate = dt2.Time.Format("2006-01-02")
		query := `
		UPDATE vulnerabilities
		SET vulnerability_last_observed_date = $1 where plugin_id =$2`
		_, err = conn.Exec(context.Background(), query, time.Now().Local().Format("2006-01-02"), vulnerability.PluginID)
	}

	defer conn.Close(context.Background())
	if err != nil {
		log.Println("Error inserting vulnerability:", err)
		return models.Vulnerability{}, err
	}
	return vulnerability, nil
}
