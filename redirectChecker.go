package main

import (
	"encoding/csv"
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

func (rr *redirectResult) LooksLikeRedirectLoop() bool {

	urls := make(map[string]int)
	for _, url := range rr.IntermediateUrls {
		if urls[url.OriginUrl] != 0 {
			return true
		}

		urls[url.OriginUrl]++
	}

	return false
}

func ReadCsv(name string) []redirectInfo {
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

func CheckUrl(info redirectInfo) redirectResult {

	result := redirectResult{
		Url:              info.Url,
		ExpectedUrl:      info.ExpectedUrl,
		Redirects:        0,
		IntermediateUrls: make([]actualRedirect, 0),
		Error:            nil,
	}

	nextUrl := info.Url
	for {
		if result.Redirects >= MaxRedirects {
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

		if result.LooksLikeRedirectLoop() {
			break
		}

		result.Redirects++

		nextUrl = redirectTo
	}

	return result
}
