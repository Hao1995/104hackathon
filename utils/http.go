package utils

import (
	"encoding/json"
	"net/http"
)

type HTTPLib struct {
	Res http.ResponseWriter
	Req *http.Request
}

func (c *HTTPLib) Init(w http.ResponseWriter, req *http.Request) {
	c.Res = w
	c.Req = req
	return
}

func (c *HTTPLib) WriteJSON(msg interface{}) {

	c.Res.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(msg)
	if err != nil {
		http.Error(c.Res, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Res.Write(js)
}
