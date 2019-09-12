package controller

import (
	"strings"

	_ "github.com/go-sql-driver/mysql" //mysql
)

func stringAddDoubleQuotation(str string) string {
	return "\"" + str + "\""
}

func stringAddSingleQuotation(str string) string {
	return "'" + str + "'"
}

func processQuote(str string) string {
	return strings.Replace(strings.Replace(str, "'", "\\'", -1), "\"", "\\\"", -1)
}
