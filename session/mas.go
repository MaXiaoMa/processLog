package session

import (
	"bytes"
	"logAnalysis/conf"
	"logAnalysis/logindex"
	"logAnalysis/monitor"
	"strings"
)

type MasManager struct {
	m_index logindex.LogIndex
	m_alert Alert
}

func (mas *MasManager) handleData(data []string, filePath string) {
	mas.m_index.BuildIndex(data, "mas", mas.parseMasSsrc, mas.parseMasTime)
	mas.m_alert.ParseLineAndAlert(data)
}

//对不同的日志格式，使用不同的处理方法
func (mas *MasManager) handleData2(data []string, filePath string) {
	mas.m_index.BuildIndex(data, "mas", mas.parseMasSsrc2, mas.parseMasTime2)
	mas.m_alert.ParseLineAndAlert(data)
}

func (mas *MasManager) Start(dir string) {
	var MasMonitor monitor.FileMonitor
	logVersion, _ := conf.GetServerLogVersion("mas")
	//根据日志格式的不同，使用不同的处理方法
	switch logVersion {
	case 1:
		MasMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, mas.handleData)
	case 2:
		MasMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, mas.handleData2)
	default:
		MasMonitor.StartMonitor(conf.Cf.StateFileDIR, dir, mas.handleData)
	}
}

func (mas *MasManager) parseMasSsrc(line string) string {
	//2018/8/7, 17:02:02.1636962  [SSRC:10223744] mediaType:1,
	index := strings.Index(line, "SSRC:")
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
func (mas *MasManager) parseMasTime(line string) string {
	//2018/8/7, 17:02:02.1636962  [SSRC:10223744] mediaType:1,
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
func (mas *MasManager) parseMasSsrc2(line string) string {
	//2018/8/7, 17:02:02.1636962  [SSRC:10223744] mediaType:1,
	index := strings.Index(line, "SSRC:")
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
func (mas *MasManager) parseMasTime2(line string) string {
	//2018/10/24 17:51:42 [E] TN:PktExtRecv0  1111111111[SSRC:4194560] FastQueue Buffer
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
