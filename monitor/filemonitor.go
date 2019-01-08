package monitor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const POLLING_INTERVALl = 5000

type FileInfo struct {
	dir      string
	fileinfo os.FileInfo
}
type FileMonitor struct {
	//read file data buf,avoid incomplete data
	m_buffer bytes.Buffer
	// state file dir
	m_stateFileDir string
	//file state map,key=state file path, value=state struct object
	m_stateMap map[string]*FileState
	//func callback for file update
	m_fileUpdateCallBack func(data []string, filePath string)
}

//statefile: file-state:path-fileName
func (monitor *FileMonitor) fileStateNameBulid(path string) string {

	fileStateName := "file-state"
	strArray := strings.Split(path, "/")
	for _, str := range strArray {
		fileStateName += "-"
		fileStateName += str
	}
	return fileStateName
}

//Check file status if file is updated return true
func (monitor *FileMonitor) fileStateCheck(srcFilePath string, fileInfo os.FileInfo) bool {

	var state *FileState
	stateFilePath := monitor.m_stateFileDir + monitor.fileStateNameBulid(srcFilePath)
	_, exsit := monitor.m_stateMap[stateFilePath]
	if exsit {
		state = monitor.m_stateMap[stateFilePath]
	} else {
		state = new(FileState)
		state.setPath(srcFilePath)
		monitor.m_stateMap[stateFilePath] = state
	}
	fileModTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")
	if strings.Compare(fileModTime, state.getReadTime()) > 0 {
		return true
	}
	return false
}

//on init, read the file state,avoid repeated reading data
func (monitor *FileMonitor) fileStateRead() {

	fileInfoArray, _ := ioutil.ReadDir(monitor.m_stateFileDir)
	for _, fileInfo := range fileInfoArray {
		if fileInfo.IsDir() == true {
			//fmt.Printf("This is a Dir:%s\n", path.Name())
			continue
		}
		stateFilePath := monitor.m_stateFileDir + fileInfo.Name()
		buf, err := ioutil.ReadFile(stateFilePath)
		if err != nil {
			continue
		}
		state := new(FileState)
		if state.parse(string(buf)) == true {
			monitor.m_stateMap[stateFilePath] = state
		}
	}
}

//Writing file status
func (monitor *FileMonitor) fileStateWrite(stateFilePath string, stateStr string) {

	err := os.MkdirAll(monitor.m_stateFileDir, os.ModePerm)
	if err != nil {
		//fmt.Println(err)
		return
	}
	fd, err := os.OpenFile(stateFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Open State File faild:%s\n", stateFilePath)
	}
	fd.WriteString(stateStr)
	fd.Close()
}

//Read updated file data
func (monitor *FileMonitor) srcFileRead(srcFilePath string, fileInfo os.FileInfo) {

	stateFilePath := monitor.m_stateFileDir + monitor.fileStateNameBulid(srcFilePath)
	_, exsit := monitor.m_stateMap[stateFilePath]
	if exsit {
		state := monitor.m_stateMap[stateFilePath]
		state.setReadTime(fileInfo.ModTime().Format("2006-01-02 15:04:05"))

		fd, err := os.Open(srcFilePath)
		if err != nil {
			fmt.Printf("Open srcFilePath Fail:%s,%s\n", srcFilePath, err)
			return //
		}
		for {
			buf := make([]byte, 8192)
			count, _ := fd.ReadAt(buf, state.getReadPosotion())
			if count <= 0 {
				break
			}
			state.addReadPosotion(int64(count))
			monitor.m_buffer.Write(buf[:count])
		}
		index := strings.LastIndex(monitor.m_buffer.String(), "\n")
		currStrings := monitor.m_buffer.String()[:index+1]
		remainStrings := monitor.m_buffer.String()[index+1:]
		monitor.m_buffer.Reset()
		monitor.m_buffer.WriteString(remainStrings)
		if monitor.m_fileUpdateCallBack != nil {
			monitor.m_fileUpdateCallBack(strings.SplitAfter(currStrings, "\n"), srcFilePath)
		}
		monitor.fileStateWrite(stateFilePath, state.format())
	}
}

func (monitor *FileMonitor) StartMonitor(stateDir string, dir string, callBack func(data []string, filePath string)) {

	monitor.m_stateFileDir = stateDir //"./state/"
	monitor.m_fileUpdateCallBack = callBack
	monitor.m_stateMap = make(map[string]*FileState)
	monitor.fileStateRead()
	for {
		//找出路径下所有的文件
		fileInfoArray := getAllFiles(dir)

		for _, file := range fileInfoArray {
			srcFilePath := file.dir + file.fileinfo.Name()
			if monitor.fileStateCheck(srcFilePath, file.fileinfo) == true {
				monitor.srcFileRead(srcFilePath, file.fileinfo)
				//log
			}
		}

		time.Sleep(POLLING_INTERVALl * time.Millisecond)
	}
}

//迭代获取该目录下所有日志文件
func getAllFiles(dir string) (fileInfoArr []FileInfo) {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range items {
		if !item.IsDir() {
			fileInfo := FileInfo{
				dir:      dir,
				fileinfo: item,
			}
			fileInfoArr = append(fileInfoArr, fileInfo)
		} else {
			fileInfoArr = append(fileInfoArr, getAllFiles(dir+item.Name())...)
		}
	}
	return
}
