package models

type JobWelfaresDBItem struct {
	ID        *int64 `json:"id"`
	JobNo     *int64 `json:"jobno"`
	WelfareNo *int64 `json:"welfare_no"`
}

type JobWelfaresGetItem struct {
	Custno   *string                    `json:"custno"`
	CustName *string                    `json:"cust_name"`
	JobNo    *int64                     `json:"jobno"`
	JobName  *string                    `json:"job_name"`
	Items    []*JobWelfaresGetChildItem `json:"items"`
	Error    *error                     `json:"error"`
}

type JobWelfaresGetChildItem struct {
	WelfareNo   *int64  `json:"welfare_no"`
	WelfareName *string `json:"welfare_name"`
}
