package logger

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type LogMsg struct {
	Level     LogLevel
	Timestamp time.Time
	Msg       string
	FileName  string
	FileLine  int
}

type Logger struct {
	level     LogLevel
	callDepth int
	signal    chan int
	msgBuffer chan *LogMsg
	msgPool   *sync.Pool
	wg        sync.WaitGroup
	executors []Executor
	formatter Formatter
	isRunning int32
}

func NewLogger() *Logger {
	return &Logger{
		level:     LevelInfo,
		signal:    make(chan int),
		msgBuffer: make(chan *LogMsg, LogMsgBufferSize),
		msgPool: &sync.Pool{New: func() interface{} {
			return &LogMsg{}
		}},
		isRunning: 0,
	}
}

func NewDefaultLogger() *Logger {
	logger := NewLogger()
	logger.AppendExecutor(NewConsoleExecutor())
	return logger
}

func (l *Logger) AppendExecutor(executor Executor) {
	l.executors = append(l.executors, executor)
}

func (l *Logger) ResetExecutor() {
	l.executors = make([]Executor, 0)
}

func (l *Logger) SetFormatter(formatter Formatter) {
	l.formatter = formatter
}

func (l *Logger) GetLevel() LogLevel {
	return l.level
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) GetCallDepth() int {
	return l.callDepth
}

func (l *Logger) SetCallDepth(callDepth int) {
	l.callDepth = callDepth
}

func (l *Logger) write(msg *LogMsg) {
	if len(l.executors) == 0 {
		return
	}
	if l.formatter == nil {
		l.formatter = &DefaultFormatter{}
	}

	for _, executorItem := range l.executors {
		if err := executorItem.WriteMsg(l.formatter.FormatMsg(msg)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to save log to destination, err:%v\n", err)
		}
	}
}

func (l *Logger) WriteMsg(msg string, level LogLevel) {
	if atomic.LoadInt32(&l.isRunning) == 0 {
		return
	}

	var fileName string
	var fileLine int
	if l.callDepth > 0 {
		_, file, line, ok := runtime.Caller(l.callDepth)
		if ok {
			fileName = file
			fileLine = line
		}
	}

	logMsg := l.msgPool.Get().(*LogMsg)
	logMsg.Level = level
	logMsg.Msg = msg
	logMsg.Timestamp = time.Now()
	logMsg.FileName = fileName
	logMsg.FileLine = fileLine
	l.msgBuffer <- logMsg
}

func (l *Logger) Close() {
	if !atomic.CompareAndSwapInt32(&l.isRunning, 1, 0) {
		return
	}

	l.signal <- 1
	l.wg.Wait()
	close(l.msgBuffer)
	close(l.signal)
	// close all
	for _, executorItem := range l.executors {
		executorItem.Close()
	}
}

func (l *Logger) Flush() {
	for {
		if len(l.msgBuffer) == 0 {
			break
		}

		msg, ok := <-l.msgBuffer
		if !ok {
			break
		}
		l.write(msg)
	}

	for _, executorItem := range l.executors {
		executorItem.Flush()
	}
}

func (l *Logger) Start() {
	if !atomic.CompareAndSwapInt32(&l.isRunning, 0, 1) {
		return
	}

	workerFunc := func() {
		defer l.wg.Done()
		for {
			select {
			case msg := <-l.msgBuffer:
				l.write(msg)
				l.msgPool.Put(msg)
			case <-l.signal:
				l.Flush()
				return
			}
		}
	}

	l.wg.Add(1)
	go workerFunc()
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.GetLevel() > LevelDebug {
		return
	}
	l.WriteMsg(fmt.Sprintf(format, v...), LevelDebug)
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.GetLevel() > LevelInfo {
		return
	}
	l.WriteMsg(fmt.Sprintf(format, v...), LevelInfo)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.GetLevel() > LevelWarn {
		return
	}
	l.WriteMsg(fmt.Sprintf(format, v...), LevelWarn)
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.GetLevel() > LevelError {
		return
	}
	l.WriteMsg(fmt.Sprintf(format, v...), LevelError)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	if l.GetLevel() > LevelFatal {
		return
	}
	l.WriteMsg(fmt.Sprintf(format, v...), LevelFatal)
}
