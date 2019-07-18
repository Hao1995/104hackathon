package models

type JobsJSONItem struct {
	// {"custno":"87a57d0d-f8da-440b-b6fd-8170f2c3ba42","jobno":"6179991","job":"業務行政助理","jobcat1":"2005003013","jobcat2":"0","jobcat3":"0","edu":60,"salary_low":990,"salary_high":990,"role":1,"language1":14444,"language2":1111,"language3":1111,"period":-1,"major_cat":"3006000000","major_cat2":"0","major_cat3":"0","industry":"1003001002","worktime":"週休二日","role_status":111,"s2":0,"s3":0,"addr_no":6001001009,"s9":1,"need_emp":1,"need_emp1":2,"startby":1,"exp_jobcat1":"0","exp_jobcat2":"0","exp_jobcat3":"0","description":"1.客戶資料處理\r\n2.商品資料處理\r\n3.報價資料處理\r\n4.訂單資料處理\r\n5.其他交辦事務\r\n6.具Photoshop CC 及 Photo Illustrator 經驗尤佳\r\n","others":"具Photoshop CC 及 Photo Illustrator 經驗尤佳"}
	Custno      *string `json:"custno"`
	Jobno       *string `json:"jobno"`
	Job         *string `json:"job"`
	Jobcat1     *string `json:"jobcat1"`
	Jobcat2     *string `json:"jobcat2"`
	Jobcat3     *string `json:"jobcat3"`
	Edu         *int64  `json:"edu"` // 最高學歷(用二進位儲存)(1: 高中以下 2:高中4:專科 8:大學16碩士 32:博士)
	SalaryLow   *int64  `json:"salary_low"`
	SalaryHigh  *int64  `json:"salary_high"`
	Role        *int64  // 1:全職2:兼職 3:高階4:殘障+全 5:殘障+兼 6:殘障+高
	Language1   *int64  `json:"language1"`
	Language2   *int64  `json:"language2"`
	Language3   *int64  `json:"language3"`
	WorkDur     *int64  `json:"period"`    // Working experience(working duration)
	MajorCat1   *string `json:"major_cat"` // 相關科系類別
	MajorCat2   *string `json:"major_cat2"`
	MajorCat3   *string `json:"major_cat3"`
	Industry    *string `json:"industry"`
	Vacation    *string `json:"worktime"`
	RoleStatus  *int64  `json:"role_status"` // 1:上班族 2:應屆畢 4:日間部 8:夜間部 16:尋找國防役 32:外籍人士 64:原住民 128:接受人力派遺-派遺單位專用 256:二度就業-派遺單位專用 512:soho
	Management  *int64  `json:"s2"`          // 管理責任：-1:該職務尚未填寫 0:無 1:4人以下 2:5-8人 3:9-12人 4:13人以上 5:未定
	BuinessTrip *int64  `json:"s3"`          // 是否出差：-1:該職務尚未填寫 0:無 1:1個月內 2:3個月以下 3:6個月以下 4:7個月以上 5:未定 (註.3)
	Addrno      *int64  `json:"addr_no"`
	WorkTime    *int64  `json:"s9"`        // 上班時段：0:預設 1:日班 2:晚班 4:大夜班 8:假日班
	NeedEmpLow  *int64  `json:"need_emp"`  // 最低需求人數
	NeedEmpHigh *int64  `json:"need_emp1"` // 最高需求人數
	StartBy     *int64  // 上班日期：0:一週內 1:兩週 2:一個月 3:不限 4:可年後上班
	ExpJobCat1  *string `json:"exp_jobcat1"` // 相關職務經驗
	ExpJobCat2  *string `json:"exp_jobcat2"`
	ExpJobCat3  *string `json:"exp_jobcat3"`
	Desc        *string `json:"description"` // 職務描述
	Others      *string `json:"others"`      // 其他條件
}
