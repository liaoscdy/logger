package logger

import (
	"bytes"
	"strconv"
	"time"
)

type DefaultFormatter struct{}

func (s *DefaultFormatter) FormatMsg(msg *LogMsg) []byte {
	// eg: 2006-01-02 15:04:05 [INFO] [code.go:90] the log message
	byBuf := bytes.Buffer{}
	// datetime + level + other + msg
	byBuf.Grow(19 + 8 + 4 + len(msg.FileName) + len(msg.Msg))
	byBuf.WriteString(msg.Timestamp.Format(time.DateTime))
	byBuf.WriteString(" [")
	byBuf.WriteString(msg.Level.ToString())
	byBuf.WriteString("] ")
	if len(msg.FileName) > 0 {
		byBuf.WriteString("[")
		byBuf.WriteString(msg.FileName)
		byBuf.WriteString(":")
		byBuf.WriteString(strconv.Itoa(msg.FileLine))
		byBuf.WriteString("] ")
	}

	byBuf.WriteString(msg.Msg)
	return byBuf.Bytes()
}
