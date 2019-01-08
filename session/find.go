package session

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"logAnalysis/logindex"
	"mmx/mediaservice/log"
	"strings"
	"time"
)

var g_sessionDir string = "../data/index/session/"

var g_index logindex.LogIndex

func FindSessionLog(sessionArray []*Session) {
	for _, session := range sessionArray {
		key := strings.Replace(session.getCallid(), ":", "", -1)
		t1 := time.Now()
		callRouteLog := g_index.FindData(key, session.GetBeginTime(), "callroute")
		session.setCallRouteLog(callRouteLog)
		t2 := time.Now()
		fmt.Println("get callroute log", t2.Sub(t1))

		mcsCallerLog := g_index.FindData(session.getCallerUuid(), session.GetBeginTime(), "mcs")
		mcsCalledLog := g_index.FindData(session.getCalledUuid(), session.GetBeginTime(), "mcs")
		session.setMcsLog(mcsCallerLog)
		session.setMcsLog(mcsCalledLog)
		t3 := time.Now()
		fmt.Println("get mcs log", t3.Sub(t2))

		masCallerLog := g_index.FindData(session.getCallerSsrc(), session.GetBeginTime(), "mas")
		masCalledLog := g_index.FindData(session.getCalledSsrc(), session.GetBeginTime(), "mas")
		session.setMasLog(masCallerLog)
		session.setMasLog(masCalledLog)
		t4 := time.Now()
		fmt.Println("get mas log", t4.Sub(t3))

		sipmgwCallerLog := g_index.FindData(session.getCallerSsrc(), session.GetBeginTime(), "sipmgw")
		sipmgwCalledLog := g_index.FindData(session.getCalledSsrc(), session.GetBeginTime(), "sipmgw")
		session.setSipmgwLog(sipmgwCallerLog)
		session.setSipmgwLog(sipmgwCalledLog)
		t5 := time.Now()
		fmt.Println("get sipmgw log", t5.Sub(t4))

		sipsgwCallerLog := g_index.FindData(key, session.GetBeginTime(), "sipsgw")
		session.setSipsgwLog(sipsgwCallerLog)
		t6 := time.Now()
		fmt.Println("get sipsgw log", t6.Sub(t5))

	}
}

func findSession(sessionFile []string, caller string, called string,
	callBeginTime string, callEndTime string, keyWords string) []*Session {
	var sessionArray []*Session
	if len(sessionFile) == 0 {
		log.Warnf("match session File empty,caller:%s , called: %s \n", caller, called)
		return nil
	}
	var str []byte
	//可能匹配到两个月份的session文件，分别读取，都放入str中
	for _, file := range sessionFile {
		sessionByte, err := ioutil.ReadFile(file)
		if err != nil {
			log.Warnf("Read session File Fail:%s\n", err)
			continue
		}
		str = append(str, sessionByte...)
	}

	byCaller := []byte(caller)
	byCalled := []byte(called)
	if keyWords == "" {
		//如果关键字为空，则查询时忽略关键字，改为特殊字符串，下面查询时就不会匹配到
		keyWords = "==nokeyword=="
	}
	byKeyWords := []byte(keyWords)
	byBeginTime := []byte(callBeginTime)
	byEndTime := []byte(callEndTime)
	for _, line := range bytes.Split(str, []byte("\n")) {
		info := bytes.Split(line, []byte(","))
		//查询时，把callid包含_1结尾的忽略
		if len(info) >= 9 &&
			bytes.Compare(info[7], byBeginTime) >= 0 &&
			bytes.Compare(info[7], byEndTime) <= 0 &&
			!bytes.HasSuffix(info[2], []byte("_1")) {
			//如果主叫和被叫匹配
			//如果主叫被叫都有值，则需全部匹配
			if len(byCaller) != 0 && len(byCalled) != 0 {
				if bytes.Compare(info[0], byCaller) == 0 &&
					bytes.Compare(info[1], byCalled) == 0 ||
					bytes.Index(line, byKeyWords) >= 0 {
					session := new(Session)
					session.parse(string(line))
					sessionArray = append(sessionArray, session)
				}
			} else {
				//如果主叫被叫不都有值
				if bytes.Compare(info[0], byCaller) == 0 ||
					bytes.Compare(info[1], byCalled) == 0 ||
					bytes.Index(line, byKeyWords) >= 0 {
					session := new(Session)
					session.parse(string(line))
					sessionArray = append(sessionArray, session)
				}
			}

		}
	}
	return sessionArray
}

//因为日质量比较大，目前暂时只支持31天跨度的查询，最多会有两个session文件会匹配到，分别为开始月份的和结束月份的
func findSessionFile(callBeginTime string, callEndTime string) (matchSessionArr []string) {
	if len(callBeginTime) < 14 || len(callEndTime) < 14 {
		return
	}
	beginIndexTime := callBeginTime[:6]
	endIndexTime := callEndTime[:6]
	fileInfoArray, err := ioutil.ReadDir(g_sessionDir)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%s\n", err)
		return
	}
	for _, fileInfo := range fileInfoArray {
		if strings.Index(fileInfo.Name(), beginIndexTime) >= 0 || strings.Index(fileInfo.Name(), endIndexTime) >= 0 {
			matchSessionArr = append(matchSessionArr, g_sessionDir+fileInfo.Name())
		}
	}
	return
}

//Find Part, diff go
//timeFormat:2018030173500
func FindLog(caller string, called string,
	callBeginTime string, callEndTime string, keyWords string) []*Session {

	t1 := time.Now()
	matchSessionArr := findSessionFile(callBeginTime, callEndTime)
	log.Info("find sessionFile :", matchSessionArr)
	sessionArray := findSession(matchSessionArr, caller, called,
		callBeginTime, callEndTime, keyWords)
	FindSessionLog(sessionArray)

	t2 := time.Now()
	fmt.Printf("Cost time %v\n", t2.Sub(t1))
	return sessionArray

}
