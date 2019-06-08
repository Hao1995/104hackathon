package controller

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Hao1995/104hackathon/log"
)

//InsertJobCategory ... Store Data
func InsertJobCategory(res http.ResponseWriter, req *http.Request) {

	filename := "C:/Users/user/Downloads/104Hackathon/category/job_category.csv"

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	lines = lines[1:len(lines)] //First line is title
	for _, line := range lines {

		queryString := "INSERT INTO job_category(`id`,`name`,`desc`,`hide`) VALUES"
		value :=
			"(" +
				stringAddSingleQuotation(processQuote(line[0])) + "," +
				stringAddSingleQuotation(processQuote(line[1])) + "," +
				stringAddSingleQuotation(processQuote(line[2])) + "," +
				stringAddSingleQuotation(processQuote(line[3])) +
				") ON DUPLICATE KEY UPDATE `id`=`id`,`name`=`name`,`desc`=`desc`,`hide`=`hide`;"

		queryString = queryString + value

		fmt.Println(queryString)

		stmt, err := db.Prepare(queryString)
		if err != nil {
			log.Errorf(err.Error())
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Errorf(err.Error())
		}
		stmt.Close()

	}

	io.WriteString(res, "Complete")
}
