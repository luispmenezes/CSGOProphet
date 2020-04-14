package demo

import "fmt"

type Team struct {
	TeamURL string
	Name    string
	Players []string
}

func (t Team) String() string {
	return fmt.Sprintf("[%s, %s, %s]", t.TeamURL, t.Name, t.Players)
}
