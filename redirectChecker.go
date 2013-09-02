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

func (rr *redirectResult) AppendIntermediate(nextUrl string, redirected string, code int) {
	intermediateUrl := actualRedirect{
		OriginUrl:     nextUrl,
		RedirectedUrl: redirected,
		ErrorCode:     code,
	}

	rr.IntermediateUrls = append(rr.IntermediateUrls, intermediateUrl)
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

	result := redirectResult{
		Url:              info.Url,
		ExpectedUrl:      info.ExpectedUrl,
		Redirects:        0,
		IntermediateUrls: make([]actualRedirect, 0),
		Error:            nil,
	}

	nextUrl := info.Url
	for {
		if result.Redirects > 5 {
			break
		}

		req, err := http.NewRequest("GET", nextUrl, nil)
		resp, err := http.DefaultTransport.RoundTrip(req)

		if err != nil {
			result.Error = err
			break
		}

		redirectTo := resp.Header.Get("Location")
		result.AppendIntermediate(nextUrl, redirectTo, resp.StatusCode)
		result.FinalUrl = nextUrl

		if resp.StatusCode != 301 {
			break
		}

		result.Redirects++

		nextUrl = redirectTo
	}

	return result
}

func main() {
	redirects := readCsv("301s.csv")

	log := make([]redirectResult, 0)

	for _, info := range redirects {
		result := checkUrl(info)
		log = append(log, result)
	}

	for _, logItem := range log {
		fmt.Printf("Original Url: %v\n", logItem.Url)
		fmt.Printf("Final Url: %v\n", logItem.FinalUrl)
		fmt.Printf("Expected Url: %v\n", logItem.ExpectedUrl)
		fmt.Printf("Number of Redirects: %v\n", logItem.Redirects)
		fmt.Printf("\n")
	}

	//fmt.Println(log)
}
