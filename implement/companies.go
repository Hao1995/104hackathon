package implement

import (
	"database/sql"
	"docker-example/model"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

//HackathonCompanies ...
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
