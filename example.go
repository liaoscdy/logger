package logger

func ConsoleLoggerExample() {
	logger := NewDefaultLogger()
	logger.SetCallDepth(2)
	logger.SetLevel(LevelDebug)
	logger.Start()
	defer logger.Close()

	// ------------------ print log msg
	logger.Info("info Msg, value:%d", 1)
	logger.Error("error Msg without any format value")
	logger.Warn("warn Msg, Value1:%d, Value2:%s", 100, "stringByFmt")
}

func FileLoggerExample() {
	logger := NewLogger()
	logger.AppendExecutor(NewFileExecutor("/tmp/logger_test.log"))
	logger.SetCallDepth(2)
	logger.SetLevel(LevelDebug)
	logger.Start()
	defer logger.Close()

	// ------------------ print log msg
	logger.Info("info Msg, value:%d", 1)
	logger.Error("error Msg without any format value")
	logger.Warn("warn Msg, Value1:%d, Value2:%s", 100, "stringByFmt")
}

func BothFileAndConsoleExample() {
	logger := NewLogger()
	logger.AppendExecutor(NewConsoleExecutor())
	logger.AppendExecutor(NewFileExecutor("/tmp/logger_test.log"))
	logger.SetCallDepth(2)
	logger.SetLevel(LevelDebug)
	logger.Start()
	defer logger.Close()

	// ------------------ print log msg
	logger.Info("info Msg, value:%d", 1)
	logger.Error("error Msg without any format value")
	logger.Warn("warn Msg, Value1:%d, Value2:%s", 100, "stringByFmt")
}
