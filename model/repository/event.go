package repository

import "csgo_prophet/model/web_crawler"

type Event struct {
	tableName struct{}          `pg:"csgo.event"`
	Id        string            `pg:"id,pk"`
	Data      web_crawler.Event `pg:"event_data"`
}
