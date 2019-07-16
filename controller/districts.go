package controller

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"

	"github.com/Hao1995/104hackathon/config"
	"github.com/Hao1995/104hackathon/glob"
	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

// SyncDistricts :
// Sync the district data to DB
// * Delete all DB data
// * Insert JSON data to DB
func SyncDistricts(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := models.APIRes{}

	// - Get the size of companies data
	var districtsIdx int
	rows, err := db.Query("SELECT COUNT(1) FROM `districts`")
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
		err := rows.Scan(&districtsIdx)
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
	logs.Debug("Will skip %v rows.", districtsIdx)

	// - Open Data File
	file, err := os.Open(config.CfgData.Data.Districts)
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
	reader := csv.NewReader(bufio.NewReader(file))
	v := reflect.ValueOf(models.DistrictsJSONItem{})
	size := glob.MySQLUpperPlaceholders / v.NumField()
	skipCount := 0
	totalCount := 0
	guard := make(chan struct{}, 2) // Max goroutines limit.
	errChan := make(chan bool)
	items := []models.DistrictsJSONItem{}
	for {

		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
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

		// - Skip first row and the existing data.
		if skipCount < districtsIdx+1 {
			skipCount++
			continue
		}

		id, err := strconv.ParseInt(line[0], 10, 64)
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

		// - Parse JSON to Item
		items = append(items, models.DistrictsJSONItem{
			ID:   &id,
			Name: &line[1],
			Desc: &line[2],
			Hide: &line[3],
		})
		// logs.Debug("JSON = %+v", item)

		totalCount++
		if len(items) >= size {
			// - Send items to channel and clear items, skipCount.
			guard <- struct{}{}
			wg.Add(1)
			go syncDistrictsInsertData(&wg, guard, errChan, items)
			logs.Debug("Send to insert %v data.", size)
			items = []models.DistrictsJSONItem{}
		}
	}
	if len(items) > 0 {
		// - Send the last data that not reach the size.
		guard <- struct{}{}
		wg.Add(1)
		go syncDistrictsInsertData(&wg, guard, errChan, items)
		logs.Debug("Last time send to insert %v data.", len(items))
	}

	logs.Debug("Waiting all goroutine")
	wg.Wait()

	// - Rev Err
	select {
	case <-errChan:
		res.Error = fmt.Sprintf("Something wrong !")
	default:
		res.Message = fmt.Sprintf("Skip %v data and insert %v data", skipCount-1, totalCount)
	}

	js, err := json.Marshal(res)
	if err != nil {
		logs.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)

	return
}

func syncDistrictsInsertData(wg *sync.WaitGroup, guard chan struct{}, errChan chan bool, data []models.DistrictsJSONItem) {
	defer wg.Done()

	sqlStr := "INSERT INTO `districts` (`id`, `name`, `desc`, `hide`) VALUES "
	vals := []interface{}{}
	for _, item := range data {
		id := utils.NewNullInt64(item.ID)
		name := utils.NewNullString(item.Name)
		desc := utils.NewNullString(item.Desc)
		hide := utils.NewNullString(item.Hide)

		sqlStr += "(?, ?, ?, ?),"
		vals = append(vals, id, name, desc, hide)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// logs.Debug(sqlStr)
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
	// logs.Debug(vals)
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
