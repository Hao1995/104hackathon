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
	case http.MethodDelete:
	default:
		res.Error = fmt.Sprintf("There is no way correspond to method '%v'", req.Method)
		httpLib.WriteJSON(res)
		return
	}
}

type UsersController struct {
}

func (c *UsersController) get(httpLib *utils.HTTPLib) {

	rows, err := db.Query("SELECT * FROM `users`")
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
		fmt.Println(item)
	}

	httpLib.WriteJSON(items)
}
