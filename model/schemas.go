package model

//Job 職務
type Job struct {
	Custno      string `json:"custno"`
	Jobno       int    `json:"jobno"`
	Job         string `json:"job"`
	Jobcat1     int    `json:"jobcat1"`
	Jobcat2     int    `json:"jobcat2"`
	Jobcat3     int    `json:"jobcat3"`
	Edu         int8   `json:"edu"`
	SalaryLow   int    `json:"salary_low"`
	SalaryHigh  int    `json:"salary_high"`
	Role        int8   `json:"role"`
	Language1   int32  `json:"language1"`
	Language2   int32  `json:"language2"`
	Language3   int32  `json:"language3"`
	Period      int8   `json:"period"`
	MajorCat    int    `json:"major_cat"`
	MajorCat2   int    `json:"major_cat2"`
	MajorCat3   int    `json:"major_cat3"`
	Industry    int    `json:"industry"`
	Worktime    string `json:"worktime"`
	RoleStatus  int8   `json:"role_status"`
	S2          int8   `json:"s2"`
	S3          int8   `json:"s3"`
	Addrno      int    `json:"addrno"`
	S9          int8   `json:"s9"`
	NeedEmp     int    `json:"need_emp"`
	NeedEmp1    int    `json:"need_emp1"`
	Startby     int8   `json:"startby"`
	ExpJobcat1  int    `json:"exp_jobcat1"`
	ExpJobcat2  int    `json:"exp_jobcat2"`
	ExpJobcat3  int    `json:"exp_jobcat3"`
	Description string `json:"description"`
	Others      string `json:"others"`
}
