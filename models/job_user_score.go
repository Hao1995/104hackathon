package models

type JobUserScoreGetItem struct {
	Custno    *string                     `json:"custno"`
	CustName  *string                     `json:"cust_name"`
	JobNo     *int64                      `json:"jobno"`
	JobName   *string                     `json:"job_name"`
	GoodScore *int64                      `json:"good_score"`
	BadScore  *int64                      `json:"bad_score"`
	Items     []*JobUserScoreGetChildItem `json:"items"`
	Error     *string                     `json:"error"`
}

type JobUserScoreGetChildItem struct {
	WelfareNo   *int64  `json:"welfare_no"`
	WelfareName *string `json:"welfare_name"`
	Score       *int64  `json:"score"`
}
