package controller

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/Hao1995/104hackathon/config"
	"github.com/Hao1995/104hackathon/glob"
	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

func Test() {
	fmt.Println("controller")
}

//HackathonJob ...
func HackathonJob(res http.ResponseWriter, req *http.Request) {

	//=====Params
	req.ParseForm()
	params := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "size":
			params[k] = strings.Join(v, "")
			// case "message":
			// 	params[k] = strings.Join(v, "")
		}
	}

	var rows *sql.Rows
	var err error
	if v, ok := params["size"]; ok {
		rows, err = db.Query("SELECT * FROM job LIMIT " + v.(string))
	} else {
		rows, err = db.Query("SELECT * FROM job LIMIT 100")
	}

	jobs := []*models.Job{}

	for rows.Next() {
		r := &models.Job{}

		err = rows.Scan(&r.Custno, &r.Jobno, &r.Job, &r.Jobcat1, &r.Jobcat2, &r.Jobcat3, &r.Edu, &r.SalaryLow, &r.SalaryHigh, &r.Role, &r.Language1, &r.Language2, &r.Language3, &r.Period, &r.MajorCat, &r.MajorCat2, &r.MajorCat3, &r.Industry, &r.Worktime, &r.RoleStatus, &r.S2, &r.S3, &r.Addrno, &r.S9, &r.NeedEmp, &r.NeedEmp1, &r.Startby, &r.ExpJobcat1, &r.ExpJobcat2, &r.ExpJobcat3, &r.Description, &r.Others)
		chechkErr(err)
		jobs = append(jobs, r)
	}

	jsonData, err := json.Marshal(jobs)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}

// SyncJobs :
// Sync the job data to DB
// * Delete all DB data
// * Insert JSON data to DB
func SyncJobs(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := models.APIRes{}

	// - Get the size of companies data
	var jobsIdx int
	rows, err := db.Query("SELECT COUNT(1) FROM `jobs`")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&jobsIdx)
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			js, err := json.Marshal(res)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(js)
			return
		}
	}
	err = rows.Err()
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	}
	logs.Debug("Will skip %v rows.", jobsIdx)

	// - Open Data File
	file, err := os.Open(config.CfgData.Data.Jobs)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	}
	defer file.Close()

	// - Scan Data Line By Line
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	v := reflect.ValueOf(models.JobsJSONItem{})
	size := glob.MySQLUpperPlaceholders / v.NumField()
	skipCount := 0
	totalCount := 0
	guard := make(chan struct{}, 2) // Max goroutines limit.
	errChan := make(chan bool)
	items := []models.JobsJSONItem{}
	for scanner.Scan() {
		// - Skip existing data.
		if skipCount < jobsIdx {
			skipCount++
			continue
		}

		// - Parse JSON to Item
		itemJSON := scanner.Text()
		item := models.JobsJSONItem{}
		json.Unmarshal([]byte(itemJSON), &item)
		// logs.Debug("JSON = %+v", item)

		items = append(items, item)
		totalCount++
		if len(items) >= size {
			// - Send items to channel and clear items, skipCount.
			guard <- struct{}{}
			wg.Add(1)
			go syncJobsInsertData(&wg, guard, errChan, items)
			logs.Debug("Send to insert %v data.", size)
			items = []models.JobsJSONItem{}
		}
	}
	if len(items) > 0 {
		// - Send the last data that not reach the size.
		guard <- struct{}{}
		wg.Add(1)
		go syncJobsInsertData(&wg, guard, errChan, items)
		logs.Debug("Last time send to insert %v data.", len(items))
	}

	// - Check Error of Scanner
	if err := scanner.Err(); err != nil {
		res.Error = err.Error()
		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	}

	logs.Debug("Waiting all goroutine")
	wg.Wait()

	// - Rev Err
	select {
	case <-errChan:
		res.Error = fmt.Sprintf("Something wrong !")
	default:
		res.Message = fmt.Sprintf("Skip %v data and insert %v data", skipCount, totalCount)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)

	return
}

