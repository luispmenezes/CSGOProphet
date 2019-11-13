package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetResults(startDate string, endDate string, stars int, demoRequired bool) []string {
	var resultsURLs []string
	var maxRecords int
	resultsURLs, _, maxRecords = getResultsPage(0, startDate, endDate, stars, demoRequired)
	for i := 100; i < maxRecords; i += 100 {
		pageURLs, _, _ := getResultsPage(0, startDate, endDate, stars, demoRequired)
		resultsURLs = append(resultsURLs, pageURLs...)
	}
	return resultsURLs
}

func getResultsPage(offset int, startDate string, endDate string, stars int, demoRequired bool) ([]string, int, int) {
	log.Printf("Get Result Page (offset:%d,start:%s,end:%s,star%d,demo:%t)", offset, startDate, endDate, stars, demoRequired)

	httpClient := &http.Client{
		Timeout: requestTimeout * time.Second,
	}

	parameters := url.Values{}

	if offset > 0 {
		parameters.Add(resultsParamOffset, strconv.Itoa(offset))
	}

	if len(startDate) > 0 {
		parameters.Add(resultsParamStartDate, startDate)
	}

	if len(endDate) > 0 {
		parameters.Add(resultsParamEndDate, endDate)
	}

	if stars > 0 && stars <= 5 {
		parameters.Add(resultsParamStars, strconv.Itoa(stars))
	}

	if demoRequired {
		parameters.Add(resultsParamContent, ContentDEMO)
	}

	req, err := http.NewRequest("GET", resultsURL, strings.NewReader(parameters.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", userAgent)

	log.Printf("New Request-> URL:%s", req.URL)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error sending http request", err)
	}

	defer resp.Body.Close()

	//TODO: Implement Status Check

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body", err)
	}

	paginationTokens := strings.Split(document.Find(".pagination-data").First().Text(), " ")
	recordCount, _ := strconv.Atoi(paginationTokens[len(paginationTokens)-2])
	lastRecord, _ := strconv.Atoi(paginationTokens[2])

	matchUrlList := []string{}

	document.Find(".result-con").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find("a").Attr("href")
		matchUrlList = append(matchUrlList, url)
	})

	return matchUrlList, lastRecord, recordCount
}

func getResultDemoLink(resultURL string) string {
	log.Printf("Get DemoLink (matchUrl:%s)", resultURL)

	httpClient := &http.Client{
		Timeout: requestTimeout * time.Second,
	}

	req, err := http.NewRequest("GET", baseURL+resultURL, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", userAgent)

	log.Printf("New Request-> URL:%s", req.URL)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//TODO: Implement Status Check

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body", err)
	}

	demoURL, _ := document.Find(".stream-box > a").First().Attr("href")

	return demoURL
}
