package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var dberr error

func init() {
	db, dberr = sql.Open("mysql", "root:password@tcp(localhost:3306)/")
	// db, dberr = sql.Open("mysql", "root:hao825_MDL7519@tcp(localhost:3306)/")
	chechkErr(dberr)

	_, dberr = db.Exec("CREATE DATABASE IF NOT EXISTS users")
	chechkErr(dberr)

	_, dberr = db.Exec("USE users")
	chechkErr(dberr)

	_, dberr = db.Exec("CREATE TABLE IF NOT EXISTS users ( id INT NOT NULL AUTO_INCREMENT, name VARCHAR(20), message VARCHAR(200), PRIMARY KEY (id))")
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
	tableDatas := make([]map[string]interface{}, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		chechkErr(err)
		tableData := make(map[string]interface{})

		for i, col := range values {
			var value interface{}
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			tableData[columns[i]] = value
		}
		tableDatas = append(tableDatas, tableData)
	}

	jsonData, err := json.Marshal(tableDatas)
	chechkErr(err)
	fmt.Println(string(jsonData))
	io.WriteString(res, string(jsonData))
}

func Create(res http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	user := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "name":
			user[k] = strings.Join(v, "")
		case "message":
			user[k] = strings.Join(v, "")
		}
	}

	insert, err := db.Prepare("INSERT users SET name=?,message=?")
	chechkErr(err)
	_, err = insert.Exec(user["name"], user["message"])
	chechkErr(err)

	str := "<h1>Success Insert</h1> <h3>Name: " + user["name"].(string) + "</h3>" + "<h3>Message: " + user["message"].(string) + "</h3>" + "\n\n" + "<a href=\"/\">Come back to home page</a>"
	io.WriteString(res, str)
}

func chechkErr(err error) {
	if err != nil {
		panic(err)
	}
}
