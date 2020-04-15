package web_crawler

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{
	Timeout: requestTimeout * time.Second,
}

var baseURL = &url.URL{Scheme: "https", Host: baseHost}

func NewRequest(method, path string, queryParameters map[string]string) (*goquery.Document, error) {
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
		return nil, errors.New(fmt.Sprintf("Error executing http request %d", response.StatusCode))
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Error("Error loading HTTP response body", err)
	}

	return document, err
}
