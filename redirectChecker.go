package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
)

type redirectInfo struct {
	Url         string
	ExpectedUrl string
}

type actualRedirect struct {
	OriginUrl     string
	RedirectedUrl string
	ErrorCode     int
}

type redirectResult struct {
	Url              string
	ExpectedUrl      string
	Redirects        int
	IntermediateUrls []actualRedirect
	FinalUrl         string
	Error            error
}

func readCsv(name string) []redirectInfo {
	csvFile, err := os.Open(name)
	defer csvFile.Close()

	if err != nil {
		panic(err)
	}

	csvReader := csv.NewReader(csvFile)
	result := make([]redirectInfo, 0)
	for {
		fields, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		info := redirectInfo{fields[0], fields[1]}

		result = append(result, info)
	}

	return result
}

func checkUrl(info redirectInfo) redirectResult {
	currentUrl := info.Url
	expected := info.ExpectedUrl

	redirects := 0
	nextUrl := currentUrl
	for {
		if redirects > 5 {
			break
		}

		req, err := http.NewRequest("GET", nextUrl, nil)
		resp, err := http.DefaultTransport.RoundTrip(req)

		if err != nil {
			break
		}

		if resp.StatusCode != 301 {
			break
		}

		redirects++
	}

	result := redirectResult{
		Url:         currentUrl,
		ExpectedUrl: expected,
		Redirects:   redirects,
		//IntermediateUrls:
		//FinalUrl: ,
		//Error: err,
	}

	return result
}

func main() {
	redirects := readCsv("301s.csv")

	fmt.Println(redirects)

	log := make([]redirectResult, 0)

	for _, info := range redirects {
		result := checkUrl(info)
		log = append(log, result)
	}

	fmt.Println(log)
}
