package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var dberr error

func init() {
	db, dberr = sql.Open("mysql", "root:hao825_MDL7519@/users")
	chechkErr(dberr)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/read", Read)
	http.HandleFunc("/create", Create)
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(res, req)
}
func Read(res http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT * FROM users")

	columns, err := rows.Columns()
	chechkErr(err)
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		chechkErr(err)
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}

}
func Create(res http.ResponseWriter, req *http.Request) {
	stmt, err := db.Prepare("INSERT users SET name=?,message=?")
	chechkErr(err)
	dbres, err := stmt.Exec("harry", "test~~~")
	chechkErr(err)

	id, err := dbres.LastInsertId()
	chechkErr(err)

	io.WriteString(res, fmt.Sprintln(id))
}

func chechkErr(err error) {
	if err != nil {
		panic(err)
	}
}
