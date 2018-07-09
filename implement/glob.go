package implement

import (
	"database/sql"
	"docker-example/config"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql" //mysql
)

var (
	db                *sql.DB
	dberr             error
	dbConnentCount    int
	dbConnentCountMax = 16382

	mu sync.Mutex

	//FailFile Store fail file name. Then we can parse again for it.
	FailFile []string
)

func init() {
	// db, dberr = sql.Open([driver name], "[user name]:[user password]@tcp([mysql host])/")
	db, dberr = sql.Open("mysql", config.CfgData.Mysql.User+":"+config.CfgData.Mysql.Password+"@tcp("+config.CfgData.Mysql.Host+":"+config.CfgData.Mysql.Port+")/") //HP
	chechkErr(dberr)

	sqlFiles, err := ioutil.ReadFile("./sql/init.sql")
	if err != nil {
		log.Fatalf(": %s", err)
	}

	splitSQLFiles := strings.Split(string(sqlFiles), ";")

	for _, v := range splitSQLFiles {
		// fmt.Println(v)
		_, dberr = db.Exec(v)
		chechkErr(dberr)
	}
}

func stringAddDoubleQuotation(str string) string {
	return "\"" + str + "\""
}

func stringAddSingleQuotation(str string) string {
	return "'" + str + "'"
}

func processQuote(str string) string {
	return strings.Replace(strings.Replace(str, "'", "\\'", -1), "\"", "\\\"", -1)
}

// func(str string) processQuote() string{
// 	str
// }

func chechkErr(err error) {
	if err != nil {
		fmt.Println("[ERROR] ", err)
	}
}
