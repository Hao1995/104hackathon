package models

type TrainActionJSONItem struct {
	// {"jobno":"10000835","date":"1527566791000","action":"viewJob","source":"web"}
	Jobno  string `json:"jobno"`
	Date   string `json:"date"`
	Action string `json:"action"`
	Source string `json:"source"`
	Device string `json:"device"`
}
