package manager

import (
	"encoding/json"
	"fmt"
	"logAnalysis/session"

	"net/http"

	"io/ioutil"
	"logAnalysis/conf"
	"mmx/mediaservice/log"
	"os"
	"strings"
	"time"
)

type ResponseLog struct {
	Caller    string `json:"caller"`
	Called    string `json:"called"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Mcs       string `json:"mcs"`
	Callroute string `json:"callroute"`
	Mas       string `json:"mas"`
	SipmGW    string `json:"sipmgw"`
	SipsGW    string `json:"sipsgw"`
}
type Response struct {
	ErrCode int           `json:"errcode"`
	Message string        `json:"message"`
	Data    []ResponseLog `json:"data"`
}

func SendHttp(w http.ResponseWriter, code int, message string, data []ResponseLog) {
	var response Response
	response.ErrCode = code
	response.Message = message
	response.Data = data
	resJson, _ := json.Marshal(response)
	//newr := bytes.Replace(resJson,[]byte("\\n"),[]byte("\n"),-1)
	w.Write(resJson)
}

func GetAllLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Warnf("request method should be GET, now is %s", r.Method)
		w.Write([]byte("request method should be GET"))
		return
	}
	vars := r.URL.Query()
	startTime := vars.Get("startTime")
	endTime := vars.Get("endTime")
	caller := vars.Get("caller")
	called := vars.Get("called")
	keyword := vars.Get("keyWord")

	var alllog []ResponseLog
	w.Header().Set("charset", "utf-8")
	//w.Header().Set("Content-Type", "application/json")

	//检查参数
	if startTime == "" || len(startTime) != 14 {
		SendHttp(w, -1, "开始时间格式不正确！", alllog)
		return
	}
	if endTime == "" || len(endTime) != 14 {
		SendHttp(w, -2, "结束时间格式不正确！", alllog)
		return
	}
	if keyword == "" && caller == "" && called == "" {
		SendHttp(w, -3, "主叫被叫关键字不能都为空！", alllog)
		return
	}
	if strings.Compare(startTime, endTime) > 0 {
		SendHttp(w, -4, "开始时间不能大于结束时间！", alllog)
		return
	}
	start, err := time.Parse("20060102150405", startTime)
	if err != nil {
		SendHttp(w, -5, "开始时间解析不正确！", alllog)
		return
	}
	end, err := time.Parse("20060102150405", endTime)
	if err != nil {
		SendHttp(w, -6, "结束时间解析不正确！", alllog)
		return
	}
	if end.Sub(start) > time.Hour*24*31 {
		SendHttp(w, -7, "时间跨度不能超过一个月！", alllog)
		return
	}
	t6 := time.Now()
	fmt.Println("request...", t6)
	//找出所有session
	sessionArray := session.FindLog(caller, called, startTime, endTime, keyword)
	if len(sessionArray) == 0 {
		SendHttp(w, 0, "no log found!", alllog)
	} else {
		for _, session := range sessionArray {
			//callrouteLog := strings.Replace(string(session.GetCallRouteLog()),"\n","\\n",-1)
			//mcsLog := strings.Replace(string(session.GetCallMcsLog()),"\n","\\n",-1)
			//masLog := strings.Replace(string(session.GetCallMasLog()),"\n","\\n",-1)
			//sipmGWLog := strings.Replace(string(session.GetSipmgwLog()),"\n","\\n",-1)
			//sipsGWLog := strings.Replace(string(session.GetSipsgwLog()),"\n","\\n",-1)
			logs := ResponseLog{
				Caller:    session.GetCaller(),
				Called:    session.GetCalled(),
				StartTime: session.GetBeginTime(),
				EndTime:   session.GetEndTime(),
				Callroute: string(session.GetCallRouteLog()),
				Mcs:       string(session.GetCallMcsLog()),
				Mas:       string(session.GetCallMasLog()),
				SipmGW:    string(session.GetSipmgwLog()),
				SipsGW:    string(session.GetSipsgwLog()),
			}
			alllog = append(alllog, logs)
		}
		//w.Write([]byte("sssss"))
		SendHttp(w, 0, "success", alllog)
	}
	t7 := time.Now()
	fmt.Println("request cost time : ", t7.Sub(t6))
}

//启动监控和索引
func StartMonitor() {
	//监控各模块日志并建立索引
	callrouteDir, err := conf.GetServerPath("callroute")
	if err != nil {
		log.Fatal("can not get callroute log dir ", err)
		panic(err)
	}
	mcsDir, err := conf.GetServerPath("mcs")
	if err != nil {
		log.Fatal("can not get mcs log dir ", err)
		panic(err)
	}
	masDir, err := conf.GetServerPath("mas")
	if err != nil {
		log.Fatal("can not get mas log dir ", err)
		panic(err)
	}
	sipmgwDir, err := conf.GetServerPath("sipmgw")
	if err != nil {
		log.Fatal("can not get sipmgw log dir ", err)
		panic(err)
	}
	sipsgwDir, err := conf.GetServerPath("sipsgw")
	if err != nil {
		log.Fatal("can not get sipsgw log dir ", err)
		panic(err)
	}
	var sm session.SessManager
	go sm.Start(callrouteDir)
	var mcs session.McsManager
	go mcs.Start(mcsDir)
	var mas session.MasManager
	go mas.Start(masDir)
	var sipmgw session.SipmgwManager
	go sipmgw.Start(sipmgwDir)
	var sipsgw session.SipsgwManager
	go sipsgw.Start(sipsgwDir)
}

//删除非今天的日志文件，只保留索引,30分钟检查一次
func DeleteOldLogFiles(paths []string) {
	for {
		time.Sleep(30 * time.Minute)
		today := time.Now().Format("2006-01-02")
		for _, path := range paths {
			//每个服务日志的目录
			serverDirs, err := ioutil.ReadDir(path)
			if err != nil {
				log.Warn("walk dir failed : ", path)
				continue
			}
			for _, serverDir := range serverDirs {
				dateDirs, err := ioutil.ReadDir(path + serverDir.Name())
				if err != nil {
					log.Warn("walk dir failed : ", path)
					continue
				}
				for _, dateDir := range dateDirs {
					if dateDir.Name() != today {
						os.RemoveAll(path + serverDir.Name() + "/" + dateDir.Name())
						log.Info("remove old logs :", path+serverDir.Name()+"/"+dateDir.Name())
					}
				}
			}
		}
	}
}
