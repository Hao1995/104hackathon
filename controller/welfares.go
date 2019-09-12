package controller

import (
	"fmt"
	"net/http"

	"github.com/Hao1995/104hackathon/glob"
	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

func Welfares(w http.ResponseWriter, req *http.Request) {

	ins := &WelfaresController{}
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

type WelfaresController struct {
}

func (c *WelfaresController) get(httpLib *utils.HTTPLib) {

	rows, err := glob.DB.Query("SELECT * FROM `welfares` ORDER BY `name`")
	if err != nil {
		logs.Error(err)
		httpLib.WriteJSON(err)
		return
	}
	defer rows.Close()

	items := []*models.WelfaresItem{}
	for rows.Next() {
		item := &models.WelfaresItem{}
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			logs.Error(err)
			httpLib.WriteJSON(err)
			return
		}
		items = append(items, item)
		// logs.Debug("id:%v, name:%v", *item.ID, *item.Name)
	}

	httpLib.WriteJSON(items)
}

func (c *WelfaresController) post(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	name := httpLib.Req.FormValueToNullString("name")
	if !name.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "name")
		httpLib.WriteJSON(res)
		return
	}

	// - Insert Data
	tx, err := glob.DB.Begin()
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	// Welfares
	welfareRes, err := tx.Exec("INSERT INTO `welfares` (`name`) VALUES (?)", name)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	welfareNo, err := welfareRes.LastInsertId()
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	// Get Users
	rows, err := tx.Query("SELECT `id` FROM `users`")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer rows.Close()

	items := []int{}
	for rows.Next() {
		var item int
		err := rows.Scan(&item)
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
		items = append(items, item)
	}

	// WelfareUserScore
	stmt, err := tx.Prepare("INSERT INTO `welfare_user_score` (`user_id`, `welfare_no`, `score`) VALUES (?, ?, ?)")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(item, welfareNo, 1)
		if err != nil {
			logs.Error(err)
			res.Error = err.Error()
			httpLib.WriteJSON(res)
			return
		}
	}

	tx.Commit()

	res.Message = fmt.Sprintf("Sucess insert welfare = {name:%v}. And sync to %v user.", name, len(items))
	httpLib.WriteJSON(res)
}

func (c *WelfaresController) delete(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	id := httpLib.Req.URL.Query().Get("id")

	// - Delete Data
	stmt, err := glob.DB.Prepare("DELETE FROM `welfares` WHERE `id` = ?")
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
