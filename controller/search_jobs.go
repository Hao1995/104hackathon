package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

// SearchJobs :
// Get > Search the Jobs by "key". ex:業務、後端
func SearchJobs(w http.ResponseWriter, req *http.Request) {

	httpLib := &utils.HTTPLib{}
	httpLib.Init(w, req)
	ins := &SearchJobsController{httpLib: httpLib}

	res := models.APIRes{}

	switch req.Method {
	case http.MethodGet:
		ins.get()
	// case http.MethodPost:
	// 	ins.post(httpLib)
	// case http.MethodDelete:
	// 	ins.delete(httpLib)
	default:
		res.Error = fmt.Sprintf("There is no way correspond to method '%v'", req.Method)
		httpLib.WriteJSON(res)
		return
	}
}

type SearchJobsController struct {
	httpLib *utils.HTTPLib
	params  models.SearchJobsParams
}

func (c *SearchJobsController) get() {

	res := make(map[string]interface{})
	res["country"] = nil
	res["jobList"] = nil
	res["Error"] = nil

	// Init Params
	if err := c.init(); err != nil {
		logs.Error(err)
		res["Error"] = err.Error()
		c.httpLib.WriteJSON(res)
		return
	}

	// - Find The Jobs That Corresponding The Key
	now := time.Now()
	jobMap, err := c.findRelativeJobs()
	if err != nil {
		logs.Error(err)
		res["Error"] = err.Error()
		c.httpLib.WriteJSON(res)
		return
	}
	logs.Info("findRelativeJobs = %v", time.Since(now))

	// - Find Length Of `jobs`
	now = time.Now()
	countJobs, err := c.countAllJobs()
	if err != nil {
		logs.Error(err)
		res["Error"] = err.Error()
		c.httpLib.WriteJSON(res)
		return
	}
	logs.Info("countAllJobs = %v", time.Since(now))

	// - Get the Jobs Score and Calculate It's PR.
	countryScore, jobListScore, err := c.getJobsPR(jobMap, countJobs)
	if err != nil {
		logs.Error(err)
		res["Error"] = err.Error()
		c.httpLib.WriteJSON(res)
		return
	}

	res["country"] = countryScore
	res["jobList"] = jobListScore

	c.httpLib.WriteJSON(res)
	return
}

func (c *SearchJobsController) init() error {

	c.params.Key = c.httpLib.Req.FormValueToNullString("key")
	if !c.params.Key.Valid {
		return fmt.Errorf("'%v' is necessary", "key")
	}

	addrNoTmp, err := c.httpLib.Req.FormValueToNullInt64("addr_no")
	if err != nil {
		return err
	}
	c.params.AddrNo = utils.NewNullString(nil)
	if addrNoTmp.Valid && utils.CountDigits(addrNoTmp.Int64) == 7 {
		tmp := strconv.FormatInt(addrNoTmp.Int64, 10) + "%" // 6001016 >> "6001016%"
		c.params.AddrNo = utils.NewNullString(&tmp)
	} else {
		logs.Debug("'addr_no' is '%v', make addr_no be null.", addrNoTmp)
	}

	c.params.UserID, err = c.httpLib.Req.FormValueToNullInt64("user_id") // pretend has user login
	if err != nil {
		return err
	}
	if !c.params.UserID.Valid {
		var tmp int64
		if err := db.QueryRow("SELECT `id` FROM `users` WHERE `email` = 'default@gmail.com'").Scan(&tmp); err != nil {
			return err
		}
		c.params.UserID = utils.NewNullInt64(&tmp)
	}

	c.params.Pi, err = c.httpLib.Req.FormValueToNullInt64("pi")
	if err != nil {
		return err
	}
	if !c.params.Pi.Valid {
		val := int64(1)
		c.params.Pi = utils.NewNullInt64(&val)
	}
	// Check 'pi' > 0
	if c.params.Pi.Int64 <= 0 {
		return fmt.Errorf("'pi' needs greater than 0")
	}

	c.params.Ps, err = c.httpLib.Req.FormValueToNullInt64("ps")
	if err != nil {
		return err
	}
	if !c.params.Ps.Valid {
		val := int64(6)
		c.params.Ps = utils.NewNullInt64(&val)
	}
	// Check 'ps' > 0
	if c.params.Ps.Int64 <= 0 {
		return fmt.Errorf("'ps' needs greater than 0")
	}

	return nil
}

