package models

type TrainClickJSONItem struct {
	Jobno       string
	Date        string
	Action      string
	Source      string
	Joblist     []string
	QueryString string
}
