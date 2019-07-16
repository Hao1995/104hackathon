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

	http.HandleFunc("/api/sync/departments", controller.SyncDepartments)
	http.HandleFunc("/api/sync/companies", controller.SyncCompanies)
	http.HandleFunc("/api/sync/jobs", controller.SyncJobs)
	http.HandleFunc("/api/sync/train_click", controller.SyncTrainClick)
	http.HandleFunc("/api/sync/train_action", controller.SyncTrainAction)

	// http.HandleFunc("/104hackathon/train_click/sync/key", controller.QueryKey)
	http.HandleFunc("/104hackathon/query_key/sync", controller.StoreQueryKey)
	http.HandleFunc("/104hackathon/job_category/sync", controller.InsertJobCategory)
	// http.HandleFunc("/104hackathon/department/sync", controller.InsertDepartment)
	http.HandleFunc("/104hackathon/district/sync", controller.InsertDistrict)
	http.HandleFunc("/104hackathon/industry/sync", controller.InsertIndustry)

	// - API for frontend
	http.HandleFunc("/104hackathon/score/area", controller.ScoreArea)

	// http.HandleFunc("/104hackathon/sync/jobkey", controller.SyncJobKey)

	// http.HandleFunc("/test", controller.Test)
	http.ListenAndServe(":8080", nil)

}
