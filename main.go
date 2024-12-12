package main

import (
	"NessusEssAutomation/services"
	"flag"
	"log"
)

func main() {
	var template string
	flag.StringVar(&template, "template", "test.html", "one of the templates in reportformats folder")
	sendgridmail := flag.Bool("sm", false, "send mail through sendgrid api")
	loop := flag.Bool("L", false, "Run continously")
	last := flag.Bool("l", false, "If true return last scan else will check if scan was performed in the last 12 hours if not will wait for scan")
	log.Print("Started nessus handling utility")
	flag.Parse()

	if *loop {
		services.Run(*sendgridmail, *last, template)
	} else {
		services.RunOnce(*sendgridmail, *last, template)
	}

}
