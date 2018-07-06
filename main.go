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
	"docker-example/model"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db       *sql.DB
	dberr    error
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
		fmt.Println(v)
		_, dberr = db.Exec(v)
		chechkErr(dberr)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/read/users/", Read)
	http.HandleFunc("/read/users/json", ReadByJSON)

	// http.HandleFunc("/create", Create)

	http.HandleFunc("/insert", Insert)
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(res, req)
}

//User Model
type User struct {
	Id, Message, Name string
}

// ReadByJSON ...
func ReadByJSON(res http.ResponseWriter, req *http.Request) {
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

//Read User
func Read(res http.ResponseWriter, req *http.Request) {
	tmpl := template.Must(template.ParseFiles("users.html"))

	id := strings.TrimPrefix(req.URL.Path, "/read/users/")
	fmt.Println(id)

	var rows *sql.Rows
	var err error
	if id != "" {
		rows, err = db.Query("SELECT * FROM users WHERE id=?", id)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		rows, err = db.Query("SELECT * FROM users")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

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

//Create User
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

//Insert User
func Insert(res http.ResponseWriter, req *http.Request) {
	directoryPath := "F:/gotool/src/test/test1"
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, file := range files {
		fileExtension := strings.Split(file.Name(), ".")
		if len(fileExtension) == 2 {
			if fileExtension[1] == "json" {
				if count < 3 {
					filePath := directoryPath + "/" + file.Name()
					ParseJsonAndInsertToMySQL(filePath)
					count++
				} else {
					break
				}
			}
		}

	}
}

//ParseJsonAndInsertToMySQL ...
func ParseJsonAndInsertToMySQL(fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c := []*model.Job{}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		fmt.Println(err.Error())
		FailFile = append(FailFile, fileName)
		return
	}

	for _, v := range c {
		insert, err := db.Prepare("INSERT INTO job() VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		chechkErr(err)
		_, err = insert.Exec(v.Custno, v.Jobno, v.Job, v.Jobcat1, v.Jobcat2, v.Jobcat3, v.Edu, v.SalaryLow, v.SalaryHigh, v.Role, v.Language1, v.Language2, v.Language3, v.Period, v.MajorCat, v.MajorCat2, v.MajorCat3, v.Industry, v.Worktime, v.RoleStatus, v.S2, v.S3, v.Addrno, v.S9, v.NeedEmp, v.NeedEmp1, v.Startby, v.ExpJobcat1, v.ExpJobcat2, v.ExpJobcat3, v.Description, v.Others)
		if err != nil {
			fmt.Printf("[ERROR][%v][%v] Content :%v \n", fileName, err, *v)
		}
	}
}

func chechkErr(err error) {
	if err != nil {
		fmt.Println("[ERROR] ", err)
	}
}
