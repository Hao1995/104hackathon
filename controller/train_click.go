package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/Hao1995/104hackathon/config"
	"github.com/Hao1995/104hackathon/glob"
	"github.com/Hao1995/104hackathon/models"
	"github.com/astaxie/beego/logs"
)

// //Get ...
// func Get(res http.ResponseWriter, req *http.Request) {

// 	// - Params
// 	req.ParseForm()
// 	params := make(map[string]interface{})
// 	for k, v := range req.Form {
// 		switch k {
// 		case "size":
// 			params[k] = strings.Join(v, "")
// 			// case "message":
// 			// 	params[k] = strings.Join(v, "")
// 		}
// 	}

// 	var rows *sql.Rows
// 	var err error
// 	if v, ok := params["size"]; ok {
// 		rows, err = db.Query("SELECT * FROM companies LIMIT " + v.(string))
// 	} else {
// 		rows, err = db.Query("SELECT * FROM companies LIMIT 100")
// 	}

// 	companies := []*models.Company{}

// 	for rows.Next() {
// 		r := &models.Company{}

// 		err = rows.Scan(&r.Custno, &r.Invoice, &r.Name, &r.Profile, &r.Management, &r.Welfare, &r.Product)
// 		chechkErr(err)
// 		companies = append(companies, r)
// 	}

// 	jsonData, err := json.Marshal(companies)
// 	if err != nil {
// 		chechkErr(err)
// 	}
// 	io.WriteString(res, string(jsonData))
// }

//QueryKey ...
func QueryKey(res http.ResponseWriter, req *http.Request) {

	//=====Get Total
	fmt.Println("===== Get Total")
	rows, err := db.Query("SELECT COUNT(1) FROM `train_click`")
	if err != nil {
		logs.Error(err.Error())
	}

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			logs.Error(err.Error())
		}
		// fmt.Printf("%v \n", count)
	}
	if err := rows.Err(); err != nil {
		logs.Error(err.Error())
	}

	//=====Get OriginQueryString
	fmt.Println("===== Get OriginQueryString")
	offset := 387092
	size := 10000
	type OriginQueryString struct {
		ID          int    `json:"id"`
		QueryString string `json:"querystring"`
	}

	for {
		query :=
			"SELECT `id`, `querystring` FROM `train_click` " +
				"ORDER BY `id` " +
				"LIMIT " + strconv.Itoa(size) + " " +
				"OFFSET " + strconv.Itoa(offset)

		rows, err = db.Query(query)
		if err != nil {
			logs.Error(err.Error())
		}

		originDatas := []OriginQueryString{}
		for rows.Next() {

			var id int
			var queryString string
			if err := rows.Scan(&id, &queryString); err != nil {
				logs.Error(err.Error())
				continue
			}

			originData := OriginQueryString{
				ID:          id,
				QueryString: queryString,
			}

			originDatas = append(originDatas, originData)
		}
		if err := rows.Err(); err != nil {
			logs.Error(err.Error())
		}

		//===== Decode
		fmt.Println("===== Decode")
		decodeKey := make(map[int]string)
		for _, v := range originDatas {
			if v.QueryString == "" {
				continue
			}
			str := "localhost/?" + v.QueryString
			u, err := url.Parse(str)
			if err != nil {
				logs.Error(err.Error())
			}
			// fmt.Println(u.String())

			m, _ := url.ParseQuery(u.RawQuery)
			if key, ok := m["keyword"]; ok {
				// fmt.Println(key[0])
				decodeKey[v.ID] = key[0]
				continue
			} else {
				logs.Error("'keyword' does not exist.")
			}
		}

		//===== Insert
		fmt.Println("===== Insert")
		if len(decodeKey) > 0 {
			for k, v := range decodeKey {
				stmt, err := db.Prepare("UPDATE `train_click` SET `key`=? WHERE `id`= ?;")
				if err != nil {
					logs.Error("[db.Prepare] " + err.Error())
				}
				_, err = stmt.Exec(v, k)
				if err != nil {
					logs.Error("[stmt.Exec] " + err.Error())
				}
				stmt.Close()
			}
		}

		offset = offset + size
		if offset > count {
			break
		}
	}
	//===== Complete
	io.WriteString(res, "Complete")
}

