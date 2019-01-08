package session

import (
	"bytes"
	"logAnalysis/conf"
	"logAnalysis/logindex"
	"logAnalysis/monitor"
	"strings"
)

type SipmgwManager struct {
	m_index logindex.LogIndex
}

func (sipmgw *SipmgwManager) handleData(data []string, filePath string) {
	sipmgw.m_index.BuildIndex(data, "sipmgw", sipmgw.parseSipmgwSsrc, sipmgw.parseSipmgwTime)
}

func (sipmgw *SipmgwManager) Start(dir string) {
	var sipmgwMonitor monitor.FileMonitor
	logVersion, _ := conf.GetServerLogVersion("sipmgw")
	//根据日志格式的不同，使用不同的处理方法
	switch logVersion {
	case 1:
		sipmgwMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, sipmgw.handleData)
	default:
		sipmgwMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, sipmgw.handleData)
	}
}

func (sipmgw *SipmgwManager) parseSipmgwSsrc(line string) string {
	//2018/8/16, 15:45:37.1853274 [I] TN:Manager      [ssrc:7471360] ManagerTask
	index := strings.Index(line, "ssrc:")
	if index < 0 {
		return ""
	}
	line = line[index+5:]
	index = strings.Index(line, "]")
	if index < 0 {
		return ""
	}
	return line[:index]
}

//only year,mounth,day,hour
func (sipmgw *SipmgwManager) parseSipmgwTime(line string) string {
	//2018/8/16, 15:45:37.1853274 [I] TN:Manager      [ssrc:7471360] ManagerTask
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
	if len(line[:index]) <= 1 {
		dateBuffer.WriteString("0")
	}
	dateBuffer.WriteString(line[:index])
	//dateBuffer.WriteString("-")

	line = line[index+1:]
	index = strings.Index(line, ",")
	if index < 0 {
		return ""
	}
	if len(line[:index]) <= 1 {
		dateBuffer.WriteString("0")
	}
	dateBuffer.WriteString(line[:index])
	dateBuffer.WriteString("/")

	line = line[index+1:]
	index = strings.Index(line, " ")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, ":")
	if index < 0 {
		return ""
	}
	dateBuffer.WriteString(line[:index])
	return dateBuffer.String()
}
