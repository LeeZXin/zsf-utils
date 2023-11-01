package sentinelutil

type EmptyLogger struct {
}

func (*EmptyLogger) Debug(msg string, keysAndValues ...interface{}) {}

func (*EmptyLogger) DebugEnabled() bool {
	return false
}

func (*EmptyLogger) Info(msg string, keysAndValues ...interface{}) {}

func (*EmptyLogger) InfoEnabled() bool {
	return false
}

func (*EmptyLogger) Warn(msg string, keysAndValues ...interface{}) {}

func (*EmptyLogger) WarnEnabled() bool {
	return false
}

func (*EmptyLogger) Error(err error, msg string, keysAndValues ...interface{}) {}

func (*EmptyLogger) ErrorEnabled() bool {
	return false
}
