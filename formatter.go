package logger

type Formatter interface {
	FormatMsg(msg *LogMsg) []byte
}
