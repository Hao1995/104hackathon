package models

type TrainActionJSONItem struct {
	// {"jobno":"10000835","date":"1527566791000","action":"viewJob","source":"web"}
	Jobno  string
	Date   string
	Action string
	Source string
	Device string
}
