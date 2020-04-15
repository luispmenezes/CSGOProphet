package web_crawler

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func getResults(startDate string, endDate string, stars int, demoRequired bool) ([]string, error) {
	var resultsURLs []string
	var maxRecords int
	var err error

	resultsURLs, _, maxRecords, err = getResultsPage(0, startDate, endDate, stars, demoRequired)

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{"recordCount": maxRecords}).Info("Result Records Count")

	for i := 100; i < maxRecords; i += 100 {
		pageURLs, _, _, pageErr := getResultsPage(i, startDate, endDate, stars, demoRequired)
		if pageErr != nil {
			return nil, err
		}
		resultsURLs = append(resultsURLs, pageURLs...)
	}
	return resultsURLs, nil
}

func getResultsPage(offset int, startDate string, endDate string, stars int, demoRequired bool) ([]string, int, int, error) {
	log.WithFields(log.Fields{
		"offset":       offset,
		"startDate":    startDate,
		"endDate":      endDate,
		"stars":        stars,
		"demoRequired": demoRequired,
	}).Info("Get Result Page")

	parameters := make(map[string]string)

	if offset > 0 {
		parameters[resultsParamOffset] = strconv.Itoa(offset)
	}

	if len(startDate) > 0 {
		parameters[resultsParamStartDate] = startDate
	}

	if len(endDate) > 0 {
		parameters[resultsParamEndDate] = endDate
	}

	if stars > 0 && stars <= 5 {
		parameters[resultsParamStars] = strconv.Itoa(stars)
	}

	if demoRequired {
		parameters[resultsParamContent] = ContentDEMO
	}

	document, err := NewRequest(http.MethodGet, resultsPath, parameters)
	if err != nil {
		log.Error("Error loading HTTP response body", err)
		return nil, 0, 0, err
	}

	paginationTokens := strings.Split(document.Find(".pagination-data").First().Text(), " ")
	recordCount, _ := strconv.Atoi(paginationTokens[len(paginationTokens)-2])
	lastRecord, _ := strconv.Atoi(paginationTokens[2])

	var matchURLList []string

	document.Find(".results-holder > .results-all > .results-sublist > .result-con > .a-reset").Each(func(i int, s *goquery.Selection) {
		matchUrl, _ := s.Attr("href")
		matchURLList = append(matchURLList, matchUrl)
	})

	return matchURLList, lastRecord, recordCount, nil
}
