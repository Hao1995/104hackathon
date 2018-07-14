package implement

import (
	"database/sql"
	"docker-example/log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Hao1995/Docker-Example/model"
)

var (
	tagScore      []*model.Tag
	areaScore     map[string]*model.AreaScore
	jobScore      map[string]*model.JobScore
	queryKeyScore map[string]*model.QueryKey
)

func init() {
	areaScore = make(map[string]*model.AreaScore)
	jobScore = make(map[string]*model.JobScore)
	queryKeyScore = make(map[string]*model.QueryKey)
}

//Score ...
func Score(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "test")
}

//ScoreArea ...
func ScoreArea(res http.ResponseWriter, req *http.Request) {
	//=====Params
	req.ParseForm()
	params := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "key":
			params[k] = strings.Join(v, "")
		}
	}

	var rows *sql.Rows
	var err error
	// if v, ok := params["size"]; ok {
	// 	rows, err = db.Query("SELECT * FROM job LIMIT " + v.(string))
	// } else {

	area := ""
	key := ""
	rows, err = db.Query("select `area`,`key`,`good_score`,`bad_score` from `area_key_score` where `area` = ? and `key` =?", area, key)
	// }

	items := []*model.AreaScore{}

	for rows.Next() {
		r := &model.AreaScore{}

		err = rows.Scan(&r.Area, &r.Key, &r.GoodScore, &r.BadScore)
		chechkErr(err)
		items = append(items, r)
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}

//ScoreJob ...
func ScoreJob(res http.ResponseWriter, req *http.Request) {
	//=====Params
	req.ParseForm()
	params := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "key":
			params[k] = strings.Join(v, "")
		}
	}

	var rows *sql.Rows
	var err error
	// if v, ok := params["size"]; ok {
	// 	rows, err = db.Query("SELECT * FROM job LIMIT " + v.(string))
	// } else {

	area := ""
	key := ""
	rows, err = db.Query("select `job`,`key`,`good_score`,`bad_score` from `area_key_score` where `area` = ? and `key` =?", area, key)
	// }

	items := []*model.JobScore{}

	for rows.Next() {
		r := &model.JobScore{}

		err = rows.Scan(&r.Job, &r.Key, &r.GoodScore, &r.BadScore)
		chechkErr(err)
		items = append(items, r)
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}

//SyncJobKey ...
func SyncJobKey(res http.ResponseWriter, req *http.Request) {

	total := 34891
	size := 1000
	offest := 0

	// var wg sync.WaitGroup
	// var mu sync.Mutex

	for {
		// go func() {
		// defer wg.Done()
		// mu.Lock()
		rows, err := db.Query("SELECT `train_click`.`key`, `job`.`job` FROM `train_click`, `job` WHERE 1 = 1 AND `train_click`.`jobno` =`job`.`jobno` AND `train_click`.`key` IS NOT NULL ORDER BY `train_click`.`key` LIMIT ? OFFSET ?", size, offest)

		queryString := "INSERT INTO job_key(`key`, `job`) VALUES"

		for rows.Next() {
			r := &model.JobKey{}

			err := rows.Scan(&r.Key, &r.Job)
			if err != nil {
				log.Errorf(err.Error())
			}
			value := "(" + stringAddSingleQuotation(processQuote(r.Key)) + "," + stringAddSingleQuotation(processQuote(r.Job)) + "),"
			queryString = queryString + value
		}

		queryString = strings.TrimRight(queryString, ",")

		fmt.Println(queryString)
		stmt, err := db.Prepare(queryString)
		defer stmt.Close()

		_, err = stmt.Exec()
		if err != nil {
			log.Errorf(err.Error())
		}

		// mu.Unlock()
		// }()

		// wg.Wait()
		offest = offest + size
		if offest > total {
			break
		}
	}
}

