package web_crawler

import (
	"csgo_prophet/model/web_crawler"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetMatchData(startDate string, endDate string, stars int, demoRequired bool) ([]web_crawler.Match, error) {
	var matchList []web_crawler.Match
	resultURLList, err := getResults(startDate, endDate, stars, demoRequired)

	if err != nil {
		return nil, errors.Wrap(err, "match results query request failed")
	}

	for _, resultURL := range resultURLList {
		match, err := parseCompletedMatch(resultURL)
		if err != nil {
			return nil, errors.Wrap(err, "error during match parsing")
		}

		matchList = append(matchList, match)
	}

	return matchList, nil
}

func parseCompletedMatch(matchUrl string) (web_crawler.Match, error) {
	log.WithFields(log.Fields{"matchUrl": matchUrl}).Info("Parsing match")

	document, err := SendRequestWithRetry(http.MethodGet, matchUrl, nil)

	if err != nil {
		return web_crawler.Match{}, errors.Wrap(err, "match overview web request failed")
	}

	match := web_crawler.Match{Url: matchUrl}
	//TIMESTAMP
	timestamp, timestampExists := document.Find(".teamsBox > .timeAndEvent > .time").First().Attr("data-unix")

	if !timestampExists {
		return web_crawler.Match{}, errors.New("match timestamp not found")
	}

	tsMilliseconds, err := strconv.ParseInt(timestamp, 10, 64)

	if err != nil {
		return web_crawler.Match{}, errors.Wrap(err, "match timestamp parsing failed")
	} else {
		match.StartTime = time.Unix(int64(float64(tsMilliseconds)/1000), 0)
	}
	//EVENT
	event, eventExists := document.Find(".teamsBox > .timeAndEvent > .event > a").First().Attr("href")

	if !eventExists {
		return web_crawler.Match{}, errors.New("match event not found")
	}

	match.Event = event

	//FORMAT
	formatStr := document.Find(".veto-box > div").First().Text()

	format, formatErr := strconv.Atoi(string(strings.Split(formatStr, " ")[2][0]))

	if formatErr != nil {
		log.Error("Invalid format ", format)
		return web_crawler.Match{}, formatErr
	}

	match.Format = format

	//DEMO URL
	demoURL, demoUrlExists := document.Find(".stream-box > a").First().Attr("href")

	if !demoUrlExists {
		return web_crawler.Match{}, errors.New("match demo link not found")
	}

	match.DemoUrl = demoURL

	//TEAM 1
	team1, team1Exists := document.Find(".teamsBox > .team > .team1-gradient > a").First().Attr("href")

	if !team1Exists {
		return web_crawler.Match{}, errors.New("match team 1 link not found")
	}

	match.Team1 = team1

	team1Score := document.Find(".teamsBox > .team > .team1-gradient > div").First().Text()

	team1ScoreInt, err := strconv.Atoi(team1Score)

	if err != nil {
		return web_crawler.Match{}, errors.Wrap(err, "team 1 score parsing failed")
	}

	match.Team1Score = team1ScoreInt

	//TEAM 2
	team2, team2Exists := document.Find(".teamsBox > .team > .team2-gradient > a").First().Attr("href")

	if !team2Exists {
		return web_crawler.Match{}, errors.New("match team 2 link not found")
	}

	match.Team2 = team2

	team2Score := document.Find(".teamsBox > .team > .team2-gradient > div").First().Text()

	team2ScoreInt, err := strconv.Atoi(team2Score)

	if err != nil {
		return web_crawler.Match{}, errors.Wrap(err, "team 2 score parsing failed")
	}

	match.Team2Score = team2ScoreInt

	//Players
	var playerList1 []string
	var playerList2 []string

	document.Find(".players > a").Each(func(i int, s *goquery.Selection) {
		playerLink, _ := s.Attr("href")

		if i%5 == 0 {
			playerList1 = append(playerList1, playerLink)
		} else {
			playerList2 = append(playerList2, playerLink)
		}
	})

	if len(playerList1) == 0 || len(playerList2) == 0 {
		return web_crawler.Match{}, errors.New("match player lineup parsing failed")
	}

	match.Team1Composition = playerList1
	match.Team2Composition = playerList2

	//MAP STATS
	var mapLinkList []string
	var mapStatsList []web_crawler.MapStats
	document.Find(".results-center-stats > .results-stats").Each(func(i int, s *goquery.Selection) {
		mapStatsLink, _ := s.Attr("href")
		mapLinkList = append(mapLinkList, mapStatsLink)
	})

	for i, mapLink := range mapLinkList {
		mapStats, err := parseMapStats(mapLink)

		if err != nil {
			return web_crawler.Match{}, errors.Wrap(err, fmt.Sprintf("error parsing map stats %d", i))
		}
		mapStatsList = append(mapStatsList, mapStats)
	}

	match.MapStats = mapStatsList

	return match, nil
}
