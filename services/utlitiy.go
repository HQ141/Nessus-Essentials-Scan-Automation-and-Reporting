package services

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

func getReceivers(recievers string) ([]string, error) {
	file, err := os.OpenFile(recievers, os.O_RDONLY, 0777)
	if err != nil {
		log.Print("Could not open receipients file")
		return nil, err
	}
	defer file.Close()
	var emailAddress []string
	filescanner := bufio.NewScanner(file)
	filescanner.Split(bufio.ScanLines)
	for filescanner.Scan() {
		mailAddress := filescanner.Text()
		if isEmailValid(mailAddress) {
			emailAddress = append(emailAddress, mailAddress)
		}
	}
	return emailAddress, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func cleanup(report string) {
	os.Remove(report)
	log.Println("Successfully removed report from local filesystem")
}

func RunOnce(mail bool, last bool, template string) {
	temp_file := "temp.html"
	template = "assets/reportformats/" + template
	context, err := authenticate()
	check(err)
	// Get scan items
	records, err := getScanItems(context, last)
	if fmt.Sprint(err) == "no new scan" {
		log.Print("no new scan found will wait for scan")
		time.Sleep(time.Hour * 1)
	} else if err != nil {
		check(err)
	} else {
		log.Print("new scan in last 12 hours")
		// Create HTML report
		err = createHtml(records, template, temp_file)
		check(err)

		// Create PDF from HTML
		pdfname := fmt.Sprintf("%s_Vulnerability Assessment Report.pdf", time.Now().Local().Format("2006-01-02"))
		err = createPdf(temp_file, pdfname)
		check(err)

		if mail {
			//Send New PDF to recipients
			err = sendgridMail(pdfname)
			check(err)
		}
		cleanup(pdfname)
		os.Exit(0)
	}
}
func Run(mail bool, last bool, template string) {
	temp_file := "temp.html"
	template = "assets/reportformats/" + template
	//Authenticate and get the context
	context, err := authenticate()
	check(err)
	for {
		// Get scan items
		records, err := getScanItems(context, last)
		if fmt.Sprint(err) == "no new scan" {
			log.Print("no new scan found will wait for scan")
			time.Sleep(time.Hour * 1)
		} else if err != nil {
			check(err)
		} else {
			log.Print("new scan in last 12 hours")
			// Create HTML report
			err = createHtml(records, template, temp_file)
			check(err)

			// Create PDF from HTML
			pdfname := fmt.Sprintf("%s_Vulnerability Assessment Report.pdf", time.Now().Local().Format("2006-01-02"))
			err = createPdf(temp_file, pdfname)
			check(err)

			//Send New PDF to recipients
			// err = services.SendgridMail(pdfname)
			// check(err)
			cleanup(pdfname)
			time.Sleep(time.Hour * 12)
		}
	}
}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
