package main

import (
	"NessusEssAutomation/services"
	"flag"
	"log"
)

// main is the entry point for the Nessus Essentials Automation utility.
// It parses command-line flags to determine the desired operation mode and parameters.
func main() {
	// Define a flag for specifying the HTML template file to use for reports.
	var template string
	flag.StringVar(&template, "template", "test.html", "Specify a template from the reportformats folder")

	// Define a flag to enable or disable sending the report via SendGrid.
	sendgridmail := flag.Bool("m", false, "Enable to send mail through SendGrid API")

	// Define a flag to enable continuous operation (loop mode).
	loop := flag.Bool("L", false, "Enable to run continuously in a loop")

	// Define a flag to enable cleaning up reports after processing.
	cleanup := flag.Bool("c", false, "Enable to clean reports after process completion")

	// Define a flag to determine whether to return the last scan or wait for a new scan.
	last := flag.Bool("l", false, "If true, return the last scan; otherwise, wait for a scan if none performed in the last 12 hours")

	// Log the start of the utility.
	log.Print("Started Nessus Essentials Automation Utility")

	// Parse the command-line flags.
	flag.Parse()

	// Check if loop mode is enabled, and run the appropriate function.
	if *loop {
		// Run in continuous loop mode
		services.Run(*sendgridmail, *last, *cleanup, template)
	} else {
		// Run once and exit
		services.RunOnce(*sendgridmail, *last, *cleanup, template)
	}
}
