package main

/*
const char* build_time(void)
{
static const char* psz_build_time = "["__DATE__ "  " __TIME__ "]";
return psz_build_time;
}
*/
// import "C"
import (
	"fmt"
	"logAnalysis/conf"
	"logAnalysis/manager"
	"mmx/mediaservice/log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

const VERSION = "1.0.0.0"

// var BUILDTIME = C.GoString(C.build_time())

//[Sep  7 2018  11:15:42]
//程序编译时间，通过上面调用c代码实现
// func BuildTime() string {
// 	t, _ := time.Parse("[Jan  2 2006  15:04:05]", BUILDTIME)
// 	return t.Format("2006-01-02 15:04:05")
// }

func init() {
	//版本号和编译时间写入文件
	f, err := os.OpenFile("./logAnalysis.version", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0444)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString("serverName:logAnalysis\n")
	f.WriteString("version:")
	f.WriteString(VERSION + "\n")
	// f.WriteString("buildTime:")
	// f.WriteString(BuildTime() + "\n")
	f.WriteString("startTime:")
	f.WriteString(time.Now().Format("2006-01-02 15:04:05") + "\n")
}
func main() {
	//配置文件初始化
	conf.Init("../conf/config.json")
	//本服务日志模块初始化
	logger := log.Start(conf.Cf.LogPath, conf.Cf.FileName, "5")
	if logger == nil {
		return
	}
	defer logger.Stop()

	//开始监控所有模块并建立索引
	manager.StartMonitor()

	//定时删除非今天的日志文件
	go manager.DeleteOldLogFiles(conf.GetAllServerPaths())

	//启动http日志查询服务
	http.HandleFunc("/search", manager.GetAllLogs)
	fmt.Println("http listen on ", conf.Cf.ListenPort)
	if err := http.ListenAndServe(":"+conf.Cf.ListenPort, nil); err != nil {
		panic(err)
	}
}
