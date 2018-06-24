package config

import (
	"fmt"
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

	fmt.Println(string(appConf))
	err = gcfg.ReadStringInto(&CfgData, string(appConf))
	if err != nil {
		log.Fatalf("Failed to parse gcfg data: %s", err)
	}
	fmt.Println(CfgData.Mysql.User)
	fmt.Println(CfgData.Mysql.Password)
}
