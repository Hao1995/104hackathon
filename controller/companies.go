package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"sync"

	"github.com/Hao1995/104hackathon/config"
	"github.com/Hao1995/104hackathon/glob"
	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

// SyncCompanies :
// Sync the company data to DB
// * Delete all DB data
// * Insert JSON data to DB
func SyncCompanies(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	res := models.APIRes{}

	// - Get the size of companies data
	var companiesIdx int
	rows, err := glob.DB.Query("SELECT COUNT(1) FROM `companies`")
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
		err := rows.Scan(&companiesIdx)
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
	logs.Debug("Will skip %v rows.", companiesIdx)

	// - Open Data File
	file, err := os.Open(config.CfgData.Data.Companies)
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
	v := reflect.ValueOf(models.CompaniesJSONItem{})
	size := glob.MySQLUpperPlaceholders / v.NumField()
	skipCount := 0
	totalCount := 0
	guard := make(chan struct{}, 2) // Max goroutines limit.
	errChan := make(chan bool)
	items := []models.CompaniesJSONItem{}
	for scanner.Scan() {
		// - Skip existing data.
		if skipCount < companiesIdx {
			skipCount++
			continue
		}

		// - Parse JSON to Item
		itemJSON := scanner.Text()
		item := models.CompaniesJSONItem{}
		json.Unmarshal([]byte(itemJSON), &item)
		// logs.Debug("JSON = %+v", item)

		items = append(items, item)
		totalCount++
		// break
		if len(items) >= size {
			// - Send items to channel and clear items, skipCount.
			guard <- struct{}{}
			wg.Add(1)
			go syncCompaniesInsertData(&wg, guard, errChan, items)
			logs.Debug("Send to insert %v data.", size)
			items = []models.CompaniesJSONItem{}
		}
	}
	if len(items) > 0 {
		// - Send the last data that not reach the size.
		guard <- struct{}{}
		wg.Add(1)
		go syncCompaniesInsertData(&wg, guard, errChan, items)
		logs.Debug("Last time send to insert %v data.", len(items))
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

	logs.Debug("Waiting all goroutine")
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

func syncCompaniesInsertData(wg *sync.WaitGroup, guard chan struct{}, errChan chan bool, data []models.CompaniesJSONItem) {
	defer wg.Done()

	sqlStr := "INSERT INTO `companies` (`custno`, `invoice`, `name`, `profile`, `management`, `welfare`, `product`) VALUES "
	vals := []interface{}{}
	for _, item := range data {
		custno := utils.NewNullString(item.Custno)
		invoice := utils.NewNullInt64(item.Invoice)
		name := utils.NewNullString(item.Name)
		profile := utils.NewNullString(item.Profile)
		management := utils.NewNullString(item.Management)
		welfare := utils.NewNullString(item.Welfare)
		product := utils.NewNullString(item.Product)
		sqlStr += "(?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, custno, invoice, name, profile, management, welfare, product)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// logs.Debug(sqlStr)
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
