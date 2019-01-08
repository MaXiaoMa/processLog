package session

import (
	"bytes"
	"logAnalysis/conf"
	"logAnalysis/logindex"
	"logAnalysis/monitor"
	"strings"
)

type SipsgwManager struct {
	m_index logindex.LogIndex
}

func (sipsgw *SipsgwManager) handleData(data []string, filePath string) {
	sipsgw.m_index.BuildIndex(data, "sipsgw", sipsgw.parseSipsgwCallid, sipsgw.parseSipsgwTime)
}

func (sipsgw *SipsgwManager) Start(dir string) {
	var sipsgwMonitor monitor.FileMonitor
	logVersion, _ := conf.GetServerLogVersion("sipsgw")
	//根据日志格式的不同，使用不同的处理方法
	switch logVersion {
	case 1:
		sipsgwMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, sipsgw.handleData)
	default:
		sipsgwMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, sipsgw.handleData)
	}
}

func (sipsgw *SipsgwManager) parseSipsgwCallid(line string) string {
	//2018/08/16 15:45:42.790450 [INFO ] [udp_revied.go:238 handPacket()]  [callId:EC4154526@00:FF:97:DA:72:68192.168.28.92] handle
	index := strings.Index(line, "callId:")
	if index < 0 {
		return ""
	}
	line = line[index+7:]
	index = strings.Index(line, "]")
	if index < 0 {
		return ""
	}
	return line[:index]
}

//only year,mounth,day,hour
func (sipsgw *SipsgwManager) parseSipsgwTime(line string) string {
	//2018/08/16 15:45:42.790450 [INFO ] [udp_revied.go:238 handPacket()]  [callId:EC4154526@00:FF:97:DA:72:68192.168.28.92] handle
	var dateBuffer bytes.Buffer
	index := strings.Index(line, "/")
	if index < 0 {
		return ""
	}
	dateBuffer.WriteString(line[:index])
	//dateBuffer.WriteString("-")

	line = line[index+1:]
	index = strings.Index(line, "/")
	if index < 0 {
		return ""
	}
	dateBuffer.WriteString(line[:index])

	line = line[index+1:]
	index = strings.Index(line, " ")
	if index < 0 {
		return ""
	}

	dateBuffer.WriteString(line[:index])
	dateBuffer.WriteString("/")

	line = line[index+1:]
	index = strings.Index(line, ":")
	if index < 0 {
		return ""
	}
	dateBuffer.WriteString(line[:index])
	return dateBuffer.String()
}
