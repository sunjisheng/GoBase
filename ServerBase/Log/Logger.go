package Log

import (
	"../Container"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Color struct {
	black        int // 黑色
	blue         int // 蓝色
	green        int // 绿色
	cyan         int // 青色
	red          int // 红色
	purple       int // 紫色
	yellow       int // 黄色
	light_gray   int // 淡灰色（系统默认值）
	gray         int // 灰色
	light_blue   int // 亮蓝色
	light_green  int // 亮绿色
	light_cyan   int // 亮青色
	light_red    int // 亮红色
	light_purple int // 亮紫色
	light_yellow int // 亮黄色
	white        int // 白色
}

// 给字体颜色对象赋值
var FontColor Color = Color{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

const (
	Log_Level_Debug = 0
	Log_Level_Info = 1
	Log_Level_Error = 2
	Log_Level_Count = 3
)

type LogFileInfo struct {
	filename string
	serialNo int
	file* os.File
}

type Logger struct {
	stopChan chan int								//关闭Chan
	writeChan chan int								//写Chan
	minLevel int									//日志等级
	curColor int                                    //当前颜色
	lastDay int                         			//最后日期
	logName string									//日志名
	currentPath string 			        			//当前路径
	buffers [Log_Level_Count]*Container.CircleBuf	//缓冲区
	mutex sync.Mutex
	fileInfos [Log_Level_Count]*LogFileInfo			//文件信息
	dirtyFlags [Log_Level_Count]bool            	//文件胀标记
}

var instance *Logger

func init() {
	instance = new(Logger)
	for level := 0; level < Log_Level_Count; level++ {
		instance.buffers[level]	= Container.NewCircleQueue(65536)
		instance.fileInfos[level] = new(LogFileInfo)
		instance.fileInfos[level].serialNo = 1
	}
	instance.stopChan = make(chan int)
	instance.writeChan = make(chan int ,1000)
	instance.curColor = FontColor.white
}

func Instance() *Logger {
	return instance
}

func StartLog(logName string, level int) {
	instance.logName = logName
	instance.minLevel = level
	instance.currentPath = GetCurrentPath()
	go instance.Loop()
}

func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

func StopLog() {
	instance.stopChan <- 1
}



func WriteLog(level int, format string, args ...interface{}) {
	now := time.Now()
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[%04d-%02d-%02d %02d:%02d:%02d] ", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	fmt.Fprintf(&buf, format, args...)
	buf.WriteByte('\n')
	instance.WriteLog(level, buf.Bytes(), uint32(buf.Len()))

	if runtime.GOOS == "windows" {
		if level == Log_Level_Error {
			PrintColorText(buf.String(), FontColor.red)
		} else {
			print(buf.String())
		}
	}
}

func (this *Logger) WriteLog(level int, buf []byte, len uint32) {
	if level < this.minLevel || level >= Log_Level_Count {
		return
	}
	this.mutex.Lock()
	this.buffers[level].Write(buf, len)
	this.mutex.Unlock()
	this.writeChan<-level
}

func (this *Logger) Loop() {
	tiker := time.NewTicker(time.Second * 5)
	var stop bool = false
	for !stop {
		select {
		case level := <-this.writeChan:
			this.WriteFile(level)
			this.FlushFile(level)
		case <- tiker.C:
			this.CheckDayChange()
			for level := 0; level < Log_Level_Count; level++ {
				this.CheckFileHuge(level)
			}
		case <-this.stopChan:
			stop = true
			break
		}
	}
	for level := 0; level < Log_Level_Count; level++ {
		this.FlushFile(level)
	}
}

func (this *Logger) CheckDayChange() bool {
	now := time.Now()
	curDay := now.Day()
	if this.lastDay == curDay {
		return false
	}
	for level := 0; level < Log_Level_Count; level++ {
		this.fileInfos[level].serialNo = 1
		if this.fileInfos[level].file != nil {
			this.fileInfos[level].file.Close()
			this.fileInfos[level].file = nil
		}
	}
	this.lastDay = curDay
	return true
}

func (this *Logger) CheckFileHuge(level int) {

	if this.fileInfos[level].file == nil {
		return
	}
	fi,err := this.fileInfos[level].file.Stat()
	if err == nil {
		if fi.Size() >= 10485760 {  //10M
			this.fileInfos[level].serialNo++
			this.fileInfos[level].file.Close()
			this.fileInfos[level].file = nil
		}
	}
}

func (this *Logger) CreateFile(level int) {
	var levelNames[Log_Level_Count]string = [Log_Level_Count]string{"Debug","Info","Error"}
	now := time.Now()
	fileName := fmt.Sprintf("%s/%s_%s_%04d%02d%02d_%d.log", this.currentPath, this.logName, levelNames[level], now.Year(), now.Month(), now.Day(), this.fileInfos[level].serialNo)
	if this.fileInfos[level].file != nil	{
		this.fileInfos[level].file.Close()
	}
	this.fileInfos[level].file, _ = os.OpenFile(fileName, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
}

func (this *Logger) WriteFile(level int) {
	if this.fileInfos[level].file == nil {
		this.CreateFile(level)
	}
	for !this.buffers[level].IsEmpty()	{
		slice := this.buffers[level].GetReadSlice()
		this.fileInfos[level].file.Write(slice)
		this.buffers[level].Skip(uint32(len(slice)))
		this.dirtyFlags[level] = true
	}
}

func (this *Logger) FlushFile(level int) {
	if !this.dirtyFlags[level]	{
		return
	}
	this.dirtyFlags[level] = false
	if this.fileInfos[level].file != nil {
		this.fileInfos[level].file.Sync()
	}
}
