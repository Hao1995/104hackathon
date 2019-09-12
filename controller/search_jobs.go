package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hao1995/104hackathon/cache"
	"github.com/Hao1995/104hackathon/glob"
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

	res := &models.SearchJobsRes{}

	// Init Params
	if err := c.init(); err != nil {
		logs.Error(err)
		res.Error = err
		c.httpLib.WriteJSON(res)
		return
	}

	// Get The Result
	res, err := c.getResult()
	if err != nil {
		logs.Error(err)
		res.Error = err
		c.httpLib.WriteJSON(res)
		return
	}

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
		if err := glob.DB.QueryRow("SELECT `id` FROM `users` WHERE `email` = 'default@gmail.com'").Scan(&tmp); err != nil {
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

func (c *SearchJobsController) getResult() (*models.SearchJobsRes, error) {

	// - Find The Jobs That Corresponding The Key
	now := time.Now()
	jobMap, err := c.findRelativeJobs()
	if err != nil {
		return nil, err
	}
	logs.Info("findRelativeJobs spend %v", time.Since(now))

	// - Find Length Of `jobs`
	now = time.Now()
	countJobs, err := c.countAllJobs()
	if err != nil {
		return nil, err
	}
	logs.Info("countAllJobs spend %v", time.Since(now))

	// Get the Jobs Score
	now = time.Now()
	jobsScore, avgGoodScore, avgBadScore, err := c.getJobsScoreAndAvgScore(jobMap)
	if err != nil {
		return nil, err
	}
	logs.Info("getJobsScoreAndAvgScore = %v", time.Since(now))

	// - Get the Number of Different Scores for Good and Bad, Respectively.
	now = time.Now()
	goodScores, badScores, err := c.getGoodBadScores()
	if err != nil {
		return nil, err
	}
	logs.Info("getGoodBadScores spend %v", time.Since(now))

	// - Calculate AVG PR
	now = time.Now()
	goodPR, badPR := c.getAvgPR(avgGoodScore, avgBadScore, goodScores, badScores, countJobs)
	logs.Info("getPRofArea spend %v", time.Since(now))

	// - Calculate Jobs PR
	now = time.Now()
	jobsPR, err := c.getJobsPR(jobsScore, goodScores, badScores, countJobs)
	if err != nil {
		return nil, err
	}
	logs.Info("getPRofJobs = %v", time.Since(now))

	return &models.SearchJobsRes{
		Country: &models.SearchJobsScoreItem{
			GoodScore: &goodPR,
			BadScore:  &badPR},
		JobList: &jobsPR}, nil
}

func (c *SearchJobsController) findRelativeJobs() (map[int64]bool, error) {

	jobMap := make(map[int64]bool)

	rows, err := glob.DB.Query("SELECT `jobno`, `joblist` FROM `train_click` WHERE `key` = ?", c.params.Key)
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

	countJobs, err := cache.GetJobsInstance().CountJobs()
	if err != nil {
		return countJobs, err
	}

	return countJobs, nil
}

func (c *SearchJobsController) getJobsScoreAndAvgScore(jobMap map[int64]bool) ([]models.JobUserScoreGetItem, float64, float64, error) {

	jobsScore := []models.JobUserScoreGetItem{}
	var goodScore, badScore float64

	stmt, err := glob.DB.Prepare("SELECT `JUS`.`jobno`, `JUS`.`good_score`, `JUS`.`bad_score` FROM `jobs` AS `J` LEFT JOIN `job_user_score` AS `JUS` ON `J`.`jobno` = `JUS`.`jobno` AND `JUS`.`user_id` = ?, `districts` AS `D` WHERE 1 = 1 AND `J`.`addr_no` = `D`.`id` AND (IFNULL(?, 1) = 1 OR `D`.`id` LIKE ?) AND `J`.`jobno` IN (?" + strings.Repeat(",?", len(jobMap)-1) + ") ORDER BY `good_score` DESC")
	if err != nil {
		return nil, goodScore, badScore, err
	}
	defer stmt.Close()

	// logs.Info(SELECT `JUS`.`jobno`, `JUS`.`good_score`, `JUS`.`bad_score` FROM `jobs` AS `J` LEFT JOIN `job_user_score` AS `JUS` ON `J`.`jobno` = `JUS`.`jobno` AND `JUS`.`user_id` = ?, `districts` AS `D` WHERE 1 = 1 AND `J`.`addr_no` = `D`.`id` AND (IFNULL(?, 1) = 1 OR `D`.`id` LIKE ?) AND `J`.`jobno` IN (?" + strings.Repeat(",?", len(jobMap)-1) + ") ORDER BY `good_score` DESC")

	vals := []interface{}{c.params.UserID, c.params.AddrNo, c.params.AddrNo}
	for key := range jobMap {
		vals = append(vals, key)
	}
	// logs.Info(vals)
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

		goodScore += float64(*jobScore.GoodScore)
		badScore += float64(*jobScore.BadScore)
	}
	logs.Debug("Get %v jobs score.", len(jobsScore))

	// Calculate Average Score
	var avgGoodScore, avgBadScore float64
	if tmpLen := len(jobsScore); tmpLen != 0 {
		avgGoodScore = float64(goodScore) / float64(tmpLen)
		avgBadScore = float64(badScore) / float64(tmpLen)
	}
	logs.Debug("Area average good-score = %v, bad-score = %v", avgGoodScore, avgBadScore)

	return jobsScore, avgGoodScore, avgBadScore, nil
}

