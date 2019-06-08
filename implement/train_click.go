package implement

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Hao1995/104hackathon/model"

	"github.com/Hao1995/104hackathon/log"
)

//HackathonTrainClick ...
func HackathonTrainClick(res http.ResponseWriter, req *http.Request) {
	//=====Params
	req.ParseForm()
	params := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "size":
			params[k] = strings.Join(v, "")
			// case "message":
			// 	params[k] = strings.Join(v, "")
		}
	}

	var rows *sql.Rows
	var err error
	if v, ok := params["size"]; ok {
		rows, err = db.Query("SELECT * FROM companies LIMIT " + v.(string))
	} else {
		rows, err = db.Query("SELECT * FROM companies LIMIT 100")
	}

	companies := []*model.Company{}

	for rows.Next() {
		r := &model.Company{}

		err = rows.Scan(&r.Custno, &r.Invoice, &r.Name, &r.Profile, &r.Management, &r.Welfare, &r.Product)
		chechkErr(err)
		companies = append(companies, r)
	}

	jsonData, err := json.Marshal(companies)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}

//QueryKey ...
func QueryKey(res http.ResponseWriter, req *http.Request) {

	//=====Get Total
	fmt.Println("===== Get Total")
	rows, err := db.Query("SELECT COUNT(1) FROM `train_click`")
	if err != nil {
		log.Errorf(err.Error())
	}

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			log.Errorf(err.Error())
		}
		// fmt.Printf("%v \n", count)
	}
	if err := rows.Err(); err != nil {
		log.Errorf(err.Error())
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
			log.Errorf(err.Error())
		}

		originDatas := []OriginQueryString{}
		for rows.Next() {

			var id int
			var queryString string
			if err := rows.Scan(&id, &queryString); err != nil {
				log.Errorf(err.Error())
				continue
			}

			originData := OriginQueryString{
				ID:          id,
				QueryString: queryString,
			}

			originDatas = append(originDatas, originData)
		}
		if err := rows.Err(); err != nil {
			log.Errorf(err.Error())
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
				log.Errorf(err.Error())
			}
			// fmt.Println(u.String())

			m, _ := url.ParseQuery(u.RawQuery)
			if key, ok := m["keyword"]; ok {
				// fmt.Println(key[0])
				decodeKey[v.ID] = key[0]
				continue
			} else {
				log.Errorf("'keyword' does not exist.")
			}
		}

		//===== Insert
		fmt.Println("===== Insert")
		if len(decodeKey) > 0 {
			for k, v := range decodeKey {
				stmt, err := db.Prepare("UPDATE `train_click` SET `key`=? WHERE `id`= ?;")
				if err != nil {
					log.Errorf("[db.Prepare] " + err.Error())
				}
				_, err = stmt.Exec(v, k)
				if err != nil {
					log.Errorf("[stmt.Exec] " + err.Error())
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

//InsertTrainClick User
func InsertTrainClick(res http.ResponseWriter, req *http.Request) {
	directoryPath := "F:/gotool/src/test/test1/data/train_click" //train_click
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Errorf(err.Error())
	}
	for _, file := range files {
		fileExtension := strings.Split(file.Name(), ".")
		if len(fileExtension) == 2 {
			if fileExtension[1] == "json" {
				filePath := directoryPath + "/" + file.Name()
				go func(filePath string) {
					ParseTrainClick(filePath)
				}(filePath)
			}
		}

	}

	io.WriteString(res, "Complete")
}

//ParseTrainClick ...
func ParseTrainClick(fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// c := []*model.Job{}
	c := []*model.TrainClick{}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		fmt.Println(err.Error())
		FailFile = append(FailFile, fileName)
		return
	}

	for _, v := range c {
		// InsertToJob(fileName, v) //job
		TrainClcikInsert(fileName, v) //companies
	}
}

//TrainClcikInsert ...
func TrainClcikInsert(fileName string, v *model.TrainClick) {
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
