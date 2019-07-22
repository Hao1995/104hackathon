package controller

import (
	"fmt"
	"net/http"

	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

func Users(w http.ResponseWriter, req *http.Request) {

	ins := &UsersController{}
	httpLib := &utils.HTTPLib{}
	httpLib.Init(w, req)

	res := models.APIRes{}

	switch req.Method {
	case http.MethodGet:
		ins.get(httpLib)
	case http.MethodPost:
		ins.post(httpLib)
	case http.MethodDelete:
		ins.delete(httpLib)
	default:
		res.Error = fmt.Sprintf("There is no way correspond to method '%v'", req.Method)
		httpLib.WriteJSON(res)
		return
	}
}

type UsersController struct {
}

func (c *UsersController) get(httpLib *utils.HTTPLib) {

	rows, err := db.Query("SELECT * FROM `users` ORDER BY `name")
	if err != nil {
		logs.Error(err)
		httpLib.WriteJSON(err)
		return
	}
	defer rows.Close()

	items := []*models.UsersItem{}
	for rows.Next() {
		item := &models.UsersItem{}
		err := rows.Scan(&item.ID, &item.Name, &item.Email)
		if err != nil {
			logs.Error(err)
			httpLib.WriteJSON(err)
			return
		}
		items = append(items, item)
		// logs.Debug("id:%v, name:%v, email:%v", *item.ID, *item.Name, *item.Email)
	}

	httpLib.WriteJSON(items)
}

func (c *UsersController) post(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	httpLib.Req.ParseForm()
	name := httpLib.Req.FormValueToNullString("name")
	if !name.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "name")
		httpLib.WriteJSON(res)
		return
	}
	email := httpLib.Req.FormValueToNullString("email")
	if !email.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "email")
		httpLib.WriteJSON(res)
		return
	}

	// - Insert Data
	stmt, err := db.Prepare("INSERT INTO `users` (`name`, `email`) VALUES (?, ?)")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, email)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	res.Message = fmt.Sprintf("Sucess insert {name:%v, email:%v}", name, email)
	httpLib.WriteJSON(res)
}

func (c *UsersController) delete(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	id := httpLib.Req.URL.Query().Get("id")

	// - Delete Data
	stmt, err := db.Prepare("DELETE FROM `users` WHERE `id` = ?")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer stmt.Close()

	execRes, err := stmt.Exec(id)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	num, err := execRes.RowsAffected()
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	res.Message = fmt.Sprintf("Delete %v data", num)
	httpLib.WriteJSON(res)
}
