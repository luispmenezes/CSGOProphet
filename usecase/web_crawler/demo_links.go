package web_crawler

import (
	"csgo_prophet/model/web_crawler"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetDemoLinks(startDate string, endDate string, stars int, demoRequired bool) ([]web_crawler.DemoLink, error) {
	var demoLinkList []web_crawler.DemoLink
	resultURLList, err := getResults(startDate, endDate, stars, demoRequired)

	if err != nil {
		return nil, err
	}

	for _, resultURL := range resultURLList {
		demoLinkList = append(demoLinkList, getMatchDemoLink(resultURL))
	}

	return demoLinkList, nil
}

func getMatchDemoLink(matchUrl string) web_crawler.DemoLink {
	log.WithFields(log.Fields{"matchUrl": matchUrl}).Info("Get DemoLink")

	document, err := SendRequestWithRetry(http.MethodGet, matchUrl, nil)

	if err != nil {
		log.Error("Error loading HTTP response body", err)
	}

	demoURL, _ := document.Find(".stream-box > a").First().Attr("href")
	timeStamp, _ := document.Find(".teamsBox > .timeAndEvent > .time").First().Attr("data-unix")

	return web_crawler.DemoLink{DemoURL: demoURL, MatchResultURL: matchUrl, Timestamp: timeStamp}
}
