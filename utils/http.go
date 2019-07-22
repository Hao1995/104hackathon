package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type HTTPLib struct {
	Res http.ResponseWriter
	Req *Request
}

func (c *HTTPLib) Init(w http.ResponseWriter, req *http.Request) {
	c.Res = w
	c.Req = &Request{Request: req}
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

// http.ResponseWriter
type Request struct {
	*http.Request
}

func (c *Request) FormValueToNullInt64(key string) (sql.NullInt64, error) {

	if tmpVal := c.Request.FormValue(key); tmpVal == "" {
		return NewNullInt64(nil), nil
	} else {
		if res, err := strconv.ParseInt(tmpVal, 10, 64); err != nil {
			return NewNullInt64(nil), err
		} else {
			return NewNullInt64(&res), nil
		}
	}

	return NewNullInt64(nil), nil
}

func (c *Request) FormValueToNullString(key string) sql.NullString {

	if tmpVal := c.Request.FormValue(key); tmpVal == "" {
		return NewNullString(nil)
	} else {
		return NewNullString(&tmpVal)
	}
}
