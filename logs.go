package main

import (
	"flag"
	"fmt"
	"github.com/aodin/logs/nginx"
	"log"
	"time"
)

func main() {
	// Allow use to specify the directory that should be parsed
	// Default to the current directory
	var path string
	flag.StringVar(&path, "i", ".", "directory to parse")

	// The user can enable a daily count of views
	var byDate bool
	flag.BoolVar(&byDate, "date", false, "count the views by date")

	// Allow the user to set the timezone, default to UTC
	// Example timezone: America/Los_Angeles
	var tz string
	flag.StringVar(&tz, "tz", "", "timezone of views")
	flag.Parse()

	// TODO Count multiple urls
	// The user must specify a URL to find
	url := flag.Arg(0)
	if url == "" {
		log.Fatal("Please specify a URL to find")
	}

	// If the user specified a timezone, load it
	var loc *time.Location
	var tzErr error
	if tz != "" {
		loc, tzErr = time.LoadLocation(tz)
		if tzErr != nil {
			log.Fatalf("Error while loading timezone: %s", tzErr)
		}
	} else {
		loc = time.UTC
	}

	// Read the specified directory and generate logging information
	logs, err := nginx.ReadDirectory(path)
	if err != nil {
		log.Fatalf("Error while reading directory %s: %s", path, err)
	}
	log.Printf("%d files examined\n", logs.Files)
	log.Printf("%d access logs\n", len(logs.Plain))
	log.Printf("%d gzipped access logs\n", len(logs.Zipped))

	// Build the records
	// TODO Stream the records
	records := make([]nginx.Record, 0)
	for _, f := range logs.Plain {
		rs, err := nginx.OpenAndParseFile(f)
		if err != nil {
			log.Fatalf("Error while parsing file %s: %s", f, err)
		}
		records = append(records, rs...)
	}
	for _, f := range logs.Zipped {
		rs, err := nginx.OpenAndParseZipFile(f)
		if err != nil {
			log.Fatalf("Error while parsing zipped file %s: %s", f, err)
		}
		records = append(records, rs...)
	}

	log.Printf("Total records parsed: %d\n", len(records))
	log.Println(byDate)
	if byDate {
		total, counts := nginx.CountDays(records, url, loc)
		fmt.Printf("\n%d Total Views\n\n", total)
		for _, count := range counts {
			fmt.Printf("%s: %d\n", count.Date, count.Count)
		}
	} else {
		total := nginx.Count(records, url)
		fmt.Printf("\n%d Total Views\n\n", total)
	}
}
