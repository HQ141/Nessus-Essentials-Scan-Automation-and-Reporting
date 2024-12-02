package services

import (
	"NessusEssAutomation/models"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/buger/jsonparser"
)

func GetScanItems(context models.Nessus) ([][]string, error) {
	Scans, err := GetScan(context)
	if err != nil {
		return nil, err
	}
	myscan := models.Scan{}
	for _, i := range Scans {
		if i.Name == context.ScanName {
			myscan = i
		}
	}
	if myscan == (models.Scan{}) {
		return nil, errors.New("no matching scan name found")
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
	if res.StatusCode != 200 {
		return nil, errors.New("http request error")
	}
	body, err := ioutil.ReadAll(res.Body)
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
		report, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("exported")
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

func CreateHtml(records [][]string, outputHTMLPath string) error {
	templateFilePath := "assets/reportformats/Vidizmoreport.html"
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		fmt.Println("Error parsing template file:", err)
		return err
	}
	var vulnerabilities []models.Vulnerability
	for _, record := range records[1:] {
		vulnerability := models.Vulnerability{
			PluginID:                         record[0],
			AssetName:                        record[4],
			PluginName:                       record[7],
			Severity:                         record[3],
			CVSS3BaseScore:                   record[13],
			VulnerabilityPriorityRating:      "N/A",
			PluginOutput:                     record[12],
			Synopsis:                         record[8],
			Description:                      record[9],
			Solution:                         record[10],
			VulnerabilityFirstDiscoveredDate: "N/A",
			VulnerabilityLastObservedDate:    "N/A",
			PluginPublicationDate:            record[15],
		}

		if strings.Contains(record[16], "Metasploit") || strings.Contains(record[16], "Core Impact") || strings.Contains(record[16], "CANVAS") {
			vulnerability.Exploitability = "Yes"
		} else {
			vulnerability.Exploitability = "No"
		}
		vulnerability, err = dbaccess(vulnerability)
		vulnerabilities = append(vulnerabilities, vulnerability)
		if err != nil {
			return err
		}
	}

	reportData := models.ReportData{
		DateGenerated:   time.Now().Format("2006-01-02"),
		Vulnerabilities: vulnerabilities,
	}

	file, err := os.Create(outputHTMLPath)
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, reportData)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	fmt.Println("Report generated successfully: " + outputHTMLPath)
	return nil
}

func CreatePdf(inputHtmlPath string, outputPDFPath string) error {
	pdf, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal("error initializing pdf generator")
		return err
	}
	pdf.AddPage(wkhtmltopdf.NewPage(inputHtmlPath))
	err = pdf.Create()
	if err != nil {
		log.Fatal("error createing pdf")
		return err
	}
	err = pdf.WriteFile(outputPDFPath)
	if err != nil {
		fmt.Println("Error writing PDF to file:", err)
		return err
	}
	err = os.Remove(inputHtmlPath)
	if err != nil {
		fmt.Println("Error removing HTML to file:", err)
		return err
	}
	fmt.Println("PDF generated successfully: " + outputPDFPath)
	return nil
}
