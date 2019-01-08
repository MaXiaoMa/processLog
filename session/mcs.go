package session

import (
	"bytes"
	"logAnalysis/conf"
	"logAnalysis/logindex"
	"logAnalysis/monitor"
	"strings"
)

type McsManager struct {
	m_index logindex.LogIndex
}

func (mcs *McsManager) handleData(data []string, filePath string) {
	mcs.m_index.BuildIndex(data, "mcs", mcs.parseMcsUuid, mcs.parseMcsTime)
}

func (mcs *McsManager) Start(dir string) {
	var McsMonitor monitor.FileMonitor
	logVersion, _ := conf.GetServerLogVersion("mcs")
	//根据日志格式的不同，使用不同的处理方法
	switch logVersion {
	case 1:
		McsMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, mcs.handleData)
	default:
		McsMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, mcs.handleData)
	}
}

func (mcs *McsManager) parseMcsUuid(line string) string {
	//[confid:p3435026975] [uuid:45d9af2c-8547-4742-a990-1a5f173f73b4] conference
	index := strings.Index(line, "uuid:")
	if index < 0 {
		return ""
	}
	line = line[index+5:]
	indexSpace := strings.Index(line, " ")
	indexBracket := strings.Index(line, "]")
	//都没找到， 返回空
	if indexSpace <= 0 && indexBracket <= 0 {
		return ""
	}
	//都找到，则取最小的
	if indexSpace > 0 && indexBracket > 0 {
		if indexSpace < indexBracket {
			index = indexSpace
		} else {
			index = indexBracket
		}
	} else {
		//找到一个，则返回找到的index

		if indexBracket > 0 {
			index = indexBracket
		}
		if indexSpace > 0 {
			index = indexSpace
		}
	}
	return line[:index]
}

//only year,mounth,day,hour
func (mcs *McsManager) parseMcsTime(line string) string {
	//2018/08/09 15:47:47 [debug]
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
	//dateBuffer.WriteString("-")

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
