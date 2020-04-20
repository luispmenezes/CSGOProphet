package web_crawler

import (
	"csgo_prophet/model/web_crawler"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func parseMapStats(mapStatsUrl string) (web_crawler.MapStats, error) {
	log.WithFields(log.Fields{"matchUrl": mapStatsUrl}).Info("Parsing match")

	document, err := NewRequest(http.MethodGet, mapStatsUrl, nil)

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "map stats overview web request failed")
	}

	mapStats := web_crawler.MapStats{Url: mapStatsUrl}

	//MAP NAME
	mapStats.MapName = document.Find(".stats-match-maps > .columns > :not(.inactive) > .stats-match-map-desc > .stats-match-map-result > .dynamic-map-name-full").First().Text()

	//TEAM 1 ROUNDS
	team1Rounds, err := strconv.Atoi(document.Find(".match-info-box > .team-left > div").First().Text())

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 1 fk")
	}

	mapStats.Team1Rounds = team1Rounds

	//TEAM 2 ROUNDS
	team2Rounds, err := strconv.Atoi(document.Find(".match-info-box > .team-right > div").First().Text())

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 2 rounds")
	}

	mapStats.Team2Rounds = team2Rounds

	// HALFS
	document.Find(".match-info-row").First().Find(".right > .ct-color, .t-color").Each(func(i int, s *goquery.Selection) {
		val, _ := strconv.Atoi(s.Text())

		if s.HasClass("ct-color") {
			if i%2 == 0 {
				mapStats.Team1CTRounds = val
			} else {
				mapStats.Team2TRounds = val
			}
		} else {
			if i%2 == 0 {
				mapStats.Team1TRounds = val
			} else {
				mapStats.Team1TRounds = val
			}
		}
	})

	// Rating
	ratingStrs := strings.Split(document.Find(".match-info-box-con > .match-info-row").Eq(1).Find(".right").First().Text(), " : ")

	team1Rating, err := strconv.ParseFloat(ratingStrs[0], 64)

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 1 rating")
	}

	mapStats.Team1Rating = team1Rating

	team2Rating, err := strconv.ParseFloat(ratingStrs[1], 64)

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 2 rating")
	}

	mapStats.Team2Rating = team2Rating

	// First Kills
	firstKillsStrs := strings.Split(document.Find(".match-info-box-con > .match-info-row").Eq(2).Find(".right").First().Text(), " : ")

	team1Fks, err := strconv.Atoi(firstKillsStrs[0])

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 1 fk")
	}

	mapStats.Team1Fks = team1Fks

	team2Fks, err := strconv.Atoi(firstKillsStrs[1])

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 2 fk")
	}

	mapStats.Team2Fks = team2Fks

	// Clutches
	clutchStrs := strings.Split(document.Find(".match-info-box-con > .match-info-row").Eq(3).Find(".right").First().Text(), " : ")

	team1Clutches, err := strconv.Atoi(clutchStrs[0])

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 1 clutches")
	}

	mapStats.Team1Clutches = team1Clutches

	team2Clutches, err := strconv.Atoi(clutchStrs[1])

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 2 fk")
	}

	mapStats.Team2Clutches = team2Clutches

	return mapStats, nil
}
