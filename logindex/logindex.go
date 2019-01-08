package logindex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var indexDIr = "../data/index/"

type IndexInfo struct {
	beginTime string
	logData   []string
}

type LogIndex struct {
}

func (li *LogIndex) saveIndex(callLogMap map[string]*IndexInfo, modeType string) {
	for key, info := range callLogMap {
		var buffer bytes.Buffer
		buffer.WriteString(indexDIr)
		buffer.WriteString(modeType)
		buffer.WriteString("/")
		buffer.WriteString(info.beginTime)
		dir := buffer.String()
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			continue
		}
		buffer.WriteString("/")
		buffer.WriteString(key)
		path := strings.Replace(buffer.String(), ":", "", -1) // for windows,
		fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Open callRouteLogMap File faild:%s\n", path)
			continue
		}
		for _, line := range info.logData {
			n, err := fd.WriteString(line)
			if err != nil {
				fmt.Println("write File faild:", err)
			}
			fmt.Println("write File number:", n)
		}
		fd.Close()
	}
}

func (li *LogIndex) BuildIndex(data []string, modeName string,
	parseKeyCallBack func(line string) string,
	parseTimeCallBack func(line string) string) {

	var key string
	var info *IndexInfo
	callLogMap := make(map[string]*IndexInfo)
	for _, line := range data {
		key = parseKeyCallBack(line)
		if len(key) > 0 {
			_, exsit := callLogMap[key]
			if exsit {
				info = callLogMap[key]
			} else {
				info = new(IndexInfo)
				callLogMap[key] = info
				info.beginTime = parseTimeCallBack(line)
			}
			info.logData = append(info.logData, line)
		}
	}
	li.saveIndex(callLogMap, modeName)
}

func (li *LogIndex) FindData(key string, beginTime string, modeName string) string {
	if len(beginTime) < 14 {
		return ""
	}
	var index string
	var indexDir bytes.Buffer
	indexDir.WriteString(indexDIr)
	//indexDir.WriteString("/")
	indexDir.WriteString(modeName)
	indexDir.WriteString("/")
	indexDir.WriteString(beginTime[:8])
	indexDir.WriteString("/")
	indexDir.WriteString(beginTime[8:10])
	indexDir.WriteString("/")
	indexFileArray, err := ioutil.ReadDir(indexDir.String())
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, indexFile := range indexFileArray {
		if strings.Index(indexFile.Name(), key) >= 0 {
			index = indexDir.String() + indexFile.Name()
			data, err := ioutil.ReadFile(index)
			if err != nil {
				fmt.Printf("FindIndexLog:%s\n", err)
				return ""
			}
			return string(data)
		}
	}
	return ""
}
