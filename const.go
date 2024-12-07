package logger

const (
	LogMsgBufferSize    int = 1024
	LogDefaultCallDepth int = 2
)

type LogLevel int

// The levels of logs.
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var LogLevelNames = []string{
	"DEBUG", "INFO", "WARN", "ERROR", "FATAL",
}

func (l LogLevel) ToString() string {
	if int(l) > len(LogLevelNames)-1 || l < 0 {
		return "UNKNOWN"
	}
	return LogLevelNames[l]
}
