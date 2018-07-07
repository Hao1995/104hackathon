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
	"sync"

	"docker-example/config"
	"docker-example/model"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db                *sql.DB
	dberr             error
	dbConnentCount    int
	dbConnentCountMax = 16382

	mu sync.Mutex

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

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/read/users/", Read)
	http.HandleFunc("/read/users/json", ReadByJSON)

	http.HandleFunc("/104hackathon/job", HackathonJob)
	http.HandleFunc("/104hackathon/companies", HackathonCompanies)

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
	// directoryPath := "F:/gotool/src/test/test1" //job
	directoryPath := "F:/gotool/src/test/test1/data/companies" //companies
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
					ParseJSONAndInsertToMySQL(filePath)
				}(filePath)
			}
		}

	}

	io.WriteString(res, "Complete")
}

//ParseJSONAndInsertToMySQL ...
func ParseJSONAndInsertToMySQL(fileName string) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// c := []*model.Job{}
	c := []*model.Company{}
	err = json.Unmarshal(raw, &c)
	if err != nil {
		fmt.Println(err.Error())
		FailFile = append(FailFile, fileName)
		return
	}

	for _, v := range c {
		// InsertToJob(fileName, v) //job
		InsertToCompanies(fileName, v) //companies
	}
}

//InsertToJob ...
func InsertToJob(fileName string, v *model.Job) {
	mu.Lock()
	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}
	stmt, err := db.Prepare("INSERT INTO job() VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	dbConnentCount++

	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}
	chechkErr(err)
	_, err = stmt.Exec(v.Custno, v.Jobno, v.Job, v.Jobcat1, v.Jobcat2, v.Jobcat3, v.Edu, v.SalaryLow, v.SalaryHigh, v.Role, v.Language1, v.Language2, v.Language3, v.Period, v.MajorCat, v.MajorCat2, v.MajorCat3, v.Industry, v.Worktime, v.RoleStatus, v.S2, v.S3, v.Addrno, v.S9, v.NeedEmp, v.NeedEmp1, v.Startby, v.ExpJobcat1, v.ExpJobcat2, v.ExpJobcat3, v.Description, v.Others)
	dbConnentCount++
	mu.Unlock()
	if err != nil {
		fmt.Printf("[ERROR][%v][%v] Content :%v \n", fileName, err, *v)
	}

	dbConnentCount--
	dbConnentCount--
}

//InsertToCompanies ...
func InsertToCompanies(fileName string, v *model.Company) {
	mu.Lock()
	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}
	stmt, err := db.Prepare("INSERT INTO companies(`custno`,`invoice`,`name`,`profile`,`management`,`welfare`,`product`) VALUES(?,?,?,?,?,?,?)")
	defer stmt.Close()
	dbConnentCount++

	for {
		if dbConnentCount < dbConnentCountMax {
			break
		}
	}
	chechkErr(err)
	_, err = stmt.Exec(v.Custno, v.Invoice, v.Name, v.Profile, v.Management, v.Welfare, v.Product)
	dbConnentCount++
	mu.Unlock()
	if err != nil {
		fmt.Printf("[ERROR][%v][%v] Content :%v \n", fileName, err, *v)
	}

	dbConnentCount--
	dbConnentCount--
}

//HackathonJob User
func HackathonJob(res http.ResponseWriter, req *http.Request) {

	//=====Params
	req.ParseForm()
	params := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "size":
			params[k] = strings.Join(v, "")
			// case "message":
			// 	params[k] = strings.Join(v, "")
		}
	}

	var rows *sql.Rows
	var err error
	if v, ok := params["size"]; ok {
		rows, err = db.Query("SELECT * FROM job LIMIT " + v.(string))
	} else {
		rows, err = db.Query("SELECT * FROM job LIMIT 100")
	}

	jobs := []*model.Job{}

	for rows.Next() {
		r := &model.Job{}

		err = rows.Scan(&r.Custno, &r.Jobno, &r.Job, &r.Jobcat1, &r.Jobcat2, &r.Jobcat3, &r.Edu, &r.SalaryLow, &r.SalaryHigh, &r.Role, &r.RoleStatus, &r.Language1, &r.Language2, &r.Language3, &r.Period, &r.MajorCat, &r.MajorCat2, &r.MajorCat3, &r.Industry, &r.Worktime, &r.RoleStatus, &r.S2, &r.S3, &r.Addrno, &r.S9, &r.NeedEmp, &r.NeedEmp1, &r.Startby, &r.ExpJobcat1, &r.ExpJobcat2, &r.ExpJobcat3, &r.Description, &r.Others)
		chechkErr(err)
		jobs = append(jobs, r)
	}

	jsonData, err := json.Marshal(jobs)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}

//HackathonCompanies User
func HackathonCompanies(res http.ResponseWriter, req *http.Request) {

	//=====Params
	req.ParseForm()
	params := make(map[string]interface{})
	for k, v := range req.Form {
		switch k {
		case "size":
			params[k] = strings.Join(v, "")
			// case "message":
			// 	params[k] = strings.Join(v, "")
		}
	}

	var rows *sql.Rows
	var err error
	if v, ok := params["size"]; ok {
		rows, err = db.Query("SELECT * FROM companies LIMIT " + v.(string))
	} else {
		rows, err = db.Query("SELECT * FROM companies LIMIT 100")
	}

	companies := []*model.Company{}

	for rows.Next() {
		r := &model.Company{}

		err = rows.Scan(&r.Custno, &r.Invoice, &r.Name, &r.Profile, &r.Management, &r.Welfare, &r.Product)
		chechkErr(err)
		companies = append(companies, r)
	}

	jsonData, err := json.Marshal(companies)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}

func chechkErr(err error) {
	if err != nil {
		fmt.Println("[ERROR] ", err)
	}
}
