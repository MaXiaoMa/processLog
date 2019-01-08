package session

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

func parseDir(dir string) []string {
	var dirArray []string
	modeDirArray, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Open Dir Fail:%s\n", err)
	}
	for _, modeDir := range modeDirArray {
		dataDirArray, err := ioutil.ReadDir(dir + modeDir.Name())
		if err != nil {
			fmt.Printf("Open dataDirArray Fail:%s\n", err)
			continue
		}
		for _, dataDir := range dataDirArray {
			dirArray = append(dirArray, dir+modeDir.Name()+"/"+dataDir.Name())
		}
	}
	return dirArray
}

//only year,mounth,day,hour
func parseCallRouteTime(line string) string {
	//[2018-05-30 09:38:06]
	var dateBuffer bytes.Buffer
	index := strings.Index(line, "[")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, "-")
	if index < 0 {
		return ""
	}
	dateBuffer.WriteString(line[:index])

	line = line[index+1:]
	index = strings.Index(line, "-")
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

	line = line[index+1:]
	index = strings.Index(line, ":")
	if index < 0 {
		return ""
	}
	dateBuffer.WriteString("/")
	dateBuffer.WriteString(line[:index])
	return dateBuffer.String()
}

func parseTime(line string) string {
	//[2018-05-30 09:38:06]
	var formtTime bytes.Buffer
	index := strings.Index(line, "[")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, "-")
	if index < 0 {
		return ""
	}
	year := line[:index]

	line = line[index+1:]
	index = strings.Index(line, "-")
	if index < 0 {
		return ""
	}
	mounth := line[:index]

	line = line[index+1:]
	index = strings.Index(line, " ")
	if index < 0 {
		return ""
	}
	day := line[:index]

	line = line[index+1:]
	index = strings.Index(line, ":")
	if index < 0 {
		return ""
	}
	hour := line[:index]

	line = line[index+1:]
	index = strings.Index(line, ":")
	if index < 0 {
		return ""
	}
	min := line[:index]

	line = line[index+1:]
	index = strings.Index(line, "]")
	if index < 0 {
		return ""
	}
	second := line[:index]

	formtTime.WriteString(year)
	formtTime.WriteString(mounth)
	formtTime.WriteString(day)
	formtTime.WriteString(hour)
	formtTime.WriteString(min)
	formtTime.WriteString(second)
	return formtTime.String()
}

func parseCallId(line string) string {
	//[2018-05-30 09:38:06] :[EC109380] <-- CallMsg_ProtoBuf_Invite
	index := strings.Index(line, "[")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, "]")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, "[")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, " ")
	if index < 0 {
		return ""
	}
	return line[:index]
}

func parseSipCaller(line string) string {
	index := strings.Index(line, "msg.m_Caller=")
	if index < 0 {
		return ""
	}
	line = line[index+13:]
	index = strings.Index(line, ",")
	if index < 0 {
		return ""
	}
	return line[:index]
}

func parseSipCalled(line string) string {
	index := strings.Index(line, "msg.m_Called=")
	if index < 0 {
		return ""
	}
	line = line[index+13:]
	index = strings.Index(line, ",")
	if index < 0 {
		return ""
	}
	return line[:index]
}

func parsePbCaller(line string) string {
	index := strings.Index(line, "pCallEvent->caller()=")
	if index < 0 {
		return ""
	}
	line = line[index+21:]
	index = strings.Index(line, ",")
	if index < 0 {
		return ""
	}
	return line[:index]
}

func parsePbCalled(line string) string {
	index := strings.Index(line, "pCallEvent->called()=")
	if index < 0 {
		return ""
	}
	line = line[index+21:]
	index = strings.Index(line, ",")
	if index < 0 {
		return ""
	}
	return line[:index]
}

func parseUuid(line string) string {
	//"uuid":"215a0da3-1bf8-4a4d-a2df-4744c11f92e8",...
	index := strings.Index(line, "\"uuid\"")
	if index < 0 {
		return ""
	}
	line = line[index+6:]
	index = strings.Index(line, "\"")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, "\"")
	if index < 0 {
		return ""
	}
	return line[:index]
}

func paseSsrc(line string) string {
	//"ssrc":"3408000",
	index := strings.Index(line, "\"ssrc\"")
	if index < 0 {
		return ""
	}
	line = line[index+6:]
	index = strings.Index(line, "\"")
	if index < 0 {
		return ""
	}
	line = line[index+1:]
	index = strings.Index(line, "\"")
	if index < 0 {
		return ""
	}
	return line[:index]
}
