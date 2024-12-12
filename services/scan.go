package services

import (
	"NessusEssAutomation/models"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func createHtml(records [][]string, templateFilePath string, outputHTMLPath string) error {
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
			VulnerabilityPriorityRating:      record[3],
			PluginOutput:                     record[12],
			Synopsis:                         record[8],
			Description:                      record[9],
			Solution:                         record[10],
			VulnerabilityFirstDiscoveredDate: "N/A",
			VulnerabilityLastObservedDate:    "N/A",
			PluginPublicationDate:            record[15],
		}
		if record[17] == "" || record[18] == "" || record[19] == "" {
			vulnerability.Exploitability = "No"
		} else {
			vulnerability.Exploitability = "Yes"
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

	log.Println("Report generated successfully: " + outputHTMLPath)
	return nil
}

func createPdf(inputHtmlPath string, outputPDFPath string) error {
	pdf, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Print(err.Error())
		log.Println("error initializing pdf generator")
		return err
	}
	pdf.AddPage(wkhtmltopdf.NewPage(inputHtmlPath))
	err = pdf.Create()
	if err != nil {
		log.Println("error createing pdf")
		return err
	}
	err = pdf.WriteFile(outputPDFPath)
	if err != nil {
		log.Println("Error writing PDF to file:", err)
		return err
	}
	err = os.Remove(inputHtmlPath)
	if err != nil {
		log.Println("Error removing HTML to file:", err)
		return err
	}
	log.Println("PDF generated successfully: " + outputPDFPath)
	return nil
}
