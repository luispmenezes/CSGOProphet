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

	//Player Stats
	mapStats.Team1PlayerStats, err = parsePlayerStats(document.Find(".stats-table > tbody").First())

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 1 player stats")
	}

	mapStats.Team2PlayerStats, err = parsePlayerStats(document.Find(".stats-table > tbody").Last())

	if err != nil {
		return web_crawler.MapStats{}, errors.Wrap(err, "failed parsing team 2 player stats")
	}

	// Economy Round Details
	mapStats.RoundDetails, err = parseEconData(strings.ReplaceAll(mapStatsUrl, "/matches/", "/matches/economy/"))

	return mapStats, nil
}

func parsePlayerStats(statsRows *goquery.Selection) ([]web_crawler.PlayerMapStats, error) {
	var playerStats []web_crawler.PlayerMapStats

	statsRows.Find("tr").Each(func(i int, selection *goquery.Selection) {
		var currentPlayer web_crawler.PlayerMapStats

		// Name
		playerName, nameExists := selection.Find(".st-player > a").Attr("href")

		if !nameExists {
			log.Error(fmt.Sprintf("Failed to get player name %d", i))
			return
		}

		currentPlayer.Name = playerName

		// Kills
		killStr := strings.TrimSpace(selection.Find(".st-kills").Clone().Children().Remove().End().Text())
		kills, err := strconv.Atoi(killStr)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse kill %d - %s", i, killStr))
			return
		}

		currentPlayer.Kills = kills

		// HeadShots
		hsStr := selection.Find(".st-kills > .gtSmartphone-only").Clone().Children().Remove().End().Text()
		hsStr = hsStr[2 : len(hsStr)-1]
		hs, err := strconv.Atoi(hsStr)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse headshots %d - %s", i, hsStr))
			return
		}

		currentPlayer.Headshots = hs

		// Assists
		assistsStr := strings.TrimSpace(selection.Find(".st-assists").Clone().Children().Remove().End().Text())
		assists, err := strconv.Atoi(assistsStr)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse assists %d - %s", i, assistsStr))
			return
		}

		currentPlayer.Assists = assists

		// Flashes
		flashesStr := selection.Find(".st-kills > .gtSmartphone-only").Clone().Children().Remove().End().Text()
		flashesStr = flashesStr[2 : len(flashesStr)-1]
		flashAss, err := strconv.Atoi(flashesStr)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse flash assists %d - %s", i, flashesStr))
			return
		}

		currentPlayer.FlashAssists = flashAss

		// Deaths
		deathsStr := strings.TrimSpace(selection.Find(".st-deaths").Clone().Children().Remove().End().Text())
		deaths, err := strconv.Atoi(deathsStr)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse deaths %d - %s", i, deathsStr))
			return
		}

		currentPlayer.Deaths = deaths

		// KAST
		kastStr := strings.TrimSpace(strings.TrimSpace(selection.Find(".st-kdratio").Clone().Children().Remove().End().Text()))
		kast, err := strconv.ParseFloat(kastStr[:len(kastStr)-1], 64)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse kast %d - %s", i, kastStr))
			return
		}

		currentPlayer.KAST = kast

		// KD-DIFF
		kdDiffStr := strings.TrimSpace(strings.TrimSpace(selection.Find(".st-kddiff").Text()))
		kdDiff, err := strconv.Atoi(kdDiffStr)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse kDiff %d - %s", i, kdDiffStr))
			return
		}

		currentPlayer.KDRatio = kdDiff

		// ADR
		adrStr := strings.TrimSpace(strings.TrimSpace(selection.Find(".st-adr").Text()))
		adr, err := strconv.ParseFloat(adrStr, 64)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse adr %d - %s", i, adrStr))
			return
		}

		currentPlayer.ADR = adr

		// FK Diff
		fkDiffStr := strings.TrimSpace(strings.TrimSpace(selection.Find(".st-fkdiff").Text()))
		fkDiff, err := strconv.Atoi(fkDiffStr)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse fk diff %d - %s", i, fkDiffStr))
			return
		}

		currentPlayer.FKDiff = fkDiff

		// HLTV Rating
		ratingStr := strings.TrimSpace(strings.TrimSpace(selection.Find(".st-rating").Text()))
		rating, err := strconv.ParseFloat(ratingStr, 64)

		if err != nil {
			log.Error(fmt.Sprintf("Failed to parse rating %d - %s", i, ratingStr))
			return
		}

		currentPlayer.Rating = rating

		// Add to list
		playerStats = append(playerStats, currentPlayer)
	})

	if len(playerStats) != 5 {
		return []web_crawler.PlayerMapStats{}, errors.New("failed to parse player stats")
	}

	return playerStats, nil
}

func parseEconData(economyUrl string) ([]web_crawler.RoundDetail, error) {
	log.WithFields(log.Fields{"matchUrl": economyUrl}).Info("Parsing economy data ")

	document, err := NewRequest(http.MethodGet, economyUrl, nil)

	if err != nil {
		return []web_crawler.RoundDetail{}, errors.Wrap(err, "economic data web request failed")
	}

	var roundDetails []web_crawler.RoundDetail
	firstHalfData := document.Find(".equipment-categories").First()

	firstHalfData.Find("tbody > .team-categories").First().Find(".equipment-category-td").Each(func(i int, selection *goquery.Selection) {
		currentRound := web_crawler.RoundDetail{Index: i}

		eqValue, err := strconv.Atoi(selection.Text()[17:])

		if err != nil {
			currentRound.EquipmentValue1 = -1
		} else {
			currentRound.EquipmentValue1 = eqValue
		}

		imgSrc, _ := selection.Find(".equipment-category").Attr("src")

		if strings.HasSuffix(imgSrc, "Win.svg") {
			currentRound.Winner = 1
		} else {
			currentRound.Winner = 2
		}

		roundDetails = append(roundDetails, currentRound)
	})

	return roundDetails, nil
}
