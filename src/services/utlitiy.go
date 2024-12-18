package services

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

// getReceivers reads email addresses from the specified file and validates them.
// It returns a slice of valid email addresses or an error if the file cannot be read.
func getReceivers(receivers string) ([]string, error) {
	file, err := os.OpenFile(receivers, os.O_RDONLY, 0777)
	if err != nil {
		log.Print("Could not open recipients file")
		return nil, err
	}
	defer file.Close()

	var emailAddress []string
	filescanner := bufio.NewScanner(file)
	filescanner.Split(bufio.ScanLines)

	// Read each line, validate, and append valid email addresses to the slice
	for filescanner.Scan() {
		mailAddress := filescanner.Text()
		if isEmailValid(mailAddress) {
			emailAddress = append(emailAddress, mailAddress)
		}
	}
	return emailAddress, nil
}

// isEmailValid checks if a given string is a valid email address using a regex.
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// cleanup removes a given file from the filesystem and logs the removal.
func cleanup(report string) {
	os.Remove(report)
	log.Println("Successfully removed report from local filesystem")
}

// RunOnce performs a single execution of the report generation process.
// It authenticates, retrieves scan items, generates reports, and optionally sends emails and cleans up.
func RunOnce(mail bool, last bool, clean bool, template string) {
	tempFile := "temp.html"
	template = "assets/reportformats/" + template

	// Authenticate and get context for Nessus operations
	context, err := authenticate()
	check(err)

	for {
		// Get scan items from Nessus
		records, err := getScanItems(context, last)
		if fmt.Sprint(err) == "no new scan" {
			log.Print("No new scan found, will wait for the next scan")
			time.Sleep(time.Hour * 1)
		} else if err != nil {
			check(err)
		} else {
			log.Print("New scan found in the last 12 hours")

			// Generate HTML report
			err = createHtml(records, template, tempFile)
			check(err)

			// Convert HTML report to PDF
			pdfName := fmt.Sprintf("%s_Vulnerability Assessment Report.pdf", time.Now().Local().Format("2006-01-02"))
			err = createPdf(tempFile, pdfName)
			check(err)

			// Optionally send the PDF report via email
			if mail {
				err = sendgridMail(pdfName)
				check(err)
			}

			// Optionally clean up the generated PDF
			if clean {
				cleanup(pdfName)
			}
			os.Exit(0)
		}
	}
}

// Run performs a continuous execution of the report generation process.
// It generates and processes reports in a loop, waiting 12 hours between iterations.
func Run(mail bool, last bool, clean bool, template string) {
	tempFile := "temp.html"
	template = "assets/reportformats/" + template

	// Authenticate and get context for Nessus operations
	context, err := authenticate()
	check(err)

	for {
		// Get scan items from Nessus
		records, err := getScanItems(context, last)
		if fmt.Sprint(err) == "no new scan" {
			log.Print("No new scan found, will wait for the next scan")
			time.Sleep(time.Hour * 1)
		} else if err != nil {
			check(err)
		} else {
			log.Print("New scan found in the last 12 hours")

			// Generate HTML report
			err = createHtml(records, template, tempFile)
			check(err)

			// Convert HTML report to PDF
			pdfName := fmt.Sprintf("%s_Vulnerability Assessment Report.pdf", time.Now().Local().Format("2006-01-02"))
			err = createPdf(tempFile, pdfName)
			check(err)

			// Optionally send the PDF report via email
			if mail {
				err = sendgridMail(pdfName)
				check(err)
			}

			// Optionally clean up the generated PDF
			if clean {
				cleanup(pdfName)
			}

			// Wait 12 hours before the next iteration
			time.Sleep(time.Hour * 12)
		}
	}
}

// check logs a fatal error and exits if the error is not nil.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Initlogger initializes logging to a specific file.
// It ensures all logs are written to "/var/log/nessus_ess".
func Initlogger(data string) {
	f, err := os.OpenFile("/var/log/nessus_ess", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()

	// Set the log output to the specified file
	log.SetOutput(f)
	log.Println(data)
}