func (c *SearchJobsController) getGoodBadScores() (map[int]int, map[int]int, error) {

	goodScores, err := c.getGoodScores()
	if err != nil {
		return nil, nil, err
	}

	badScores, err := c.getBadScores()
	if err != nil {
		return nil, nil, err
	}

	return goodScores, badScores, nil
}

func (c *SearchJobsController) getAvgPR(avgGoodScore, avgBadScore float64, goodScores, badScores map[int]int, countJobs int64) (float64, float64) {

	// - Calculate AVG PR
	var goodCount, badCount int
	for score, count := range goodScores {
		if float64(score) < avgGoodScore {
			goodCount += count
		}
	}
	goodPR := float64(goodCount) / float64(countJobs)
	for score, count := range badScores {
		if float64(score) > avgBadScore {
			badCount += count
		}
	}
	badPR := float64(badCount) / float64(countJobs)

	logs.Debug("Area Good PR = %v, Bad PR = %v", goodPR, badPR)
	return goodPR, badPR
}

func (c *SearchJobsController) getJobsPR(jobsScore []models.JobUserScoreGetItem, goodScores, badScores map[int]int, countJobs int64) ([]models.SearchJobsListItem, error) {

	jobListScore := []models.SearchJobsListItem{}
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

		if err := glob.DB.QueryRow("SELECT `J`.`job` AS `job_name`, `C`.`name` AS `cust_name` FROM `jobs` AS `J`, `companies` AS `C` WHERE 1 = 1 AND `J`.`custno` = `C`.`custno` AND `J`.`jobno` = ?", item.JobNo).Scan(&jobScore.JobName, &jobScore.JobCompany); err != nil {
			return nil, err
		}

		jobListScore = append(jobListScore, jobScore)
	}

	return jobListScore, nil
}

func (c *SearchJobsController) getGoodScores() (map[int]int, error) {

	// Good
	goodScores := make(map[int]int)
	rows, err := glob.DB.Query("SELECT `good_score`, COUNT(1) AS `count` FROM `job_user_score` WHERE `user_id` = ? GROUP BY `good_score`", c.params.UserID)
	if err != nil {
		return goodScores, err
	}
	defer rows.Close()

	for rows.Next() {
		var score, count int
		if err := rows.Scan(&score, &count); err != nil {
			return goodScores, err
		}
		goodScores[score] = count
	}

	return goodScores, nil
}

func (c *SearchJobsController) getBadScores() (map[int]int, error) {

	// Bad
	badScores := make(map[int]int)
	rows, err := glob.DB.Query("SELECT `bad_score`, COUNT(1) AS `count` from `job_user_score` WHERE `user_id` = ? GROUP BY `bad_score`", c.params.UserID)
	if err != nil {
		return badScores, err
	}
	defer rows.Close()

	for rows.Next() {
		var score, count int
		if err := rows.Scan(&score, &count); err != nil {
			return badScores, err
		}
		badScores[score] = count
	}

	return badScores, nil
}
