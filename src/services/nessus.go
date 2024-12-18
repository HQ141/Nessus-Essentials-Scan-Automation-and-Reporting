package services

import (
	"NessusEssAutomation/config"
	"NessusEssAutomation/models"
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

// getScanItems retrieves the scan report from Nessus Essentials in CSV format.
// If `last` is true, it ensures the scan has been performed in the last 12 hours.
func getScanItems(context models.Nessus, last bool) ([][]string, error) {
	myscan := models.Scan{}

	// Retrieve all scans
	Scans, err := getScan(context)
	if err != nil {
		return nil, err
	}

	// Find the scan matching the specified name
	for _, i := range Scans {
		if i.Name == context.ScanName {
			myscan = i
		}
	}

	// If no matching scan is found, return an error
	if myscan == (models.Scan{}) {
		return nil, errors.New("no matching scan name found")
	}

	// Check if the scan is recent (within 12 hours)
	if !last && time.Since(time.Unix(myscan.LastModificationDate, 0)) >= 12*time.Hour {
		return [][]string{}, errors.New("no new scan")
	}

	// Prepare the request payload to export the scan report in CSV format
	values := map[string]interface{}{
		"format":      "csv",
		"template_id": "",
		"reportContents": map[string]interface{}{
			"csvColumns": map[string]bool{
				"id":                 true,
				"cve":                true,
				"cvss":               true,
				"risk":               true,
				"hostname":           true,
				"protocol":           true,
				"port":               true,
				"plugin_name":        true,
				"synopsis":           true,
				"description":        true,
				"solution":           true,
				"see_also":           true,
				"plugin_output":      true,
				"risk_factor":        true,
				"plugin_information": true,
				"exploitable_with":   true,
			},
		},
	}

	// Marshal the payload into JSON format
	jsonValues, _ := json.Marshal(values)

	// Create an HTTP request to initiate the report export
	r, err := http.NewRequest("POST", context.Url+fmt.Sprintf("/scans/%d/export", myscan.ID), bytes.NewBuffer(jsonValues))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", context.AccessKey, context.SecretKey))

	// Execute the HTTP request
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}

	// Handle HTTP status codes
	if res.StatusCode == 409 {
		return [][]string{}, errors.New("no new scan")
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("http request error, received %d", res.StatusCode)
	}

	// Parse the response to get the export token
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	token, err := jsonparser.GetString(body, "token")
	if err != nil {
		return nil, err
	}

	// Poll the export endpoint until the report is ready
	for {
		time.Sleep(5 * time.Second)

		req, err := http.NewRequest("GET", context.Url+fmt.Sprintf("/tokens/%s/download", token), nil)
		req.Header.Set("Accept", "text/csv")
		req.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", context.AccessKey, context.SecretKey))

		if err != nil {
			return nil, err
		}

		// Execute the HTTP request to download the report
		resp, err := context.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		// If the report is not ready, retry
		if resp.StatusCode == http.StatusNotFound {
			continue
		}

		// Read the report content as a CSV
		report, _ := io.ReadAll(resp.Body)
		rec := csv.NewReader(strings.NewReader(string(report)))
		records, _ := rec.ReadAll()

		// Filter records to exclude rows with "None" risk
		fields := make(map[string]int)
		for i, name := range records[0] {
			fields[name] = i
		}
		var filteredRecords [][]string
		filteredRecords = append(filteredRecords, records[0])
		for _, record := range records[1:] {
			if record[fields["Risk"]] != "None" {
				filteredRecords = append(filteredRecords, record)
			}
		}
		return filteredRecords, nil
	}
}

// authenticate establishes a connection to Nessus and retrieves access and secret keys.
func authenticate() (models.Nessus, error) {
	// Load Nessus configuration
	context, err := config.LoadNessusConfig()
	if err != nil {
		return models.Nessus{}, err
	}

	// Configure HTTP client to skip TLS verification
	customTransport := (http.DefaultTransport.(*http.Transport))
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	context.HttpClient = &http.Client{Transport: customTransport}

	// Retrieve authentication token
	context, err = getToken(context)
	if err != nil {
		return models.Nessus{}, err
	}

	// Retrieve access and secret keys
	keys, err := getSessionKeys(context)
	if err != nil {
		return models.Nessus{}, err
	}
	return keys, nil
}

// getToken authenticates with Nessus and retrieves a session token.
func getToken(context models.Nessus) (models.Nessus, error) {
	body := []byte(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, context.Username, context.Password))
	posturl := context.Url + "/session"

	// Create an HTTP request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		return models.Nessus{}, err
	}
	r.Header.Add("Content-Type", "application/json")

	// Execute the request
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return models.Nessus{}, err
	}
	defer res.Body.Close()

	// Parse the response to extract the token
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Nessus{}, err
	}
	var responseobj models.Session
	err = json.Unmarshal(bodyBytes, &responseobj)
	if err != nil {
		return models.Nessus{}, err
	}
	context.Token = responseobj.Token
	return context, nil
}

// getSessionKeys retrieves access and secret keys from Nessus.
func getSessionKeys(context models.Nessus) (models.Nessus, error) {
	puturl := context.Url + "/session/keys"

	// Create an HTTP request
	r, err := http.NewRequest("PUT", puturl, nil)
	if err != nil {
		return models.Nessus{}, err
	}
	token := fmt.Sprintf("token=%s", context.Token)
	r.Header.Add("X-Cookie", token)

	// Execute the request
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return models.Nessus{}, err
	}
	if res.StatusCode != 200 {
		return models.Nessus{}, errors.New("not able to authenticate with server")
	}
	defer res.Body.Close()

	// Parse the response to extract access and secret keys
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return models.Nessus{}, err
	}
	var responseobj models.SessionKeys
	err = json.Unmarshal(bodyBytes, &responseobj)
	if err != nil {
		return models.Nessus{}, err
	}
	context.AccessKey = responseobj.AccessKey
	context.SecretKey = responseobj.SecretKey
	return context, nil
}

// getScan retrieves all scans from Nessus Essentials.
func getScan(_context models.Nessus) ([]models.Scan, error) {
	geturl := _context.Url + "/scans?"

	// Create an HTTP request
	r, err := http.NewRequest("GET", geturl, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", _context.AccessKey, _context.SecretKey))

	// Execute the request
	res, err := _context.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("not able to authenticate with server")
	}
	defer res.Body.Close()

	// Parse the response to extract scan details
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseobj models.ScanResponse
	err = json.Unmarshal(bodyBytes, &responseobj)
	if err != nil {
		return nil, err
	}
	return responseobj.Scans, nil
}