func (c *SearchJobsController) findRelativeJobs() (map[int64]bool, error) {

	jobMap := make(map[int64]bool)

	rows, err := db.Query("SELECT `jobno`, `joblist` FROM `train_click` WHERE `key` = ?", c.params.Key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var job int64
		var jobList string
		err := rows.Scan(&job, &jobList)
		if err != nil {
			return nil, err
		}
		if _, ok := jobMap[job]; !ok {
			jobMap[job] = true
		}

		jobs := strings.Split(jobList, ",")
		for _, val := range jobs {
			num, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			if _, ok := jobMap[num]; !ok {
				jobMap[num] = true
				// fmt.Printf("%v, ", num)
			}
		}
	}
	logs.Debug("Key jobs = %v", len(jobMap))

	return jobMap, nil
}

func (c *SearchJobsController) countAllJobs() (int64, error) {

	var countJobs int64
	if err := db.QueryRow("SELECT COUNT(1) FROM `jobs`").Scan(&countJobs); err != nil {
		return countJobs, err
	}
	logs.Debug("All job = %v", countJobs)

	return countJobs, nil
}

func (c *SearchJobsController) getJobsPR(jobMap map[int64]bool, countJobs int64) (models.SearchJobsScoreItem, []models.SearchJobsListItem, error) {

	// Get the Jobs Score
	now := time.Now()
	jobsScore, goodScore, badScore, err := c.getJobsScoreAndTotalScore(jobMap)
	if err != nil {
		return models.SearchJobsScoreItem{}, nil, err
	}
	logs.Info("getJobsScoreAndTotalScore = %v", time.Since(now))

	countryScore, jobListScore, err := c.getJobsPRAndTotalPR(jobsScore, goodScore, badScore, countJobs)
	if err != nil {
		return countryScore, nil, err
	}

	return countryScore, jobListScore, err
}

func (c *SearchJobsController) getJobsScoreAndTotalScore(jobMap map[int64]bool) ([]models.JobUserScoreGetItem, int64, int64, error) {

	jobsScore := []models.JobUserScoreGetItem{}
	var goodScore, badScore int64

	stmt, err := db.Prepare("SELECT `J`.`jobno`, `JUS`.`good_score`, `JUS`.`bad_score` FROM `job_user_score` AS `JUS`, `jobs` AS `J`, `districts` AS `D` WHERE `JUS`.`jobno` = `J`.`jobno` AND `J`.`addr_no` = `D`.`id` AND `JUS`.`user_id` = ? AND (IFNULL(?, 1) = 1 OR `D`.`id` LIKE ?) AND `J`.`jobno` IN (?" + strings.Repeat(",?", len(jobMap)-1) + ") ORDER BY `good_score` DESC")
	if err != nil {
		return nil, goodScore, badScore, err
	}
	defer stmt.Close()

	vals := []interface{}{c.params.UserID, c.params.AddrNo, c.params.AddrNo}
	for key := range jobMap {
		vals = append(vals, key)
	}
	rows, err := stmt.Query(vals...)
	if err != nil {
		return nil, goodScore, badScore, err
	}
	defer rows.Close()

	for rows.Next() {
		jobScore := models.JobUserScoreGetItem{}
		if err := rows.Scan(&jobScore.JobNo, &jobScore.GoodScore, &jobScore.BadScore); err != nil {
			return nil, goodScore, badScore, err
		}
		jobsScore = append(jobsScore, jobScore)

		goodScore += *jobScore.GoodScore
		badScore += *jobScore.BadScore
	}
	logs.Debug("Get %v jobs score. GoodScore = %v, BadScore = %v", len(jobsScore), goodScore, badScore)

	return jobsScore, goodScore, badScore, nil
}

