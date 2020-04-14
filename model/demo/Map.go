package demo

import "fmt"

type Map struct {
	MapName string
	Score1  int
	Score2  int
	Rounds  []Round
}

func (m Map) String() string {
	return fmt.Sprintf("[%s, %d, %d, %v]", m.MapName, m.Score1, m.Score2, m.Rounds)
}
