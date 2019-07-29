package models

import "database/sql"

type SearchJobsParams struct {
	Key    sql.NullString
	UserID sql.NullInt64
	AddrNo sql.NullString
	Pi     sql.NullInt64
	Ps     sql.NullInt64
}

type SearchJobsScoreItem struct {
	GoodScore *float64 `json:"goodScore"`
	BadScore  *float64 `json:"badScore"`
}

type SearchJobsWelfaresItem struct {
	WelfareNo   *int64  `json:"welfare_no"`
	WelfareName *string `json:"welfare_name"`
}

type SearchJobsListItem struct {
	JobName    *string `json:"jobName"`
	JobCompany *string `json:"jobCompany"`
	SearchJobsScoreItem
}
