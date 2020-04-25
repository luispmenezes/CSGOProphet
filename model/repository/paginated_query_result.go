package repository

type EventQueryResult struct {
	Page       int
	TotalPages int
	PageSize   int
	Result     []Event
}

type MatchQueryResult struct {
	Page       int
	TotalPages int
	PageSize   int
	Result     []Match
}
