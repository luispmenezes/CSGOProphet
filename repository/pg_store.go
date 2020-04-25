package repository

import (
	"csgo_prophet/model/repository"
	"fmt"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"math"
)

const defaultPageSize = 50

type PGStore struct {
	database *pg.DB
}

func NewPGStore(host, port, database, username, password string) *PGStore {
	log.Printf("Starting Persistence manager (host:%s, port:%s, db:%s user:%s)", host, port, database, username)
	return &PGStore{database: pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Database: database,
		User:     username,
		Password: password,
	})}
}

func (P *PGStore) AddMatch(matches *[]repository.Match) error {
	return P.database.Insert(matches)
}

func (P *PGStore) GetMatch(pageIdx, pageSize int, filter map[string]interface{}) (repository.MatchQueryResult, error) {
	if pageSize == 0 || pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}

	recordCount, err := P.database.Model((*repository.Event)(nil)).Count()

	if err != nil {
		return repository.MatchQueryResult{}, err
	}

	totalPages := int(math.Ceil(float64(recordCount) / float64(pageSize)))

	if pageIdx >= totalPages {
		return repository.MatchQueryResult{
			Page:       pageIdx,
			TotalPages: totalPages,
			PageSize:   pageSize,
		}, nil
	}

	data := []repository.Match{}
	query := P.database.Model(&data).Order("open_time ASC")

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	err = query.Limit(pageSize).Offset(pageIdx * pageSize).Select()

	return repository.MatchQueryResult{
		Page:       pageIdx,
		TotalPages: totalPages,
		PageSize:   pageSize,
		Result:     data,
	}, err
}

func (P *PGStore) RemoveMatch(filter map[string]interface{}) (int, error) {
	query := P.database.Model(&repository.Match{})

	if len(filter) > 0 {
		for key, value := range filter {
			query = query.Where(key+" = ?", value)
		}
	} else {
		query = query.Where("TRUE")
	}

	res, err := query.Delete()

	if err == nil && res != nil {
		return res.RowsAffected(), err
	}

	return 0, err
}

func (P *PGStore) AddEvent(events *[]repository.Event) error {
	return P.database.Insert(events)
}

func (P *PGStore) GetEvent(pageIdx, pageSize int, filter map[string]interface{}) (repository.EventQueryResult, error) {
	if pageSize == 0 || pageSize > defaultPageSize {
		pageSize = defaultPageSize
	}

	recordCount, err := P.database.Model((*repository.Event)(nil)).Count()

	if err != nil {
		return repository.EventQueryResult{}, err
	}

	totalPages := int(math.Ceil(float64(recordCount) / float64(pageSize)))

	if pageIdx >= totalPages {
		return repository.EventQueryResult{
			Page:       pageIdx,
			TotalPages: totalPages,
			PageSize:   pageSize,
		}, nil
	}

	data := []repository.Event{}
	query := P.database.Model(&data).Order("open_time ASC")

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	err = query.Limit(pageSize).Offset(pageIdx * pageSize).Select()

	return repository.EventQueryResult{
		Page:       pageIdx,
		TotalPages: totalPages,
		PageSize:   pageSize,
		Result:     data,
	}, err
}

func (P *PGStore) RemoveEvent(filter map[string]interface{}) (int, error) {
	query := P.database.Model(&repository.Event{})

	if len(filter) > 0 {
		for key, value := range filter {
			query = query.Where(key+" = ?", value)
		}
	} else {
		query = query.Where("TRUE")
	}

	res, err := query.Delete()

	if err == nil && res != nil {
		return res.RowsAffected(), err
	}

	return 0, err
}
