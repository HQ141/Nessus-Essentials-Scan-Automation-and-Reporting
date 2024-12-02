package main

import (
	"NessusEssAutomation/services"
	"log"
)

func main() {
	// Authenticate and get the context
	context, err := services.Authenticate()
	if err != nil {
		log.Fatal(err)
		return
	}
	// Get scan items
	records, err := services.GetScanItems(context)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Create HTML report
	err = services.CreateHtml(records, "hello.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create PDF from HTML
	err = services.CreatePdf("hello.html", "hello.pdf")
	if err != nil {
		log.Fatal(err)
		return
	}

}
