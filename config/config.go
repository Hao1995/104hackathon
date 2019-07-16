package config

import (
	"io/ioutil"
	"log"

	"github.com/astaxie/beego/logs"

	gcfg "gopkg.in/gcfg.v1"
)

//Cfg : Configure Struct
type Cfg struct {
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     string
		Name     string
	}
	Logs struct {
		Level     int
		Level_Abc string
	}
	Data struct {
		Companies      string
		Jobs           string
		Train_Click    string
		Train_Action   string
		Departments    string
		Districts      string
		Industries     string
		Job_Categories string
	}
}

//CfgData : Can be use by other package
var (
	CfgData = Cfg{}
)

func init() {

	appConf, err := ioutil.ReadFile("./app.conf")
	if err != nil {
		log.Fatalf("Failed to read app.conf file: %s", err)
	}

	logs.Debug(string(appConf))

	err = gcfg.ReadStringInto(&CfgData, string(appConf))
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}
}
