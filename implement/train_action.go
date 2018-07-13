package implement

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Hao1995/Docker-Example/model"
)

//InsertTrainAction User
func InsertTrainAction(res http.ResponseWriter, req *http.Request) {
	directoryPath := "C:/Users/user/Downloads/104Hackathon/chunk/train_action" //train_action
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fileExtension := strings.Split(file.Name(), ".")
		if len(fileExtension) == 2 {
			if fileExtension[1] == "json" {
				filePath := directoryPath + "/" + file.Name()
				go func(filePath string) {
					ParseTrainAction(filePath)
				}(filePath)
			}
		}

	}

	io.WriteString(res, "Complete")
}

//ParseTrainAction ...
func ParseTrainAction(fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// c := []*model.Job{}
	c := []*model.TrainAction{}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		fmt.Println(err.Error())
		FailFile = append(FailFile, fileName)
		return
	}

	if len(c) > 0 {
		queryString := "INSERT INTO train_action(`action`, `jobno`, `date`, `source`, `device`) VALUES"

		fmt.Printf("[INFO] %v length : %v\n", fileName, len(c))
		for _, v := range c {
			// TrainActionInsert(fileName, v)
			value := "(" + stringAddDoubleQuotation(v.Action) + "," + stringAddDoubleQuotation(v.Jobno) + "," + stringAddDoubleQuotation(v.Date) + "," + stringAddDoubleQuotation(v.Source) + "," + stringAddDoubleQuotation(v.Device) + "),"
			queryString = queryString + value
		}

		queryString = strings.TrimRight(queryString, ",")
		// fmt.Println(queryString)
		// fmt.Println()

		TrainActionInsert(queryString, fileName)
	}

	fmt.Printf("[INFO] %v completed\n", fileName)
}

//TrainActionInsert ...
func TrainActionInsert(queryString, fileName string) {
	mu.Lock()
	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}

	stmt, err := db.Prepare(queryString)
	defer stmt.Close()
	dbConnentCount++

	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}
	chechkErr(err)

	_, err = stmt.Exec()
	dbConnentCount++
	mu.Unlock()
	if err != nil {
		fmt.Printf("[ERROR][%v][%v] Content :%v \n", fileName, err, queryString)
	}

	dbConnentCount--
	dbConnentCount--
}
