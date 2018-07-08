package implement

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

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
