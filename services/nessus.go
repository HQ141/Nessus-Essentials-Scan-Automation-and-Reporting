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
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

func getScanItems(context models.Nessus, last bool) ([][]string, error) {
	myscan := models.Scan{}
	Scans, err := getScan(context)
	if err != nil {
		return nil, err
	}
	for _, i := range Scans {
		if i.Name == context.ScanName {
			myscan = i
		}
	}
	if myscan == (models.Scan{}) {
		return nil, errors.New("no matching scan name found")
	}
	if !last {
		if time.Since(time.Unix(myscan.LastModificationDate, 0)) >= 12*time.Hour {
			return [][]string{}, errors.New("no new scan")
		}
	}
	values := map[string]interface{}{
		"format":      "csv",
		"template_id": "",
		"reportContents": map[string]interface{}{
			"csvColumns": map[string]bool{
				"id":                   true,
				"cve":                  true,
				"cvss":                 true,
				"risk":                 true,
				"hostname":             true,
				"protocol":             true,
				"port":                 true,
				"plugin_name":          true,
				"synopsis":             true,
				"description":          true,
				"solution":             true,
				"see_also":             true,
				"plugin_output":        true,
				"stig_severity":        false,
				"cvss4_base_score":     false,
				"cvss4_bt_score":       false,
				"cvss3_base_score":     true,
				"cvss_temporal_score":  false,
				"cvss3_temporal_score": false,
				"vpr_score":            false,
				"epss_score":           false,
				"risk_factor":          true,
				"references":           false,
				"plugin_information":   true,
				"exploitable_with":     true,
			},
		},
	}
	jsonValues, _ := json.Marshal(values)
	r, err := http.NewRequest("POST", context.Url+fmt.Sprintf("/scans/%d/export", myscan.ID), bytes.NewBuffer(jsonValues))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", context.AccessKey, context.SecretKey))
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 409 {
		return [][]string{}, errors.New("no new scan")
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("http request error, recieved %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	token, err := jsonparser.GetString(body, "token")
	if err != nil {
		return nil, err
	}
	for {
		time.Sleep(time.Second * 5)

		req, err := http.NewRequest("GET", context.Url+fmt.Sprintf("/tokens/%s/download", token), nil)
		req.Header.Set("Accept", "text/csv")

		if err != nil {
			return nil, err
		}
		req.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", context.AccessKey, context.SecretKey))
		req.Header.Set("Content-Type", "application/json")

		log.Println("trying to export")
		resp, err := context.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusNotFound {
			continue
		}
		report, _ := io.ReadAll(resp.Body)
		log.Println("exported")
		rec := csv.NewReader(strings.NewReader(string(report)))
		rec.ReuseRecord = true
		records, _ := rec.ReadAll()
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

// Authenticate performs authentication and returns Nessus client
func authenticate() (models.Nessus, error) {
	context, err := config.LoadNessusConfig()
	if err != nil {
		return models.Nessus{}, err
	}
	customTransport := (http.DefaultTransport.(*http.Transport))
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	context.HttpClient = &http.Client{Transport: customTransport}
	context, err = getToken(context)
	if err != nil {
		return models.Nessus{}, err
	}
	keys, err := getSessionKeys(context)
	if err != nil {
		return models.Nessus{}, err
	}
	return keys, nil
}

// GetToken retrieves the authentication token
func getToken(context models.Nessus) (models.Nessus, error) {
	body := []byte(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, context.Username, context.Password))
	posturl := context.Url + "/session"
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		return models.Nessus{}, err
	}
	r.Header.Add("Content-Type", "application/json")
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return models.Nessus{}, err
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseobj models.Session
	err = json.Unmarshal(bodyBytes, &responseobj)
	if err != nil {
		return models.Nessus{}, err
	}
	context.Token = responseobj.Token
	return context, nil
}

// GetSessionKeys retrieves the session keys for Nessus
func getSessionKeys(context models.Nessus) (models.Nessus, error) {
	puturl := context.Url + "/session/keys"
	r, err := http.NewRequest("PUT", puturl, nil)
	if err != nil {
		return models.Nessus{}, err
	}
	token := fmt.Sprintf("token=%s", context.Token)
	r.Header.Add("X-Cookie", token)
	res, err := context.HttpClient.Do(r)
	if err != nil {
		return models.Nessus{}, err
	}
	if res.StatusCode != 200 {
		return models.Nessus{}, errors.New("not able to authenticate with server")
	}
	defer res.Body.Close()
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
func getScan(_context models.Nessus) ([]models.Scan, error) {
	geturl := _context.Url + "/scans?"
	r, err := http.NewRequest("GET", geturl, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", _context.AccessKey, _context.SecretKey))

	res, err := _context.HttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {

		return nil, errors.New("not able to authenticate with server")
	}
	defer res.Body.Close()
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