// func (c *SearchJobsController) getJobsPRAndTotalPR(jobsScore []models.JobUserScoreGetItem, goodScore, badScore, countJobs int64) (models.SearchJobsScoreItem, []models.SearchJobsListItem, error) {

// 	countryScore := models.SearchJobsScoreItem{}
// 	jobListScore := []models.SearchJobsListItem{}

// 	// Calculate Country PR
// 	var avgGoodScore, avgBadScore float64
// 	if tmpLen := len(jobsScore); tmpLen != 0 {
// 		avgGoodScore = float64(goodScore) / float64(tmpLen)
// 		avgBadScore = float64(badScore) / float64(tmpLen)
// 	}
// 	logs.Debug("Average good-score = %v, bad-score = %v", avgGoodScore, avgBadScore)

// 	now := time.Now()
// 	var goodPR, badPR float64
// 	if err := db.QueryRow("SELECT COUNT(1) AS `count` FROM `job_user_score` WHERE `user_id` = ? AND `good_score` < ?", c.params.UserID, avgGoodScore).Scan(&goodPR); err != nil {
// 		return countryScore, nil, err
// 	}
// 	goodPR /= float64(countJobs)
// 	countryScore.GoodScore = &goodPR
// 	if err := db.QueryRow("SELECT COUNT(1) AS `count` FROM `job_user_score` WHERE `user_id` = ? AND `bad_score` > ?", c.params.UserID, avgBadScore).Scan(&badPR); err != nil {
// 		return countryScore, nil, err
// 	}
// 	badPR /= float64(countJobs)
// 	countryScore.BadScore = &badPR
// 	logs.Debug("Country Good PR = %v, Bad PR = %v", goodPR, badPR)
// 	logs.Info("getPRofTotal = %v", time.Since(now))

// 	// Calculate jobList PR
// 	now = time.Now()
// 	startIdx := (c.params.Pi.Int64 - 1) * c.params.Ps.Int64
// 	endIdx := c.params.Pi.Int64 * c.params.Ps.Int64
// 	for _, item := range jobsScore[startIdx:endIdx] {
// 		jobScore := models.SearchJobsListItem{}
// 		var goodPR, badPR float64
// 		if err := db.QueryRow("SELECT COUNT(1) AS `count` FROM `job_user_score` WHERE `user_id` = ? AND `good_score` < ?", c.params.UserID, item.GoodScore).Scan(&goodPR); err != nil {
// 			return countryScore, nil, err
// 		}
// 		goodPR /= float64(countJobs)
// 		jobScore.GoodScore = &goodPR
// 		if err := db.QueryRow("SELECT COUNT(1) AS `count` FROM `job_user_score` WHERE `user_id` = ? AND `bad_score` > ?", c.params.UserID, item.BadScore).Scan(&badPR); err != nil {
// 			return countryScore, nil, err
// 		}
// 		badPR /= float64(countJobs)
// 		jobScore.BadScore = &badPR

// 		if err := db.QueryRow("SELECT `J`.`job` AS `job_name`, `C`.`name` AS `cust_name` FROM `jobs` AS `J`, `companies` AS `C` WHERE 1 = 1 AND `J`.`custno` = `C`.`custno` AND `J`.`jobno` = ?", item.JobNo).Scan(&jobScore.JobName, &jobScore.JobCompany); err != nil {
// 			return countryScore, nil, err
// 		}

// 		jobListScore = append(jobListScore, jobScore)
// 	}
// 	logs.Info("getPRofJobs = %v", time.Since(now))

// 	return countryScore, jobListScore, nil
// }

