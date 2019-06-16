package glob

import (
	"github.com/Hao1995/104hackathon/config"
	"github.com/astaxie/beego/logs"
)

func init() {

	if config.CfgData.Logs.Level == 0 || config.CfgData.Logs.Level > 7 {
		logs.Warn("Logs level is wrong = %v. Using 'debug' level.", config.CfgData.Logs.Level)
		logs.SetLevel(logs.LevelDebug)
	} else {
		logs.SetLevel(config.CfgData.Logs.Level)
	}

	logs.SetLogFuncCallDepth(3) //print the position that we wrote the logs

	// Logs example
	// logs.Debug("my book is bought in the year of ", 2016)
	// logs.Info("this %s cat is %v years old", "yellow", 3)
	// logs.Notice("Notice check")
	// logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	// logs.Error(1024, "is a very", "good game")
	// logs.Alert("Alert Test")
	// logs.Critical("oh,crash")
	// logs.Emergency("Emergency test")
}
