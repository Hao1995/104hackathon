package models

// type JobUserScoreDBItem struct {
// 	ID        *int64 `json:"id"`
// 	JobNo     *int64 `json:"jobno"`
// 	WelfareNo *int64 `json:"welfare_no"`
// }

type JobUserScoreGetItem struct {
	Custno   *string                     `json:"custno"`
	CustName *string                     `json:"cust_name"`
	JobNo    *int64                      `json:"jobno"`
	JobName  *string                     `json:"job_name"`
	Score    *int64                      `json:"score"`
	Items    []*JobUserScoreGetChildItem `json:"items"`
	Error    *error                      `json:"error"`
}

type JobUserScoreGetChildItem struct {
	WelfareNo   *int64  `json:"welfare_no"`
	WelfareName *string `json:"welfare_name"`
	Score       *int64  `json:"score"`
}
