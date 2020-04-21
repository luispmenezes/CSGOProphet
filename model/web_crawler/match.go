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
	Url              string
	MapName          string
	Team1Rounds      int
	Team1TRounds     int
	Team1CTRounds    int
	Team1Rating      float64
	Team1Fks         int
	Team1Clutches    int
	Team2Rounds      int
	Team2TRounds     int
	Team2CTRounds    int
	Team2Rating      float64
	Team2Fks         int
	Team2Clutches    int
	Team1PlayerStats []PlayerMapStats
	Team2PlayerStats []PlayerMapStats
	RoundDetails     []RoundDetail
}

type PlayerMapStats struct {
	Name         string
	Kills        int
	Headshots    int
	Assists      int
	FlashAssists int
	Deaths       int
	KAST         float64
	KDRatio      int
	ADR          float64
	FKDiff       int
	Rating       float64
}

type RoundDetail struct {
	Index           int
	Winner          int
	EquipmentValue1 int
	EquipmentValue2 int
}
