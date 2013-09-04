package main

import (
	"fmt"
)

const (
	MaxRedirects = 10
)

func main() {
	redirects := ReadCsv("301s.csv")

	log := make([]redirectResult, 0)

	for _, info := range redirects {
		result := CheckUrl(info)
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
