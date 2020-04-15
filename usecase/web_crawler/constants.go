package web_crawler

const baseHost = "www.hltv.org"
const resultsPath = "/results"
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
const requestTimeout = 30

const resultsParamOffset = "offset"
const resultsParamStartDate = "startDate"
const resultsParamEndDate = "endDate"
const resultsParamStars = "stars"
const resultsParamContent = "content"
const resultsParamType = "matchType"
const resultsParamMap = "map"

const (
	ContentDEMO       = "demo"
	ContentHIGHLIGHTS = "highlights"
	ContentVOD        = "vod"
	ContentSTATS      = "stats"
)

const (
	MatchTypeLAN    = "Lan"
	MatchTypeONLINE = "Online"
)

const (
	MapCACHE       = "de_cache"
	MapSEASON      = "de_season"
	MapDUST2       = "de_dust2"
	MapMIRAGE      = "de_mirage"
	MapINFERNO     = "de_inferno"
	MapNUKE        = "de_nuke"
	MapTRAIN       = "de_train"
	MapCOBBLESTONE = "de_cobblestone"
	MapOVERPASS    = "de_overpass"
	MapTUSCAN      = "de_tuscan"
	MapVERTIGO     = "de_vertigo"
)
