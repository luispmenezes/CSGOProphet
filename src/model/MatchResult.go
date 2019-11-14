package model

import "fmt"

type MatchResult struct {
	ResultURL string
	EventURL  string
	TimeStamp string
	Format    string
	Team1     Team
	Team2     Team
	Score1    int
	Score2    int
	Maps      []Map
}

func (m MatchResult) String() string {
	return fmt.Sprintf("[%s, %s, %s, %s, %v, %v, %d, %d, %v]",
		m.ResultURL, m.EventURL, m.TimeStamp, m.Format, m.Team1, m.Team2, m.Score1, m.Score2, m.Maps)
}