func CalKeyScore(res http.ResponseWriter, req *http.Request) {

	fmt.Println("===== Get All Tag")
	rows, err := db.Query("SELECT `id`,`name`,`score` FROM tag;")

	tagScore = []*model.Tag{}
	for rows.Next() {
		r := &model.Tag{}

		err := rows.Scan(&r.ID, &r.Name, &r.Score)
		if err != nil {
			log.Errorf(err.Error())
		}

		tagScore = append(tagScore, r)
	}

	fmt.Println("===== Get All Key")
	rows, err = db.Query("SELECT `name` FROM query_key;")

	queryKeys := []*model.QueryKey{}
	for rows.Next() {
		r := &model.QueryKey{}

		err := rows.Scan(&r.Name)
		if err != nil {
			log.Errorf(err.Error())
		}

		wg.Add(1)
		go CalKeyScoreGetOriginInfoOfKey(r)

		queryKeys = append(queryKeys, r)
	}
	wg.Wait()

	fmt.Println("===== Insert Data")
	queryString := "INSERT INTO job_key(`key`, `job`) VALUES"

	for rows.Next() {
		r := &model.JobKey{}

		err := rows.Scan(&r.Key, &r.Job)
		if err != nil {
			log.Errorf(err.Error())
		}
		value := "(" + stringAddSingleQuotation(processQuote(r.Key)) + "," + stringAddSingleQuotation(processQuote(r.Job)) + "),"
		queryString = queryString + value
	}

	queryString = strings.TrimRight(queryString, ",")

	fmt.Println(queryString)
	stmt, err := db.Prepare(queryString)
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Errorf(err.Error())
	}
}

//CalKeyScoreGetOriginInfoOfKey ...
func CalKeyScoreGetOriginInfoOfKey(r *model.QueryKey) {

	key := r.Name

	defer wg.Done()

	mu.Lock()
	fmt.Println("===== Get Info of the Key : ", key)
	rows, err := db.Query("SELECT  e.`action` AS 'job_action',`e`.`key`,`e`.`job`,`e`.welfare AS 'company_walfare',`f`.id AS 'districk_id',`f`.name AS 'districk_name' FROM `district` AS f RIGHT JOIN (SELECT d.`key`,`c`.name,`c`.profile,`c`.welfare,`d`.`addr_no`,`d`.`job`,`d`.`action` FROM `companies` AS c RIGHT JOIN(SELECT a.`key`, custno, `b`.addr_no, `b`.`job`, `a`.`action` FROM `job` AS b RIGHT JOIN (SELECT  `train_click`.key, jobno, `action` FROM `train_click` WHERE `train_click`.key = ? AND `train_click`.`action` IN ('clickApply' , 'clickSave')) AS a ON b.jobno = a.jobno) AS d ON c.custno = d.custno) AS e ON e.addr_no = f.id", key)

	if err != nil {
		log.Errorf(err.Error())
	}

	for rows.Next() {
		r := &model.ScoreOriginData{}

		err := rows.Scan(&r.JobAction, &r.Key, &r.JobName, &r.CompanyWelfare, &r.DistrictID, &r.DistrictName)
		if err != nil {
			log.Errorf(err.Error())
		}
		goodScore, badScore := CalScore(r.CompanyWelfare)

		if r.JobAction == "clickApply" {
			goodScore *= 3
			badScore *= 3
		} else if r.JobAction == "clickSave" {
			goodScore *= 2
			badScore *= 2
		}

		if _, ok := areaScore[r.DistrictName]; !ok {
			areaScore[r.DistrictName] = &model.AreaScore{}
		}
		areaScore[r.DistrictName].GoodScore = areaScore[r.DistrictName].GoodScore + goodScore
		areaScore[r.DistrictName].BadScore = areaScore[r.DistrictName].BadScore + badScore

		if _, ok := jobScore[r.JobName]; !ok {
			jobScore[r.JobName] = &model.JobScore{}
		}
		jobScore[r.JobName].GoodScore = jobScore[r.JobName].GoodScore + goodScore
		jobScore[r.JobName].BadScore = jobScore[r.JobName].BadScore + badScore

		if _, ok := queryKeyScore[r.Key]; !ok {
			queryKeyScore[r.Key] = &model.QueryKey{}
		}
		queryKeyScore[r.Key].GoodScore = queryKeyScore[r.Key].GoodScore + goodScore
		queryKeyScore[r.Key].BadScore = queryKeyScore[r.Key].BadScore + badScore
	}

	mu.Unlock()
}

func CalScore(welfare string) (int, int) {

	goodScore := 0
	badScore := 0
	for _, v := range tagScore {
		if strings.Contains(welfare, v.Name) {
			tagScore := v.Score
			if tagScore >= 0 {
				goodScore = goodScore + tagScore
			} else {
				badScore = badScore + tagScore
			}
		}
	}

	return goodScore, badScore
}
