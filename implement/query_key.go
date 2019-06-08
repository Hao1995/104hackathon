package implement

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Hao1995/104hackathon/log"
)

//StoreQueryKey ...
func StoreQueryKey(res http.ResponseWriter, req *http.Request) {

	type key struct {
		Key string `json:"key"`
	}
	//=====Search Key
	fmt.Println("===== Search Key")

	var rows *sql.Rows
	var err error

	rows, err = db.Query("SELECT DISTINCT(`key`) AS `key` FROM `train_click`;")

	keys := []*key{}

	for rows.Next() {
		r := &key{}

		err = rows.Scan(&r.Key)
		if err != nil {
			log.Errorf(err.Error())
		} else {
			// if strings.HasPrefix(r.Key, "\"") && strings.HasSuffix(r.Key, "\"") {

			// }
			r.Key = strings.Trim(r.Key, " ")
			r.Key = strings.Trim(r.Key, "\"")
			keys = append(keys, r)
		}
	}
	if rows.Err() != nil {
		log.Errorf(rows.Err().Error())
	}

	//=====Store Query Key
	fmt.Println("===== Store Query Key")

	queryStringStart := "INSERT INTO query_key(`name`) VALUES"
	// size := len(keys)
	limit := 500
	offest := 0

	queryString := queryStringStart

	count := 0
	for i := offest; i < len(keys); i++ {
		if i < offest+limit {
			count++
			v := keys[i]
			value := "(" + stringAddSingleQuotation(processQuote(v.Key)) + "),"
			queryString = queryString + value
		}
		if i == offest+limit-1 {
			queryString = strings.TrimRight(queryString, ",")

			queryString = queryString + " ON DUPLICATE KEY UPDATE `name`=`name`;"
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

			queryString = queryStringStart
			offest = offest + limit
		}
	}

	io.WriteString(res, "Complete")
}
