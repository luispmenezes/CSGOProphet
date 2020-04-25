package repository

import "csgo_prophet/model/repository"

type Store interface {
	AddMatch(matches *[]repository.Match) error
	GetMatch(pageIdx, pageSize int, filter map[string]interface{}) (repository.MatchQueryResult, error)
	RemoveMatch(filter map[string]interface{}) (int, error)

	AddEvent(events *[]repository.Event) error
	GetEvent(pageIdx, pageSize int, filter map[string]interface{}) (repository.EventQueryResult, error)
	RemoveEvent(filter map[string]interface{}) (int, error)
}
