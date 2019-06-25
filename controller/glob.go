package controller

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/Hao1995/104hackathon/config"

	_ "github.com/go-sql-driver/mysql" //mysql
)

var (
	db                *sql.DB
	dberr             error
	dbConnentCount    int
	dbConnentCountMax = 16382

	wg sync.WaitGroup
	mu sync.Mutex

	//FailFile Store fail file name. Then we can parse again for it.
	FailFile []string
)

func init() {

	// - Init DB
	// db, dberr = sql.Open([driver name], "[user name]:[user password]@tcp([mysql host])/")
	db, dberr = sql.Open("mysql", config.CfgData.Mysql.User+":"+config.CfgData.Mysql.Password+"@tcp("+config.CfgData.Mysql.Host+":"+config.CfgData.Mysql.Port+")/"+config.CfgData.Mysql.Name) //HP
	chechkErr(dberr)

	// max connection. 25%-50% of CPU threads
	db.SetMaxOpenConns(2)

	// - Init Schemas
	// sqlFiles, err := ioutil.ReadFile("./sql/init.sql")
	// if err != nil {
	// 	logs.Fatalf(": %s", err)
	// }

	// splitSQLFiles := strings.Split(string(sqlFiles), ";")

	// for _, v := range splitSQLFiles {
	// 	_, dberr = db.Exec(v)
	// 	chechkErr(dberr)
	// }
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

func chechkErr(err error) {
	if err != nil {
		fmt.Println("[ERROR] ", err)
	}
}
