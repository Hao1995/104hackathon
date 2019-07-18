package utils

import (
	"database/sql"
	"strconv"
)

func StrToNewNullInt64(val *string) (sql.NullInt64, error) {
	valInt64, err := strconv.ParseInt(*val, 10, 64)
	if err != nil {
		return NewNullInt64(&valInt64), err
	}
	var valPtr *int64
	if valInt64 != 0 {
		valPtr = &valInt64
	}
	return NewNullInt64(valPtr), nil
}