func (c *SearchJobsController) getJobsPRAndTotalPR(jobsScore []models.JobUserScoreGetItem, goodScore, badScore, countJobs int64) (models.SearchJobsScoreItem, []models.SearchJobsListItem, error) {

	countryScore := models.SearchJobsScoreItem{}
	jobListScore := []models.SearchJobsListItem{}

	// Calculate Country PR
	var avgGoodScore, avgBadScore float64
	if tmpLen := len(jobsScore); tmpLen != 0 {
		avgGoodScore = float64(goodScore) / float64(tmpLen)
		avgBadScore = float64(badScore) / float64(tmpLen)
	}
	logs.Debug("Average good-score = %v, bad-score = %v", avgGoodScore, avgBadScore)

	// - Load all score data.
	now := time.Now()
	// Good
	rows, err := db.Query("SELECT `good_score`, COUNT(1) AS `count` from `job_user_score` WHERE `user_id` = ? GROUP BY `good_score`", c.params.UserID)
	if err != nil {
		return countryScore, nil, err
	}
	defer rows.Close()

	goodScores := make(map[int]int)
	for rows.Next() {
		var score, count int
		if err := rows.Scan(&score, &count); err != nil {
			return countryScore, nil, err
		}
		goodScores[score] = count
	}

	// Bad
	rows, err = db.Query("SELECT `bad_score`, COUNT(1) AS `count` from `job_user_score` WHERE `user_id` = ? GROUP BY `bad_score`", c.params.UserID)
	if err != nil {
		return countryScore, nil, err
	}
	defer rows.Close()

	badScores := make(map[int]int)
	for rows.Next() {
		var score, count int
		if err := rows.Scan(&score, &count); err != nil {
			return countryScore, nil, err
		}
		badScores[score] = count
	}
	logs.Info("goodScores & badScores = %v", time.Since(now))

	// - Calculate AVG PR
	now = time.Now()
	var goodCount, badCount int
	for score, count := range goodScores {
		if float64(score) < avgGoodScore {
			goodCount += count
		}
	}
	goodPR := float64(goodCount) / float64(countJobs)
	countryScore.GoodScore = &goodPR
	for score, count := range badScores {
		if float64(score) > avgBadScore {
			badCount += count
		}
	}
	badPR := float64(badCount) / float64(countJobs)
	countryScore.BadScore = &badPR
	logs.Debug("Country Good PR = %v, Bad PR = %v", goodPR, badPR)
	logs.Info("getPRofTotal = %v", time.Since(now))

	// - Calculate Jobs PR
	now = time.Now()
	startIdx := (c.params.Pi.Int64 - 1) * c.params.Ps.Int64
	endIdx := c.params.Pi.Int64 * c.params.Ps.Int64
	for _, item := range jobsScore[startIdx:endIdx] {
		jobScore := models.SearchJobsListItem{}

		var goodCount, badCount int
		for score, count := range goodScores {
			if int64(score) < *item.GoodScore {
				goodCount += count
			}
		}
		goodPR := float64(goodCount) / float64(countJobs)
		jobScore.GoodScore = &goodPR
		for score, count := range badScores {
			if int64(score) > *item.BadScore {
				badCount += count
			}
		}
		badPR := float64(badCount) / float64(countJobs)
		jobScore.BadScore = &badPR

		if err := db.QueryRow("SELECT `J`.`job` AS `job_name`, `C`.`name` AS `cust_name` FROM `jobs` AS `J`, `companies` AS `C` WHERE 1 = 1 AND `J`.`custno` = `C`.`custno` AND `J`.`jobno` = ?", item.JobNo).Scan(&jobScore.JobName, &jobScore.JobCompany); err != nil {
			return countryScore, nil, err
		}

		jobListScore = append(jobListScore, jobScore)
	}
	logs.Info("getPRofJobs = %v", time.Since(now))

	return countryScore, jobListScore, nil
}
