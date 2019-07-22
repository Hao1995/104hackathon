package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

func JobWelfares(w http.ResponseWriter, req *http.Request) {

	ins := &JobWelfaresController{}
	httpLib := &utils.HTTPLib{}
	httpLib.Init(w, req)

	res := models.APIRes{}

	switch req.Method {
	case http.MethodGet:
		ins.get(httpLib)
	case http.MethodPost:
		ins.post(httpLib)
	case http.MethodDelete:
		ins.delete(httpLib)
	default:
		res.Error = fmt.Sprintf("There is no way correspond to method '%v'", req.Method)
		httpLib.WriteJSON(res)
		return
	}
}

type JobWelfaresController struct {
}

func (c *JobWelfaresController) get(httpLib *utils.HTTPLib) {

	res := models.JobWelfaresGetItem{}

	jobNo, err := httpLib.Req.FormValueToNullInt64("jobno")
	if err != nil {
		res.Error = &err
		httpLib.WriteJSON(res)
		return
	}
	if !jobNo.Valid {
		err := fmt.Errorf("'%v' is necessary", "jobno")
		res.Error = &err
		httpLib.WriteJSON(res)
		return
	}

	// - Company Name and Job Name
	row := db.QueryRow("SELECT `C`.`custno` AS `custno`, `C`.`name` AS `cust_name`, `J`.`jobno` AS `jobno` , `J`.`job` AS `job_name` FROM `job_welfares` AS `JW`, `jobs` AS `J`, `companies` AS `C` WHERE `JW`.`jobno` = `J`.`jobno` AND `J`.`custno` = `C`.`custno` AND `JW`.`jobno` = ?", jobNo)
	err = row.Scan(&res.Custno, &res.CustName, &res.JobNo, &res.JobName)
	if err != nil {
		logs.Error(err)
		res.Error = &err
		httpLib.WriteJSON(res)
		return
	}

	// - Get the Welfares
	stmt, err := db.Prepare("SELECT `W`.`id` AS `welare_no`, `W`.`name` AS `welare_name` FROM `job_welfares` AS `JW`, `welfares` AS `W` WHERE `JW`.`welfare_no` = `W`.`id` AND `JW`.`jobno` = ?")
	if err != nil {
		logs.Error(err)
		res.Error = &err
		httpLib.WriteJSON(err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(jobNo)
	if err != nil {
		logs.Error(err)
		res.Error = &err
		httpLib.WriteJSON(err)
		return
	}
	defer rows.Close()

	items := []*models.JobWelfaresGetChildItem{}
	for rows.Next() {
		item := &models.JobWelfaresGetChildItem{}
		err := rows.Scan(&item.WelfareNo, &item.WelfareName)
		if err != nil {
			logs.Error(err)
			httpLib.WriteJSON(err)
			return
		}
		items = append(items, item)
		// logs.Debug("welfare_no:%v, welfare_name:%v", *item.WelfareNo, *item.WelfareName)
	}

	res.Items = items

	httpLib.WriteJSON(res)
	return
}

func (c *JobWelfaresController) post(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	jobno, err := httpLib.Req.FormValueToNullInt64("jobno")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	// - Sync Data
	tx, err := db.Begin()
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	// Get job info
	rows, err := tx.Query("SELECT `jobno`, `desc`, `others` FROM `jobs` WHERE (IFNULL(?, 1) = 1 OR `jobno` = ?)", jobno, jobno)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer rows.Close()

	jobs := []models.JobsDBItem{}
	for rows.Next() {
		job := models.JobsDBItem{}
		err := rows.Scan(&job.Jobno, &job.Desc, &job.Others)
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
		jobs = append(jobs, job)
	}

	// Get welfares
	rows, err = tx.Query("SELECT `id`,`name` FROM `welfares`")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer rows.Close()

	welfares := []models.WelfaresItem{}
	for rows.Next() {
		welfare := models.WelfaresItem{}
		err := rows.Scan(&welfare.ID, &welfare.Name)
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
		welfares = append(welfares, welfare)
	}

	// Insert. Does the job include those welfares ?
	stmt, err := tx.Prepare("INSERT INTO `job_welfares` (`jobno`, `welfare_no`) VALUES (?, ?)")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer stmt.Close()

	for _, job := range jobs {
		// logs.Debug("Job:%v.", *job.Jobno)
		for _, welfare := range welfares {
			if job.Desc != nil && strings.Contains(*job.Desc, *welfare.Name) || job.Others != nil && strings.Contains(*job.Others, *welfare.Name) {
				// logs.Debug("Welfare: %v", *welfare.Name)
				if _, err := stmt.Exec(job.Jobno, welfare.ID); err != nil {
					logs.Error(err)
					res.Error = err.Error()
					httpLib.WriteJSON(res)
					return
				}
			}
		}
	}

	tx.Commit()

	res.Message = fmt.Sprintf("Sucess synchronize the welfares of %v jobs.", len(jobs))
	httpLib.WriteJSON(res)
	return
}

func (c *JobWelfaresController) delete(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	jobno, err := httpLib.Req.FormValueToNullInt64("jobno")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	// - Delete Data
	exeRes, err := db.Exec("DELETE FROM `job_welfares` WHERE `jobno` = ?", jobno)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	num, err := exeRes.RowsAffected()
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	res.Message = fmt.Sprintf("Delete %v data", num)
	httpLib.WriteJSON(res)
	return
}
