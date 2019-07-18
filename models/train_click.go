package models

type TrainClickJSONItem struct {
	// {"jobno": "5314432", "date": "1527330445000", "action": "clickJob", "source": "web", "joblist": ["5314432", "3494303", "3859882", "4439691", "7416542", "7484247", "9063545", "8318922", "9431734", "9127494", "9904656", "9926223", "9950320", "10079727", "10378122", "6830159", "6856925", "10343328", "9730978", "8510715"], "querystring": "ro=0&jobcat=2005003003%2C2005003002%2C2005003005&kwop=7&keyword=%E6%97%A5%E6%96%87&area=6001001000%2C6001005000&order=1&asc=0&mode=s&jobsource=n104bank1"}
	Jobno       string   `json:"jobno"`
	Date        string   `json:"date"`
	Action      string   `json:"action"`
	Source      string   `json:"source"`
	Joblist     []string `json:"joblist"`
	QueryString string   `json:"querystring"`
}
