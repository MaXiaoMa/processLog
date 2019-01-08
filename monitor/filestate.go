package monitor

import (
	"fmt"
	"strconv"
	"strings"
)

type FileState struct {
	// monitor file path
	m_path string
	// monitor file last read time
	m_readTime string
	// monitor file last read pos
	m_readPosition int64
}

func (fs *FileState) setPath(path string) {
	fs.m_path = path
}

func (fs *FileState) setReadTime(readTime string) {
	fs.m_readTime = readTime
}

func (fs *FileState) setReadPosotion(pos int64) {
	fs.m_readPosition = pos
}

func (fs *FileState) addReadPosotion(pos int64) {
	fs.m_readPosition += pos
}

func (fs *FileState) getPath() string {
	return fs.m_path
}

func (fs *FileState) getReadTime() string {
	return fs.m_readTime
}

func (fs *FileState) getReadPosotion() int64 {
	return fs.m_readPosition
}

func (fs *FileState) format() string {
	return fmt.Sprintf("%s,%s,%d", fs.m_path, fs.m_readTime, fs.m_readPosition)
}

func (fs *FileState) parse(fsStr string) bool {
	strArray := strings.Split(fsStr, ",")
	if len(strArray) < 3 {
		return false
	}
	i, err := strconv.ParseInt(strArray[2], 10, 64)
	if err != nil {
		return false
	}
	fs.setPath(strArray[0])
	fs.setReadTime(strArray[1])
	fs.setReadPosotion(i)
	return true
}
