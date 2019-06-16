package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hao1995/104hackathon/log"
	"github.com/Hao1995/104hackathon/model"
)

var (
	tagScore      []*model.Tag
	areaScore     map[string]*model.AreaScore
	jobScore      map[string]*model.JobScore
	queryKeyScore map[string]*model.QueryKey

	areaMappingId map[string]string

	PRMapping map[int]int
)

func init() {
	areaScore = make(map[string]*model.AreaScore)
	jobScore = make(map[string]*model.JobScore)
	queryKeyScore = make(map[string]*model.QueryKey)

	areaMappingId = make(map[string]string)

	areaMappingId["6001001"] = "台北市"
	areaMappingId["6001002"] = "新北市"
	areaMappingId["6001003"] = "宜蘭縣"
	areaMappingId["6001004"] = "基隆市"
	areaMappingId["6001005"] = "桃園市"
	areaMappingId["6001006"] = "新竹縣市"
	areaMappingId["6001007"] = "苗栗縣"
	areaMappingId["6001008"] = "台中市"
	areaMappingId["6001009"] = "台中市(原台中縣)"
	areaMappingId["6001010"] = "彰化縣"
	areaMappingId["6001011"] = "南投縣"
	areaMappingId["6001012"] = "雲林縣"
	areaMappingId["6001013"] = "嘉義縣市"
	areaMappingId["6001014"] = "台南市"
	areaMappingId["6001015"] = "台南市(原台南縣)"
	areaMappingId["6001016"] = "高雄市"
	areaMappingId["6001017"] = "高雄市(原高雄縣)"
	areaMappingId["6001018"] = "屏東縣"
	areaMappingId["6001019"] = "台東縣"
	areaMappingId["6001020"] = "花蓮縣"

	PRMapping = make(map[int]int)

}

// ScoreArea :
// The score of this job in a certain area.
func ScoreArea(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Content-Type", "application/json")

	start := time.Now()

	finalReturn := &model.FinalReturn{}
	finalReturnCountry := &model.FinalReturnCountry{}
	finalReturnJobList := []*model.FinalReturnJobList{}

	// - Params
	fmt.Println("=== Parse Parameters")
	req.ParseForm()
	params := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "key":
			params[k] = strings.Join(v, "")
		case "country_id":
			if _, ok := areaMappingId[v[0]]; ok {
				params[k] = strings.Join(v, "")
			}
		case "size":
			params[k] = strings.Join(v, "")
		case "page":
			params[k] = strings.Join(v, "")
		}
	}

	key, ok := params["key"]
	if !ok {
		io.WriteString(res, "Parameter 'key' is necessary.")
		return
	}

	countryID, ok := params["country_id"]
	if !ok {
		io.WriteString(res, "Parameter 'key' is necessary.")
		return
	}

	size, ok := params["size"]
	if !ok {
		io.WriteString(res, "Parameter 'key' is necessary.")
		return
	}

	page, ok := params["page"]
	if !ok {
		io.WriteString(res, "Parameter 'key' is necessary.")
		return
	}

	//=== Original Data From `104hackathon-welfare`.`area_job_key_score`
	fmt.Println("=== Original Data From `104hackathon-welfare`.`area_job_key_score`")
	var rows *sql.Rows
	var err error

	countryIDStr := countryID.(string) + "%"
	sizeInt, err := strconv.Atoi(size.(string))
	if err != nil {
		log.Errorf(err.Error())
	}
	if sizeInt < 0 {
		io.WriteString(res, "Parameter [size] can not be negative number.")
	}
	pageInt, err := strconv.Atoi(page.(string))
	if err != nil {
		log.Errorf(err.Error())
	}
	if pageInt < 0 {
		io.WriteString(res, "Parameter [page] can not be negative number.")
	}
	offset := (pageInt - 1) * sizeInt

	rows, err = db.Query("SELECT `job`, `good_score`, `bad_score` FROM `104hackathon-welfare`.`area_job_key_score` WHERE `addr_no` like ? AND `key` = ? GROUP BY `addr_no`,`jobno` LIMIT ? OFFSET ? ", countryIDStr, key, size, offset)

	for rows.Next() {
		r := &model.FinalReturnJobList{}
		err = rows.Scan(&r.JobName, &r.GoodScore, &r.BadScore)
		if err != nil {
			log.Errorf(err.Error())
		}
		r.JobCompany = ""
		finalReturnJobList = append(finalReturnJobList, r)
	}
	fmt.Printf("%s took %v\n", "Load data from `area_job_key_socre`", time.Since(start))

	//=== Average Data Of The Area
	fmt.Println("=== Average Data Of The Area")
	start = time.Now()

	rows, err = db.Query("SELECT AVG(`good_score`) AS `good_score`, AVG(`bad_score`) AS `bad_score` FROM ( SELECT `good_score`, `bad_score` FROM `104hackathon-welfare`.`area_job_key_score` WHERE `addr_no` LIKE ? AND `key` = ? GROUP BY `addr_no`,`jobno` ) AS `tmp`", countryIDStr, key)

	for rows.Next() {
		err = rows.Scan(&finalReturnCountry.GoodScore, &finalReturnCountry.BadScore)
		if err != nil {
			log.Errorf(err.Error())
		}
	}
	fmt.Printf("%s took %v\n", "Average Data Of The Area", time.Since(start))

	//=== Marshal Data to JSON
	fmt.Println("=== Marshal Data to JSON")
	start = time.Now()

	finalReturn.Country = finalReturnCountry
	finalReturn.JobList = finalReturnJobList

	jsonData, err := json.Marshal(finalReturn)
	if err != nil {
		log.Errorf(err.Error())
	}
	fmt.Printf("%s took %v\n", "Marshal Data to JSON", time.Since(start))
	io.WriteString(res, string(jsonData))
	return
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
