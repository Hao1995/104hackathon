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
