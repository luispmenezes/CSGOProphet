package web_crawler

import "fmt"

type DemoLink struct {
	DemoURL        string
	MatchResultURL string
	Timestamp      string
}

func (d DemoLink) String() string {
	return fmt.Sprintf("[%s, %s, %s]", d.DemoURL, d.MatchResultURL, d.Timestamp)
}
