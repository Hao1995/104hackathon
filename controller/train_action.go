package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
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

// SyncTrainAction :
// Sync the train-action data to DB
// * Delete all DB data
// * Insert JSON data to DB
func SyncTrainAction(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := models.APIRes{}

	// - Get the size of train_action data
	var trainActionIdx int
	rows, err := db.Query("SELECT COUNT(1) FROM `train_action`")
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
		err := rows.Scan(&trainActionIdx)
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
	logs.Trace("Will skip %v rows.", trainActionIdx)

	// - Open Data File
	file, err := os.Open(config.CfgData.Data.Train_Action)
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
	v := reflect.ValueOf(models.TrainActionJSONItem{})
	size := glob.MySQLUpperPlaceholders / v.NumField()
	skipCount := 0
	totalCount := 0
	guard := make(chan struct{}, 2) // Max goroutines limit.
	errChan := make(chan bool)
	items := []models.TrainActionJSONItem{}
	for scanner.Scan() {
		// - Skip existing data.
		if skipCount < trainActionIdx {
			skipCount++
			continue
		}

		// - Parse JSON to Item
		itemJSON := scanner.Text()
		item := models.TrainActionJSONItem{}
		json.Unmarshal([]byte(itemJSON), &item)
		// logs.Trace("JSON = %v", item)

		items = append(items, item)
		totalCount++
		if len(items) >= size {
			// - Send items to channel and clear items, skipCount.
			guard <- struct{}{}
			wg.Add(1)
			go syncTrainActionInsertData(&wg, guard, errChan, items)
			logs.Trace("Send to insert %v data.", size)
			items = []models.TrainActionJSONItem{}
		}
	}
	if len(items) > 0 {
		// - Send the last data that not reach the size.
		guard <- struct{}{}
		wg.Add(1)
		go syncTrainActionInsertData(&wg, guard, errChan, items)
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

func syncTrainActionInsertData(wg *sync.WaitGroup, guard chan struct{}, errChan chan bool, data []models.TrainActionJSONItem) {
	defer wg.Done()

	sqlStr := "INSERT INTO `train_action` (`jobno`, `date`, `action`, `source`, `device`) VALUES "
	vals := []interface{}{}
	for _, item := range data {
		action := item.Action
		jobno, err := strconv.ParseInt(item.Jobno, 10, 64)
		if err != nil {
			logs.Error(err)
			select {
			case errChan <- true:
			default:
				<-guard
				return
			}
		}
		tmpDate, err := strconv.ParseInt(item.Date, 10, 64)
		if err != nil {
			logs.Error(err)
			select {
			case errChan <- true:
			default:
				<-guard
				return
			}
		}
		tmpDate /= 1000 // without millseconds. ex 1527330445000 -> 1527330445
		date := time.Unix(tmpDate, 0)
		source := item.Source
		device := item.Device
		sqlStr += "(?, ?, ?, ?, ?),"
		vals = append(vals, jobno, date, action, source, device)
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
