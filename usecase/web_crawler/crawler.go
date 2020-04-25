package web_crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"math"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{
	Timeout: requestTimeout * time.Second,
}

var baseURL = &url.URL{Scheme: "https", Host: baseHost}

func newRequest(method, path string, queryParameters map[string]string) (*goquery.Document, error) {
	log.WithFields(log.Fields{
		"method":      method,
		"path":        path,
		"queryParams": queryParameters,
	}).Debug("New Crawler Http request")

	requestURL := baseURL
	requestURL.Path = path
	requestQuery := requestURL.Query()

	for paramKey, paramValue := range queryParameters {
		requestQuery.Set(paramKey, paramValue)
	}

	requestURL.RawQuery = requestQuery.Encode()

	request, err := http.NewRequest(method, requestURL.String(), http.NoBody)

	if err != nil {
		return nil, err
	}

	request.Header.Add("User-Agent", userAgent)

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusTooManyRequests {
			return nil, errors.New("too many requests")
		} else {
			return nil, errors.New(fmt.Sprintf("error executing http request status: %d", response.StatusCode))
		}
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Error loading HTTP response body", err)
	}

	return document, err
}

func SendRequestWithRetry(method, path string, queryParameters map[string]string) (*goquery.Document, error) {
	var result *goquery.Document
	var err error
	idx := 0.0

	for {
		result, err = newRequest(method, path, queryParameters)
		if err != nil && err.Error() == "too many requests" {
			log.Debug("Too many requests backing off -> idx :", idx)
			time.Sleep(time.Duration(math.Min(math.Pow(2.0, idx), 128.0)) * time.Second)
			idx++
		} else {
			time.Sleep(1 * time.Second)
			break
		}
	}

	return result, err
}
