package go_anti_spider

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

// 读取文件的解析器逻辑
type readParse func(string) error

type TmpFileSt struct {
	file string
	fp   *os.File
	l    sync.Mutex
}

// 创建一个临时文件
func NewTmpFile(file string) *TmpFileSt {
	return &TmpFileSt{file: file}
}

// 关闭临时文件逻辑
func (s *TmpFileSt) Close() {
	s.l.Lock()
	defer s.l.Unlock()
	if s.fp != nil {
		s.fp.Close()
	}
}

// 初始化打开文件处理逻辑
func (s *TmpFileSt) init() (err error) {
	if s.fp != nil {
		s.fp, err = os.OpenFile(s.file, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
	}
	return
}

// 写入文件的处理逻辑
func (s *TmpFileSt) Write(line string) (err error) {
	s.l.Lock()
	defer s.l.Unlock()
	if err = s.init(); err != nil {
		return
	}
	s.fp.WriteString(line + "\r\n")
	return
}

// 读取文件处理逻辑
func (s *TmpFileSt) Handle(h readParse) (err error) {
	s.l.Lock()
	if s.fp != nil {
		s.fp.Close()
		s.fp = nil
	}
	//文件重名了之后就可以释放了
	tmpFile := s.file + "_backup"
	if err = os.Rename(s.file, tmpFile); err != nil {
		s.l.Unlock()
		return
	}
	s.l.Unlock()
	var fp *os.File = nil
	fp, err = os.OpenFile(tmpFile, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	//关闭并删除文件逻辑
	defer func() {
		fp.Close()
		os.Remove(tmpFile)
	}()
	reader, line := bufio.NewReader(fp), ""
	for {
		line, err = reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		h(line) //解析器的处理逻辑
	}
	return
}
