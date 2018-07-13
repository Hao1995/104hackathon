package implement

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Hao1995/Docker-Example/model"
)

//HackathonJob ...
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

		err = rows.Scan(&r.Custno, &r.Jobno, &r.Job, &r.Jobcat1, &r.Jobcat2, &r.Jobcat3, &r.Edu, &r.SalaryLow, &r.SalaryHigh, &r.Role, &r.Language1, &r.Language2, &r.Language3, &r.Period, &r.MajorCat, &r.MajorCat2, &r.MajorCat3, &r.Industry, &r.Worktime, &r.RoleStatus, &r.S2, &r.S3, &r.Addrno, &r.S9, &r.NeedEmp, &r.NeedEmp1, &r.Startby, &r.ExpJobcat1, &r.ExpJobcat2, &r.ExpJobcat3, &r.Description, &r.Others)
		chechkErr(err)
		jobs = append(jobs, r)
	}

	jsonData, err := json.Marshal(jobs)
	if err != nil {
		chechkErr(err)
	}
	io.WriteString(res, string(jsonData))
}