func syncJobsInsertData(wg *sync.WaitGroup, guard chan struct{}, errChan chan bool, data []models.JobsJSONItem) {
	defer wg.Done()

	sqlStr := "INSERT INTO `jobs` (`custno`, `jobno`, `job`, `jobcat1`, `jobcat2`, `jobcat3`, `edu`, `salary_low`, `salary_high`, `role`, `language1`, `language2`, `language3`, `work_dur`, `major_cat1`, `major_cat2`, `major_cat3`, `industry`, `vacation`, `role_status`, `management`, `buiness_trip`, `addr_no`, `work_time`, `need_emp_low`, `need_emp_high`, `startby`, `exp_jobcat1`, `exp_jobcat2`, `exp_jobcat3`, `desc`, `others`) VALUES "
	vals := []interface{}{}
	var errInData error
	for _, item := range data {
		custno := utils.NewNullString(item.Custno)

		jobno, err := utils.StrToNewNullInt64(item.Jobno)
		if err != nil {
			errInData = err
			break
		}

		job := utils.NewNullString(item.Job)

		jobcat1, err := utils.StrToNewNullInt64(item.Jobcat1)
		if err != nil {
			errInData = err
			break
		}

		jobcat2, err := utils.StrToNewNullInt64(item.Jobcat2)
		if err != nil {
			errInData = err
			break
		}

		jobcat3, err := utils.StrToNewNullInt64(item.Jobcat3)
		if err != nil {
			errInData = err
			break
		}

		edu := utils.NewNullInt64(item.Edu)
		salaryLow := utils.NewNullInt64(item.SalaryLow)
		salaryHigh := utils.NewNullInt64(item.SalaryLow)
		role := utils.NewNullInt64(item.Role)
		language1 := utils.NewNullInt64(item.Language1)
		language2 := utils.NewNullInt64(item.Language2)
		language3 := utils.NewNullInt64(item.Language3)
		workDur := utils.NewNullInt64(item.WorkDur)

		majorCat1, err := utils.StrToNewNullInt64(item.MajorCat1)
		if err != nil {
			errInData = err
			break
		}

		majorCat2, err := utils.StrToNewNullInt64(item.MajorCat2)
		if err != nil {
			errInData = err
			break
		}

		majorCat3, err := utils.StrToNewNullInt64(item.MajorCat3)
		if err != nil {
			errInData = err
			break
		}

		industry, err := utils.StrToNewNullInt64(item.Industry)
		if err != nil {
			errInData = err
			break
		}

		vacation := utils.NewNullString(item.Vacation)
		roleStatus := utils.NewNullInt64(item.RoleStatus)
		management := utils.NewNullInt64(item.Management)
		buinessTrip := utils.NewNullInt64(item.BuinessTrip)

		if *item.Addrno == 0 {
			item.Addrno = nil
		}
		addrNo := utils.NewNullInt64(item.Addrno)

		workTime := utils.NewNullInt64(item.WorkTime)
		needEmpLow := utils.NewNullInt64(item.NeedEmpLow)
		needEmpHigh := utils.NewNullInt64(item.NeedEmpHigh)
		startby := utils.NewNullInt64(item.StartBy)

		expJobcat1, err := utils.StrToNewNullInt64(item.ExpJobCat1)
		if err != nil {
			errInData = err
			break
		}

		expJobcat2, err := utils.StrToNewNullInt64(item.ExpJobCat2)
		if err != nil {
			errInData = err
			break
		}

		expJobcat3, err := utils.StrToNewNullInt64(item.ExpJobCat3)
		if err != nil {
			errInData = err
			break
		}

		desc := utils.NewNullString(item.Desc)
		others := utils.NewNullString(item.Others)

		sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, custno, jobno, job, jobcat1, jobcat2, jobcat3, edu, salaryLow, salaryHigh, role, language1, language2, language3, workDur, majorCat1, majorCat2, majorCat3, industry, vacation, roleStatus, management, buinessTrip, addrNo, workTime, needEmpLow, needEmpHigh, startby, expJobcat1, expJobcat2, expJobcat3, desc, others)
	}
	if errInData != nil {
		logs.Error(errInData)
		select {
		case errChan <- true:
		default:
			<-guard
			return
		}
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// logs.Debug(sqlStr)
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		logs.Error(err)
		select {
		case errChan <- true:
		default:
			<-guard
			return
		}
	}
	defer stmt.Close()
	// logs.Debug(vals)
	_, err = stmt.Exec(vals...)
	if err != nil {
		logs.Error(err)
		select {
		case errChan <- true:
		default:
			<-guard
			return
		}
	}

	<-guard
	return
}
