package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"docker-example/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var dberr error

func init() {
	// db, dberr = sql.Open([driver name], "[user name]:[user password]@tcp([mysql host])/")
	db, dberr = sql.Open("mysql", config.CfgData.Mysql.User+":"+config.CfgData.Mysql.Password+"@tcp("+config.CfgData.Mysql.Host+":"+config.CfgData.Mysql.Port+")/") //HP
	chechkErr(dberr)

	sqlFiles, err := ioutil.ReadFile("./sql/init.sql")
	if err != nil {
		log.Fatalf(": %s", err)
	}

	splitSQLFiles := strings.Split(string(sqlFiles), ";\n")

	for _, v := range splitSQLFiles {
		fmt.Println(v)
		_, dberr = db.Exec(v)
		chechkErr(dberr)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/read/users", Read)
	http.HandleFunc("/read/users/json", ReadByJson)

	http.HandleFunc("/create", Create)
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(res, req)
}

type User struct {
	Id, Message, Name string
}

func ReadByJson(res http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT * FROM users")

	users := []*User{}

	for rows.Next() {
		r := &User{}

		err = rows.Scan(&r.Id, &r.Name, &r.Message)
		chechkErr(err)
		users = append(users, r)
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}

type Todo struct {
	Title string
	Done  bool
}
type TodoPageData struct {
	Todos []Todo
}

func Read(res http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("users.html"))

	rows, err := db.Query("SELECT * FROM users")

	data := struct {
		Users []*User
	}{}

	for rows.Next() {
		r := &User{}

		err = rows.Scan(&r.Id, &r.Name, &r.Message)
		chechkErr(err)
		data.Users = append(data.Users, r)
	}
	tmpl.Execute(res, data)
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
		fmt.Println("[ERROR] ", err)
	}
}
