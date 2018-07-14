package model

//Job 職務
type Job struct {
	Custno      string `json:"custno"`
	Jobno       string `json:"jobno"`
	Job         string `json:"job"`
	Jobcat1     string `json:"jobcat1"`
	Jobcat2     string `json:"jobcat2"`
	Jobcat3     string `json:"jobcat3"`
	Edu         int8   `json:"edu"`
	SalaryLow   int    `json:"salary_low"`
	SalaryHigh  int    `json:"salary_high"`
	Role        int8   `json:"role"`
	Language1   int32  `json:"language1"`
	Language2   int32  `json:"language2"`
	Language3   int32  `json:"language3"`
	Period      int8   `json:"period"`
	MajorCat    string `json:"major_cat"`
	MajorCat2   string `json:"major_cat2"`
	MajorCat3   string `json:"major_cat3"`
	Industry    string `json:"industry"`
	Worktime    string `json:"worktime"`
	RoleStatus  int16  `json:"role_status"`
	S2          int8   `json:"s2"`
	S3          int8   `json:"s3"`
	Addrno      int    `json:"addr_no"`
	S9          int8   `json:"s9"`
	NeedEmp     int    `json:"need_emp"`
	NeedEmp1    int    `json:"need_emp1"`
	Startby     int8   `json:"startby"`
	ExpJobcat1  string `json:"exp_jobcat1"`
	ExpJobcat2  string `json:"exp_jobcat2"`
	ExpJobcat3  string `json:"exp_jobcat3"`
	Description string `json:"description"`
	Others      string `json:"others"`
}

//Company 公司
type Company struct {
	Custno     string `json:"custno"`
	Invoice    int    `json:"invoice"`
	Name       string `json:"name"`
	Profile    string `json:"profile"`
	Management string `json:"management"`
	Welfare    string `json:"welfare"`
	Product    string `json:"product"`
}

//TrainClick 公司
type TrainClick struct {
	Action      string   `json:"action"`
	Jobno       string   `json:"jobno"`
	Date        string   `json:"date"`
	Joblist     []string `json:"joblist"`
	QueryString string   `json:"querystring"`
	Source      string   `json:"source"`
}

//TrainAction ...
type TrainAction struct {
	Action string `json:"action"`
	Jobno  string `json:"jobno"`
	Date   string `json:"date"`
	Source string `json:"source"`
	Device string `json:"device"`
}

//JobCategory ...
type JobCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Hide string `json:"hide"`
}

//Department ...
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Hide string `json:"hide"`
}

//District ...
type District struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Hide string `json:"hide"`
}

//Industry ...
type Industry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Hide string `json:"hide"`
}

//=====Score

//ScoreOriginData ...
type ScoreOriginData struct {
	Key            string `json:"key"`
	JobName        string `json:"jobname"`
	JobAction      string `json:"jobaction"`
	CompanyName    string `json:"companyname"`
	CompanyWelfare string `json:"companywelfare"`
	DistrictID     string `json:"districtid"`
	DistrictName   string `json:"districtname"`
}

//JobKey ...
type JobKey struct {
	Key string `json:"key"`
	Job string `json:"job"`
}

//JobScore ...
type JobScore struct {
	Job       string `json:"job"`
	Key       string `json:"key"`
	GoodScore int    `json:"goodscore"`
	BadScore  int    `json:"badscore"`
}

//AreaScore ...
type AreaScore struct {
	Area      string `json:"area"`
	Key       string `json:"key"`
	GoodScore int    `json:"goodscore"`
	BadScore  int    `json:"badscore"`
}

//QueryKey ...
type QueryKey struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	GoodScore int    `json:"goodscore"`
	BadScore  int    `json:"badscore"`
}

//Tag ...
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}