// SyncTrainClick :
// Sync the train-click data to DB
// * Delete all DB data
// * Insert JSON data to DB
func SyncTrainClick(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := models.APIRes{}

	// - Get the size of train_click data
	var trainClickIdx int
	rows, err := db.Query("SELECT COUNT(1) FROM `train_click`")
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
		err := rows.Scan(&trainClickIdx)
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
	logs.Trace("Will skip %v rows.", trainClickIdx)

	// - Open Data File
	file, err := os.Open(config.CfgData.Data.Train_Click)
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
	v := reflect.ValueOf(models.TrainClickJSONItem{})
	size := glob.MySQLUpperPlaceholders / v.NumField()
	skipCount := 0
	totalCount := 0
	guard := make(chan struct{}, 2)
	errChan := make(chan bool)
	items := []models.TrainClickJSONItem{}
	for scanner.Scan() {
		// - Skip existing data.
		if skipCount < trainClickIdx {
			skipCount++
			continue
		}

		// - Parse JSON to Item
		itemJSON := scanner.Text()
		item := models.TrainClickJSONItem{}
		json.Unmarshal([]byte(itemJSON), &item)
		// logs.Debug("JSON = %v", item)

		items = append(items, item)
		if len(items) >= size {
			// - Send items to channel and clear items, count.
			wg.Add(1)
			guard <- struct{}{}
			go syncTrainClickInsertData(&wg, guard, errChan, items)
			logs.Trace("Send to insert %v data.", size)
			items = []models.TrainClickJSONItem{}
		}
	}
	if len(items) > 0 {
		// - Send the last data that not reach the size.
		guard <- struct{}{}
		wg.Add(1)
		go syncTrainClickInsertData(&wg, guard, errChan, items)
		logs.Trace("Last time send to insert %v data.", len(items))
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

	logs.Trace("Waiting all goroutine")
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

func syncTrainClickInsertData(wg *sync.WaitGroup, guard chan struct{}, errChan chan bool, data []models.TrainClickJSONItem) {
	defer wg.Done()

	sqlStr := "INSERT INTO `train_click` (`action`, `jobno`, `date`, `joblist`, `querystring`, `source`) VALUES "
	vals := []interface{}{}
	for _, item := range data {
		action := item.Action
		jobno, err := strconv.ParseInt(item.Jobno, 10, 64)
		if err != nil {
			select {
			case errChan <- true:
			default:
				<-guard
				return
			}
		}
		tmpDate, err := strconv.ParseInt(item.Date, 10, 64)
		if err != nil {
			select {
			case errChan <- true:
			default:
				<-guard
				return
			}
		}
		tmpDate /= 1000 // without millseconds. ex 1527330445000 -> 1527330445
		date := time.Unix(tmpDate, 0)
		// fmt.Print(date)
		jobList := ""
		for _, job := range item.Joblist {
			// Check 'job' is integer.
			_, err := strconv.ParseInt(job, 10, 64)
			if err != nil {
				select {
				case errChan <- true:
				default:
					return
				}
			}
			// Store as string
			jobList += job + ","
		}
		jobList = jobList[0 : len(jobList)-1]
		queryString := item.QueryString
		source := item.Source
		sqlStr += "(?, ?, ?, ?, ?, ?),"
		vals = append(vals, action, jobno, date, jobList, queryString, source)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// logs.Trace(sqlStr)
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
	// logs.Trace(vals)
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

//ParseTrainClick ...
func ParseTrainClick(fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// c := []*models.Job{}
	c := []*models.TrainClick{}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		fmt.Println(err.Error())
		FailFile = append(FailFile, fileName)
		return
	}

	for _, v := range c {
		// InsertToJob(fileName, v) //job
		trainClcikInsert(fileName, v) //companies
	}
}

//TrainClcikInsert ...
func trainClcikInsert(fileName string, v *models.TrainClick) {
	mu.Lock()
	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}
	stmt, err := db.Prepare("INSERT INTO train_click(`action`, `jobno`, `date`, `joblist`, `querystring`, `source`) VALUES(?,?,?,?,?,?)")
	defer stmt.Close()
	dbConnentCount++

	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}
	chechkErr(err)

	jobList := "["
	for _, job := range v.Joblist {
		jobList = jobList + job + ","
	}
	jobListByte := []byte(jobList)
	jobListByte = jobListByte[:len(jobList)-1]

	jobListFinal := string(jobListByte) + "]"
	// fmt.Println("[jobListByte] ", jobListFinal)

	_, err = stmt.Exec(v.Action, v.Jobno, v.Date, jobListFinal, v.QueryString, v.Source)
	dbConnentCount++
	mu.Unlock()
	if err != nil {
		fmt.Printf("[ERROR][%v][%v] Content :%v \n", fileName, err, *v)
	}

	dbConnentCount--
	dbConnentCount--
}
