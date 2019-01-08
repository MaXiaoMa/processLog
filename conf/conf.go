package conf

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

//全局配置变量
var Cf *Conf

type alarm struct {
	SendTo   string `json:"sendTo"`
	Port     string `json:"port"`
	Oid      string `json:"oid"`
	Rtt      int    `json:"rtt"`
	LostRate int    `json:"lostRate"`
}

type log struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	LogVersion int    `json:"logVersion"`
}
type Conf struct {
	ListenPort   string `json:"listenPort"`
	LogPath      string `json:"logpath"`
	FileName     string `json:"filename"`
	Alarm        alarm  `json:"alarm"`
	StateFileDIR string `json:"stateFileDir"`
	LogFile      []log  `json:"logfile"`
}

func newConfig() *Conf {
	return &Conf{}
}

//过滤json文件中的注释
func filterComment(jsonStr string) (string, error) {
	jsonBuf := strings.NewReader(jsonStr)
	scanner := bufio.NewScanner(jsonBuf)
	var res string
	for scanner.Scan() {
		piece := scanner.Text()
		raw := strings.TrimLeft(piece, " \t")
		if !strings.HasPrefix(raw, "//") {
			//处理行内的//注释
			rawArr := strings.Split(raw, "//")
			res += rawArr[0]
		}
	}
	return res, nil
}

//从文件读取，初始化全局配置
func Init(path string) (err error) {
	Cf = newConfig()
	configStr, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("open file failed: ", err)
		return
	}
	configJson, err := filterComment(string(configStr))
	if err != nil {
		fmt.Println("filter comment error: ", err)
	}
	err = json.Unmarshal([]byte(configJson), Cf)
	if err != nil {
		fmt.Println("unmarshal json failed : ", err)
	}
	return nil
}

//获取指定服务的日志路径
func GetServerPath(serverName string) (path string, err error) {
	for _, v := range Cf.LogFile {
		if v.Name == serverName {
			path = v.Path
		}
	}
	if path == "" {
		err = errors.New("can not find path")
		return
	}
	return
}

//获取指定服务的日志版本
func GetServerLogVersion(serverName string) (version int, err error) {
	for _, v := range Cf.LogFile {
		if v.Name == serverName {
			version = v.LogVersion
		}
	}
	if version == 0 {
		err = errors.New("can not find log version,use default 1")
		version = 1
		return
	}
	return
}

//获取所有服务的日志路径
func GetAllServerPaths() (path []string) {
	for _, v := range Cf.LogFile {
		path = append(path, v.Path)
	}
	return
}
