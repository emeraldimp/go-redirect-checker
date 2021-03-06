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

func ReadCsv(name string) ([]redirectInfo, error) {
	csvFile, err := os.Open(name)

	if err != nil {
		return nil, err
	}

	defer csvFile.Close()

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

	return result, nil
}

func CheckUrl(info redirectInfo, maxRedirects int) redirectResult {

	result := redirectResult{
		Url:              info.Url,
		ExpectedUrl:      info.ExpectedUrl,
		Redirects:        0,
		IntermediateUrls: make([]actualRedirect, 0),
		Error:            nil,
	}

	nextUrl := info.Url
	for {
		if result.Redirects >= maxRedirects {
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
