package session

import (
	"fmt"
	"io/ioutil"
	"logAnalysis/conf"
	"logAnalysis/logindex"
	"logAnalysis/monitor"
	"os"
	"strings"
)

type SessManager struct {
	m_sessionDir string
	m_sessionMap map[string]*Session
	m_index      logindex.LogIndex
}

func (sm *SessManager) addSipCall(line string) {
	session := sm.getSession(line)
	if session != nil {
		session.setCallid(parseCallId(line))
		session.setCaller(parseSipCaller(line))
		session.setCalled(parseSipCalled(line))
		session.setCallBeginTime(parseTime(line))
	}
}

func (sm *SessManager) addPbCall(line string) {
	session := sm.getSession(line)
	if session != nil {
		session.setCallid(parseCallId(line))
		session.setCaller(parsePbCaller(line))
		session.setCalled(parsePbCalled(line))
		session.setCallBeginTime(parseTime(line))
	}
}

func (sm *SessManager) addUuidAndSsrc(line string) {
	session := sm.getSession(line)
	if session != nil {
		session.setSsrc(paseSsrc(line))
		session.setUuid(parseUuid(line))
	}
}

func (sm *SessManager) addEndTime(line string) {
	session := sm.getSession(line)
	if session != nil {
		session.setCallEndTime(parseTime(line))
	}
}

func (sm *SessManager) handleData(data []string, filePath string) {

	typeFilter := "\"type\":\"0\""
	uuidFilter := "\"method\":\"release\""
	sipCallerFilter := "msg.m_Caller="
	pbCallerFilter := "pCallEvent->caller()="
	sessionCloesFilter := "Session will be deleted and close"
	for _, line := range data {
		if strings.Index(line, pbCallerFilter) >= 0 {
			sm.addPbCall(line)
		} else if strings.Index(line, sipCallerFilter) >= 0 {
			sm.addSipCall(line)
		} else if strings.Index(line, typeFilter) >= 0 &&
			strings.Index(line, uuidFilter) >= 0 {
			sm.addUuidAndSsrc(line)
		} else if strings.Index(line, sessionCloesFilter) >= 0 {
			sm.addEndTime(line)
		}
	}
	sm.writeSession()
	sm.m_index.BuildIndex(data, "callroute", parseCallId, parseCallRouteTime)
}

func (sm *SessManager) getSession(line string) *Session {
	var session *Session
	callid := parseCallId(line)
	if len(callid) == 0 {
		return nil
	}
	_, exsit := sm.m_sessionMap[callid]
	if exsit {
		session = sm.m_sessionMap[callid]
	} else {
		session = new(Session)
		sm.m_sessionMap[callid] = session
	}
	return session
}

func (sm *SessManager) writeTempSession() {
	path := fmt.Sprintf("%s.index", sm.m_sessionDir)
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Open writeSessionMapFile File faild\n")
		return
	}
	for _, session := range sm.m_sessionMap {
		fd.WriteString(session.format())
		fd.WriteString("\n")
	}
	fd.Close()
}

func (sm *SessManager) writeCompleteSession() {

	comSessionMap := make(map[string][]string) //fileName, sessionFormat
	for key, session := range sm.m_sessionMap {
		if len(session.GetEndTime()) > 0 {
			filePath := fmt.Sprintf("%s%s%s.index", sm.m_sessionDir,
				session.getBeginTimeYear(), session.getBeginTimeMonth())
			comSessionMap[filePath] = append(comSessionMap[filePath], session.format())
			delete(sm.m_sessionMap, key)
		}
	}
	for filePath, sessionStrArray := range comSessionMap {
		fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Open writeSessionMapFile File faild\n")
			continue
		}
		for _, sessionStr := range sessionStrArray {
			fd.WriteString(sessionStr)
			fd.WriteString("\n")
		}
		fd.Close()
	}
}

func (sm *SessManager) writeSession() {
	err := os.MkdirAll(sm.m_sessionDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	sm.writeCompleteSession()
	sm.writeTempSession()
}

func (sm *SessManager) readSession() {

	path := fmt.Sprintf("%s.index", sm.m_sessionDir)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("readSessionMapFile:%s\n", err)
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		session := new(Session)
		if session.parse(line) == true {
			sm.m_sessionMap[session.getCallid()] = session
		}
	}
}

func (sm *SessManager) Start(dir string) {
	sm.m_sessionDir = "../data/index/session/"
	sm.m_sessionMap = make(map[string]*Session)
	sm.readSession()

	var SessionMonitor monitor.FileMonitor
	logVersion, _ := conf.GetServerLogVersion("callroute")
	//根据日志格式的不同，使用不同的处理方法
	switch logVersion {
	case 1:
		SessionMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, sm.handleData)
	default:
		SessionMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, sm.handleData)
	}
}
