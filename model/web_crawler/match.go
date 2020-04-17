package web_crawler

import "time"

type Match struct {
	Url              string
	Format           int
	StartTime        time.Time
	Event            string
	Team1            string
	Team1Composition []string
	Team1Score       int
	Team2            string
	Team2Composition []string
	Team2Score       int
	DemoUrl          string
	MapStats         []MapStats
}

type MapStats struct {
	MapName       string
	Team1Rounds   int
	Team1TRounds  int
	Team1CTRounds int
	Team1Rating   int
	Team1Fks      int
	Team1Clutches int
	Team2Rounds   int
	Team2TRounds  int
	Team2CTRounds int
	Team2Rating   int
	Team2Fks      int
	Team2Clutches int
	PlayerStats1  []PlayerMapStats
	PlayerStats2  []PlayerMapStats
}

type PlayerMapStats struct {
	Kills        int
	Headshots    int
	Assists      int
	FlashAssists int
	Deaths       int
	KAST         int
	KDRatio      int
	ADR          int
	FKDiff       int
	Rating       int
}
