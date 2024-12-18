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

// createHtml generates an HTML report from a template and a list of vulnerabilities.
// It populates the template with vulnerability data and saves it to the specified output path.
func createHtml(records [][]string, templateFilePath string, outputHTMLPath string) error {
	// Parse the HTML template file
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		fmt.Println("Error parsing template file:", err)
		return err
	}

	// Process the records to create a slice of vulnerabilities
	var vulnerabilities []models.Vulnerability
	for _, record := range records[1:] { // Skip the header row
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
			VulnerabilityFirstDiscoveredDate: "N/A", // Default value
			VulnerabilityLastObservedDate:    "N/A", // Default value
			PluginPublicationDate:            record[15],
		}

		// Determine exploitability based on certain fields
		if record[17] == "" && record[18] == "" && record[19] == "" {
			vulnerability.Exploitability = "No"
		} else {
			vulnerability.Exploitability = "Yes"
		}

		// Update or insert the vulnerability into the database
		vulnerability, err = dbaccess(vulnerability)
		if err != nil {
			return err
		}

		// Append the processed vulnerability to the slice
		vulnerabilities = append(vulnerabilities, vulnerability)
	}

	// Prepare the data for the report
	reportData := models.ReportData{
		DateGenerated:   time.Now().Format("2006-01-02"),
		Vulnerabilities: vulnerabilities,
	}

	// Create the output HTML file
	file, err := os.Create(outputHTMLPath)
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return err
	}
	defer file.Close()

	// Execute the template with the report data and write to the file
	err = tmpl.Execute(file, reportData)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	log.Println("Report generated successfully: " + outputHTMLPath)
	return nil
}

// createPdf generates a PDF report from an HTML file.
// It uses wkhtmltopdf to convert the HTML file to a PDF and saves it to the specified output path.
func createPdf(inputHtmlPath string, outputPDFPath string) error {
	// Initialize the PDF generator
	pdf, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Print(err.Error())
		log.Println("Error initializing PDF generator")
		return err
	}

	// Add the input HTML file as a page in the PDF
	pdf.AddPage(wkhtmltopdf.NewPage(inputHtmlPath))

	// Create the PDF
	err = pdf.Create()
	if err != nil {
		log.Println("Error creating PDF")
		return err
	}

	// Write the PDF to the specified output file
	err = pdf.WriteFile(outputPDFPath)
	if err != nil {
		log.Println("Error writing PDF to file:", err)
		return err
	}

	// Remove the temporary HTML file after generating the PDF
	err = os.Remove(inputHtmlPath)
	if err != nil {
		log.Println("Error removing HTML file:", err)
		return err
	}

	log.Println("PDF generated successfully: " + outputPDFPath)
	return nil
}
