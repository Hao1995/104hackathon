package glob

import (
	"database/sql"

	"github.com/Hao1995/104hackathon/config"

	_ "github.com/go-sql-driver/mysql" //mysql
)

const (
	MySQLUpperPlaceholders = 65535
)

var (
	DB *sql.DB
)

func init() {

	// - Init DB
	// db, dberr = sql.Open([driver name], "[user name]:[user password]@tcp([mysql host])/")
	var err error
	DB, err = sql.Open("mysql", config.CfgData.Mysql.User+":"+config.CfgData.Mysql.Password+"@tcp("+config.CfgData.Mysql.Host+":"+config.CfgData.Mysql.Port+")/"+config.CfgData.Mysql.Name) //HP
	if err != nil {
		panic(err)
	}

	// max connection. 25%-50% of CPU threads
	DB.SetMaxOpenConns(2)

	// - Init Schemas
	// sqlFiles, err := ioutil.ReadFile("./sql/init.sql")
	// if err != nil {
	// 	logs.Fatalf(": %s", err)
	// }

	// splitSQLFiles := strings.Split(string(sqlFiles), ";")

	// for _, v := range splitSQLFiles {
	// 	_, dberr = DB.Exec(v)
	// 	chechkErr(dberr)
	// }
}
