package logger

type Executor interface {
	WriteMsg(msg []byte) error
	Flush()
	Close()
}
