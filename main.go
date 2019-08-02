package main

import (
	"net/http"

	_ "github.com/Hao1995/104hackathon/config"
	_ "github.com/Hao1995/104hackathon/glob"

	"github.com/Hao1995/104hackathon/controller"
)

func main() {

	// - Sync Data From 104 Open Data
	http.HandleFunc("/api/sync/categories", controller.SyncCategories)
	http.HandleFunc("/api/sync/companies", controller.SyncCompanies)
	http.HandleFunc("/api/sync/jobs", controller.SyncJobs)
	http.HandleFunc("/api/sync/train_click", controller.SyncTrainClick)
	http.HandleFunc("/api/sync/train_click/key", controller.SyncTrainClickKey)
	http.HandleFunc("/api/sync/train_action", controller.SyncTrainAction)

	// - Add Extra Data
	http.HandleFunc("/api/user", controller.Users)
	http.HandleFunc("/api/welfare", controller.Welfares)
	http.HandleFunc("/api/user/welfare/score", controller.WelfareUserScore)
	http.HandleFunc("/api/job/welfare", controller.JobWelfares)
	http.HandleFunc("/api/user/job/score", controller.JobUserScore)

	// - API for frontend
	http.HandleFunc("/104hackathon/jobs", controller.SearchJobs)

	http.HandleFunc("/test", controller.Test)
	http.ListenAndServe(":8080", nil)

}
