package config

import (
	"io/ioutil"
	"log"

	gcfg "gopkg.in/gcfg.v1"
)

//Cfg : Configure Struct
type Cfg struct {
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     string
	}
	Logs struct {
		Level int
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

	err = gcfg.ReadStringInto(&CfgData, string(appConf))
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}
}
