package repository

import "csgo_prophet/model/web_crawler"

type Match struct {
	tableName struct{}          `pg:"csgo.match"`
	Id        string            `pg:"id,pk"`
	Data      web_crawler.Match `pg:"match_data"`
}
