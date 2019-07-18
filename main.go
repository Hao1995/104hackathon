package main

import (
	"net/http"

	_ "github.com/Hao1995/104hackathon/config"
	_ "github.com/Hao1995/104hackathon/glob"

	"github.com/Hao1995/104hackathon/controller"
)

func main() {

	http.HandleFunc("/104hackathon/job", controller.HackathonJob)
	http.HandleFunc("/104hackathon/companies", controller.HackathonCompanies)

	http.HandleFunc("/api/sync/categories", controller.SyncCategories)
	http.HandleFunc("/api/sync/companies", controller.SyncCompanies)
	http.HandleFunc("/api/sync/jobs", controller.SyncJobs)
	http.HandleFunc("/api/sync/train_click", controller.SyncTrainClick)
	http.HandleFunc("/api/sync/train_action", controller.SyncTrainAction)

	http.HandleFunc("/api/user", controller.Users)
	http.HandleFunc("/api/welfare", controller.Welfares)

	http.HandleFunc("/104hackathon/query_key/sync", controller.StoreQueryKey)

	// - API for frontend
	http.HandleFunc("/104hackathon/score/area", controller.ScoreArea)

	// http.HandleFunc("/104hackathon/sync/jobkey", controller.SyncJobKey)

	http.HandleFunc("/test", controller.Test)
	http.ListenAndServe(":8080", nil)

}
