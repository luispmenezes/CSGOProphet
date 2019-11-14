package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"./model"

	"github.com/PuerkitoBio/goquery"
)

func GetDemoLinks(startDate string, endDate string, stars int, demoRequired bool) []model.DemoLink {
	var demoLinkList []model.DemoLink
	resultURLList := getResults(startDate, endDate, stars, demoRequired)

	for _, resultURL := range resultURLList {
		demoLinkList = append(demoLinkList, getResultDemoLink(resultURL))
	}

	return demoLinkList
}

func getResults(startDate string, endDate string, stars int, demoRequired bool) []string {
	var resultsURLs []string
	var maxRecords int
	resultsURLs, _, maxRecords = getResultsPage(0, startDate, endDate, stars, demoRequired)
	log.Printf("Result Records Count: %d", maxRecords)
	for i := 100; i < maxRecords; i += 100 {
		pageURLs, _, _ := getResultsPage(i, startDate, endDate, stars, demoRequired)
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

	requestURL := resultsURL
	if len(parameters) > 0 {
		requestURL += "?" + parameters.Encode()
	}

	log.Printf("New Request-> URL:%s", requestURL)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", userAgent)

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

	matchURLList := []string{}

	document.Find(".results-holder > .results-all > .results-sublist > .result-con > .a-reset").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		matchURLList = append(matchURLList, url)
	})

	return matchURLList, lastRecord, recordCount
}

func getResultDemoLink(resultURL string) model.DemoLink {
	log.Printf("Get DemoLink (matchUrl:%s)", resultURL)

	httpClient := &http.Client{
		Timeout: requestTimeout * time.Second,
	}

	requestURL := baseURL + resultURL

	log.Printf("New Request-> URL:%s", requestURL)

	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", userAgent)

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
	timeStamp, _ := document.Find(".teamsBox > .timeAndEvent > .time").First().Attr("data-unix")

	return model.DemoLink{DemoURL: demoURL, MatchResultURL: resultURL, Timestamp: timeStamp}
}
