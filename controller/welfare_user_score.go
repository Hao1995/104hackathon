package controller

import (
	"fmt"
	"net/http"

	"github.com/Hao1995/104hackathon/models"
	"github.com/Hao1995/104hackathon/utils"
	"github.com/astaxie/beego/logs"
)

func WelfareUserScore(w http.ResponseWriter, req *http.Request) {

	ins := &WelfareUserScoreController{}
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

type WelfareUserScoreController struct {
}

func (c *WelfareUserScoreController) get(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	httpLib.Req.ParseForm()
	userID, err := httpLib.Req.FormValueToNullInt64("user_id")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	if !userID.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "user_id")
		httpLib.WriteJSON(res)
		return
	}

	stmt, err := db.Prepare("SELECT `W`.`name`, `WU`.`score` FROM `welfare_user_score` AS `WU`, `welfares` AS `W` WHERE `WU`.`welfare_no` = `W`.`id` AND `WU`.`user_id` = ? ORDER BY `name`")
	if err != nil {
		logs.Error(err)
		httpLib.WriteJSON(err)
		return
	}
	defer stmt.Close()

	items := []*models.WelfareUserScoreItem{}
	rows, err := stmt.Query(userID)
	for rows.Next() {
		item := &models.WelfareUserScoreItem{}
		err := rows.Scan(&item.Name, &item.Score)
		if err != nil {
			logs.Error(err)
			httpLib.WriteJSON(err)
			return
		}
		items = append(items, item)
		// logs.Debug("name:%v, score:%v", *item.Name, *item.Score)
	}

	httpLib.WriteJSON(items)
	return
}

func (c *WelfareUserScoreController) post(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	httpLib.Req.ParseForm()
	userID, err := httpLib.Req.FormValueToNullInt64("user_id")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	if !userID.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "user_id")
		httpLib.WriteJSON(res)
		return
	}
	welfareNo, err := httpLib.Req.FormValueToNullInt64("welfare_no")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	if !welfareNo.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "welfare_no")
		httpLib.WriteJSON(res)
		return
	}
	score, err := httpLib.Req.FormValueToNullInt64("score")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	if !score.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "score")
		httpLib.WriteJSON(res)
		return
	}

	// - Insert Data
	stmt, err := db.Prepare("INSERT INTO `welfare_user_score` (`user_id`, `welfare_no`, `score`) VALUES (?, ?, ?)")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, welfareNo, score)
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}

	res.Message = fmt.Sprintf("Sucess insert {user_idx:%v, welfare_no:%v, score:%v}", userID, welfareNo, score)
	httpLib.WriteJSON(res)
	return
}

func (c *WelfareUserScoreController) delete(httpLib *utils.HTTPLib) {

	res := models.APIRes{}

	userID, err := httpLib.Req.FormValueToNullInt64("user_id")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	if !userID.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "user_id")
		httpLib.WriteJSON(res)
		return
	}
	welfareNo, err := httpLib.Req.FormValueToNullInt64("welfare_no")
	if err != nil {
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	if !welfareNo.Valid {
		res.Error = fmt.Sprintf("'%v' is necessary", "welfare_no")
		httpLib.WriteJSON(res)
		return
	}

	// - Delete Data
	stmt, err := db.Prepare("DELETE FROM `welfare_user_score` WHERE `user_id` = ? AND `welfare_no` = ?")
	if err != nil {
		logs.Error(err)
		res.Error = err.Error()
		httpLib.WriteJSON(res)
		return
	}
	defer stmt.Close()

	execRes, err := stmt.Exec(userID, welfareNo)
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
	return
}
