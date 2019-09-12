package controller

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
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
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

// SyncTrainClick :
// Sync the train-click data to DB
// * Delete all DB data
// * Insert JSON data to DB
func SyncTrainClick(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := models.APIRes{}

	// - Get the size of train_click data
	var trainClickIdx int
	rows, err := glob.DB.Query("SELECT COUNT(1) FROM `train_click`")
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
	stmt, err := glob.DB.Prepare(sqlStr)
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

// SyncTrainClickKey :
// Sync the key that within querystring to DB

func SyncTrainClickKey(w http.ResponseWriter, req *http.Request) {

	ins := &SyncTrainClickKeyController{}
	httpLib := &utils.HTTPLib{}
	httpLib.Init(w, req)

	res := models.APIRes{}

	switch req.Method {
	// case http.MethodGet:
	// ins.get(httpLib)
	case http.MethodPost:
		ins.post(httpLib)
	// case http.MethodDelete:
	// ins.delete(httpLib)
	default:
		res.Error = fmt.Sprintf("There is no way correspond to method '%v'", req.Method)
		httpLib.WriteJSON(res)
		return
	}
}

type SyncTrainClickKeyController struct {
}

func (c *SyncTrainClickKeyController) post(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	// - Sync Data
	// Get job info
	rows, err := glob.DB.Query("SELECT `id`, `querystring` FROM `train_click`")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer rows.Close()

	tx, err := glob.DB.Begin()
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	var notFindKeyword, notFindQueryString, findKeyword int
	for rows.Next() {
		var id *int
		var querystring *string
		err := rows.Scan(&id, &querystring)
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}

		// Get keyword
		var key sql.NullString
		if *querystring != "" {
			str := "localhost/?" + *querystring
			u, err := url.Parse(str)
			if err != nil {
				logs.Error(err)
				res.Error = err.Error()
				httpLib.WriteJSON(res)
				return
			}
			m, _ := url.ParseQuery(u.RawQuery)
			if _, ok := m["keyword"]; ok {
				val := m["keyword"][0]
				key = utils.NewNullString(&val)
				findKeyword++
			} else {
				notFindKeyword++
			}
		} else {
			notFindQueryString++
		}

		// Insert Key
		_, err = tx.Exec("UPDATE `train_click` SET `key` = ? WHERE `id` = ?", key, *id)
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
	}

	tx.Commit()

	res.Message = fmt.Sprintf("No querystring = %v. No keyword = %v. Success insert = %v", notFindQueryString, notFindKeyword, findKeyword)
	httpLib.WriteJSON(res)
	return
}
