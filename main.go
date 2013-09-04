package main

import (
	"fmt"
	"flag"
)

const (
	MaxRedirects = 10
)

func main() {
	maxRedirects := flag.Int("max-redirects", MaxRedirects, "The maximum number of redirects to follow before giving up (excluding detected loops)")

	flag.Parse()

	filename := flag.Arg(0)

	if filename == "" {
		filename = "301s.csv"
	}

	redirects := ReadCsv(filename)

	log := make([]redirectResult, 0)

	for _, info := range redirects {
		result := CheckUrl(info, *maxRedirects)
		log = append(log, result)
	}

	for _, logItem := range log {

		if logItem.FinalUrl == logItem.ExpectedUrl {
			fmt.Printf("OK: %v Matched\n", logItem.Url)
			continue
		}

		if logItem.LooksLikeRedirectLoop() {
			fmt.Printf("LOOP: %v Redirect Loop? Stopped after %v redirects\n", logItem.Url, logItem.Redirects)
			continue
		}

		fmt.Printf("ERR: %v Unexpected destination: %v\n", logItem.Url, logItem.FinalUrl)
	}
}
