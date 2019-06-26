package controller

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

	// - Begin Transaction
	tx, err := db.Begin()
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

	// - Delete DB data
	if _, err := tx.Exec("DELETE FROM `train_action`"); err != nil {
		tx.Rollback()
		log.Fatal(err)
		res.Error = err.Error()

		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	}

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
	count := 0
	errChan := make(chan bool)
	items := []models.TrainActionJSONItem{}
	for scanner.Scan() {
		// - Parse JSON to Item
		itemJSON := scanner.Text()
		item := models.TrainActionJSONItem{}
		json.Unmarshal([]byte(itemJSON), &item)
		// logs.Debug("JSON = %v", item)

		items = append(items, item)
		if len(items) >= size {
			// - Send items to channel and clear items, count.
			wg.Add(1)
			go syncTrainActionInsertData(&wg, errChan, tx, items)
			items = []models.TrainActionJSONItem{}
			count += size
		}
	}
	if len(items) > 0 {
		// - Send the last data that not reach the size.
		wg.Add(1)
		go syncTrainActionInsertData(&wg, errChan, tx, items)
		count += len(items)
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

	wg.Wait()

	err = tx.Commit()
	if err != nil {
		res.Error = err.Error()
		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
		return
	}

	// - Rev Err
	select {
	case <-errChan:
		res.Error = fmt.Sprintf("Something wrong !")
	default:
		res.Message = fmt.Sprintf("Success insert %v data", count)
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)

	return
}

func syncTrainActionInsertData(wg *sync.WaitGroup, errChan chan bool, tx *sql.Tx, data []models.TrainActionJSONItem) {
	defer wg.Done()

	sqlStr := "INSERT INTO `train_action` (`jobno`, `date`, `action`, `source`, `device`) VALUES "
	vals := []interface{}{}
	for _, item := range data {
		action := item.Action
		jobno, err := strconv.ParseInt(item.Jobno, 10, 64)
		if err != nil {
			select {
			case errChan <- true:
			default:
				return
			}
		}
		tmpDate, err := strconv.ParseInt(item.Date, 10, 64)
		if err != nil {
			select {
			case errChan <- true:
			default:
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
	logs.Debug(sqlStr)
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		logs.Error(err)
		select {
		case errChan <- true:
		default:
			return
		}
	}
	logs.Debug(vals)
	_, err = stmt.Exec(vals...)
	if err != nil {
		tx.Rollback()
		logs.Error(err)
		select {
		case errChan <- true:
		default:
			return
		}
	}

	return
}
