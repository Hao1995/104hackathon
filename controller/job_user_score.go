package controller

import (
	"fmt"
	"net/http"

	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

func JobUserScore(w http.ResponseWriter, req *http.Request) {

	ins := &JobUserScoreController{}
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

type JobUserScoreController struct {
}

func (c *JobUserScoreController) get(httpLib *utils.HTTPLib) {

	res := models.JobUserScoreGetItem{}

	// Init Params
	jobno, err := httpLib.Req.FormValueToNullInt64("jobno")
	if err != nil {
		tmp := err.Error()
		res.Error = &tmp
		httpLib.WriteJSON(res)
		return
	}
	if !jobno.Valid {
		err := fmt.Sprintf("'%v' is necessary", "jobno")
		res.Error = &err
		httpLib.WriteJSON(res)
		return
	}

	userID, err := httpLib.Req.FormValueToNullInt64("user_id")
	if err != nil {
		tmp := err.Error()
		res.Error = &tmp
		httpLib.WriteJSON(res)
		return
	}
	if !userID.Valid {
		err := fmt.Sprintf("'%v' is necessary", "user_id")
		res.Error = &err
		httpLib.WriteJSON(res)
		return
	}

	// - Company Name and Job Name and Total Score
	err = db.QueryRow("SELECT `C`.`custno`, `C`.`name` AS `cust_name`, `J`.`jobno` AS `jobno`, `J`.`job` AS `job_name`, `JUS`.`good_score`, `JUS`.`bad_score` FROM `job_user_score` AS `JUS`, `jobs` AS `J`, `companies` AS `C` WHERE 1 = 1 AND `JUS`.`jobno` = `J`.`jobno` AND `J`.`custno` = `C`.`custno` AND `JUS`.`jobno` = ? AND `JUS`.`user_id` = ?", jobno, userID).Scan(&res.Custno, &res.CustName, &res.JobNo, &res.JobName, &res.GoodScore, &res.BadScore)
	if err != nil {
		tmp := err.Error()
		logs.Error(tmp)
		res.Error = &tmp
		httpLib.WriteJSON(res)
		return
	}

	// - Get the Welfares and Each Score
	rows, err := db.Query("SELECT `W`.`id` AS `welfare_no`, `W`.`name` AS `welfare_name`, `WUS`.`score` FROM `job_welfares` AS `JW`, `welfares` AS `W`, `welfare_user_score` AS `WUS` WHERE 1 = 1 AND `JW`.`welfare_no` = `W`.`id` AND `W`.`id` = `WUS`.`welfare_no` AND `JW`.`jobno` = ? AND `WUS`.`user_id` = ? ORDER BY `welfare_no`", jobno, userID)
	if err != nil {
		tmp := err.Error()
		logs.Error(tmp)
		res.Error = &tmp
		httpLib.WriteJSON(err)
		return
	}
	defer rows.Close()

	items := []*models.JobUserScoreGetChildItem{}
	for rows.Next() {
		item := &models.JobUserScoreGetChildItem{}
		if err := rows.Scan(&item.WelfareNo, &item.WelfareName, &item.Score); err != nil {
			logs.Error(err)
			httpLib.WriteJSON(err)
			return
		}
		items = append(items, item)
		// logs.Debug("welfare_no:%v, welfare_name:%v, score:%v", *item.WelfareNo, *item.WelfareName, *item.Score)
	}

	res.Items = items

	httpLib.WriteJSON(res)
	return
}

func (c *JobUserScoreController) post(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	// - Init Parameters
	jobno, err := httpLib.Req.FormValueToNullInt64("jobno")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	userID, err := httpLib.Req.FormValueToNullInt64("user_id")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	// - Sync Data
	// Get job info
	rows, err := db.Query("SELECT `jobno` FROM `jobs` WHERE (IFNULL(?, 1) = 1 OR `jobno` = ?)", jobno, jobno)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer rows.Close()

	var jobs []int
	for rows.Next() {
		var job int
		if err := rows.Scan(&job); err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
		jobs = append(jobs, job)
	}

	// Get Users Info
	rows, err = db.Query("SELECT `id` FROM `users` WHERE (IFNULL(?, 1) = 1 OR `id` = ?)", userID, userID)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var user int
		if err := rows.Scan(&user); err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
		users = append(users, user)
	}

	// Get Score of each User and Job
	tx, err := db.Begin()
	for _, user := range users {
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
		for _, job := range jobs {
			// Calculate Score
			rows, err := tx.Query("SELECT `WUS`.`score` AS `score` FROM `job_welfares` AS `JW`, `welfare_user_score` AS `WUS` WHERE 1 = 1 AND `JW`.`welfare_no` = `WUS`.`welfare_no` AND `WUS`.`user_id` = ? AND `JW`.`jobno` = ?", user, job)
			if err != nil {
				logs.Error(err)
				res.Error = err.Error()
				httpLib.WriteJSON(res)
				return
			}
			defer rows.Close()

			var goodScore, badScore int64
			for rows.Next() {
				var score *int64
				if err := rows.Scan(&score); err != nil {
					logs.Error(err)
					res.Error = err.Error()
					httpLib.WriteJSON(res)
					return
				}

				switch {
				case *score > 0:
					goodScore += *score
				case *score < 0:
					badScore += *score
				}
			}

			// Insert Score Data.
			if _, err := tx.Exec("INSERT INTO `job_user_score` (`jobno`, `user_id`, `good_score`, `bad_score`) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE `good_score` = VALUES(`good_score`), `bad_score` = VALUES(`bad_score`);", job, user, goodScore, badScore); err != nil {
				logs.Error(err)
				res.Error = err.Error()
				httpLib.WriteJSON(res)
				return
			}
		}
	}
	tx.Commit()

	res.Message = fmt.Sprintf("Sucess synchronize the welfares score of %v jobs and %v users.", len(jobs), len(users))
	httpLib.WriteJSON(res)
	return
}

func (c *JobUserScoreController) delete(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	jobno, err := httpLib.Req.FormValueToNullInt64("jobno")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	if !jobno.Valid {
		err := fmt.Errorf("'%v' is necessary", "jobno")
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	// - Delete Data
	exeRes, err := db.Exec("DELETE FROM `job_user_score` WHERE `jobno` = ?", jobno)
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
