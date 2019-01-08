package session

import (
	"bytes"
	"fmt"
	"strings"
)

type Session struct {
	m_caller        string
	m_called        string
	m_callId        string
	m_callerUuid    string
	m_calledUuid    string
	m_callerSsrc    string
	m_calledSsrc    string
	m_callBeginTime string
	m_callEndTime   string
	m_callRouteLog  bytes.Buffer
	m_mcsLog        bytes.Buffer
	m_masLog        bytes.Buffer
	m_sipmgwLog     bytes.Buffer
	m_sipsgwLog     bytes.Buffer
}

func (ss *Session) GetCallRouteLog() []byte {
	return ss.m_callRouteLog.Bytes()
}

func (ss *Session) GetCallMcsLog() []byte {
	return ss.m_mcsLog.Bytes()
}

func (ss *Session) GetCallMasLog() []byte {
	return ss.m_masLog.Bytes()
}
func (ss *Session) GetSipmgwLog() []byte {
	return ss.m_sipmgwLog.Bytes()
}
func (ss *Session) GetSipsgwLog() []byte {
	return ss.m_sipsgwLog.Bytes()
}

func (ss *Session) getBeginTimeMonth() string {
	if len(ss.m_callBeginTime) < 14 {
		return ""
	}
	return ss.m_callBeginTime[4:6]
}

func (ss *Session) getBeginTimeYear() string {
	if len(ss.m_callBeginTime) < 14 {
		return ""
	}
	return ss.m_callBeginTime[:4]
}

func (ss *Session) getCallid() string {
	return ss.m_callId
}

func (ss *Session) getCallerUuid() string {
	return ss.m_callerUuid
}

func (ss *Session) getCalledUuid() string {
	return ss.m_calledUuid
}

func (ss *Session) getCallerSsrc() string {
	return ss.m_callerSsrc
}

func (ss *Session) getCalledSsrc() string {
	return ss.m_calledSsrc
}

func (ss *Session) GetBeginTime() string {
	return ss.m_callBeginTime
}

func (ss *Session) GetCaller() string {
	return ss.m_caller
}
func (ss *Session) GetCalled() string {
	return ss.m_called
}

func (ss *Session) GetEndTime() string {
	return ss.m_callEndTime
}
func (ss *Session) setCallid(callid string) {
	ss.m_callId = callid
}

func (ss *Session) setCaller(caller string) {
	ss.m_caller = caller
}

func (ss *Session) setCalled(called string) {
	ss.m_called = called
}

func (ss *Session) setCallBeginTime(beginTime string) {
	ss.m_callBeginTime = beginTime
}

func (ss *Session) setCallEndTime(endTime string) {
	ss.m_callEndTime = endTime
}

func (ss *Session) setUuid(uuid string) {
	if len(ss.m_callerUuid) == 0 {
		ss.m_callerUuid = uuid
	} else {
		ss.m_calledUuid = uuid
	}
}

func (ss *Session) setSsrc(ssrc string) {
	if len(ss.m_callerSsrc) == 0 {
		ss.m_callerSsrc = ssrc
	} else {
		ss.m_calledSsrc = ssrc
	}
}

func (ss *Session) setCallRouteLog(log string) {
	ss.m_callRouteLog.WriteString(log)
}

func (ss *Session) setMcsLog(log string) {
	ss.m_mcsLog.WriteString(log)
}

func (ss *Session) setMasLog(log string) {
	ss.m_masLog.WriteString(log)
}

func (ss *Session) setSipmgwLog(log string) {
	ss.m_sipmgwLog.WriteString(log)
}

func (ss *Session) setSipsgwLog(log string) {
	ss.m_sipsgwLog.WriteString(log)
}

func (ss *Session) format() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s",
		ss.m_caller, ss.m_called,
		ss.m_callId,
		ss.m_callerUuid, ss.m_calledUuid,
		ss.m_callerSsrc, ss.m_calledSsrc,
		ss.m_callBeginTime, ss.m_callEndTime)
}

func (ss *Session) parse(line string) bool {
	strArray := strings.Split(line, ",")
	if len(strArray) < 9 {
		return false
	}
	ss.m_caller = strArray[0]
	ss.m_called = strArray[1]
	ss.m_callId = strArray[2]
	ss.m_callerUuid = strArray[3]
	ss.m_calledUuid = strArray[4]
	ss.m_callerSsrc = strArray[5]
	ss.m_calledSsrc = strArray[6]
	ss.m_callBeginTime = strArray[7]
	ss.m_callEndTime = strArray[8]
	return true
}
