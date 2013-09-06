/***********
*
* redirectChecker - Checks 301 redirects from a CSV file
*
* Written in 2013 by Geoff Lehr jademystery@hotmail.com
*
* To the extent possible under law, the author(s) have dedicated all
* copyright and related and neighboring rights to this software to the 
* public domain worldwide. This software is distributed without any warranty. 
*
* You should have received a copy of the CC0 Public Domain Dedication along
* with this software. 
* If not, see <http://creativecommons.org/publicdomain/zero/1.0/>. 
*
************/

package main

import (
	"fmt"
	"flag"
	"os"
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

	redirects, err := ReadCsv(filename)

	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

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
