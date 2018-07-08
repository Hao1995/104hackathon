package main

import (
	"html/template"
	"net/http"

	"docker-example/implement"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/read/users/", implement.Read)
	http.HandleFunc("/read/users/json", implement.ReadByJSON)
	http.HandleFunc("/create", implement.Create)

	http.HandleFunc("/104hackathon/job", implement.HackathonJob)
	http.HandleFunc("/104hackathon/companies", implement.HackathonCompanies)
	http.HandleFunc("/104hackathon/train_click", implement.HackathonTrainClick)

	// http.HandleFunc("/insert", Insert)
	// http.HandleFunc("/insert/train_click", implement.InsertTrainClick)
	// http.HandleFunc("/insert/train_action", InsertTrainAction)
	http.HandleFunc("/104hackathon/train_click/sync/key", implement.QueryKey)

	// http.HandleFunc("/test", implement.Test)
	http.ListenAndServe(":8080", nil)

}

func index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(res, req)
}
