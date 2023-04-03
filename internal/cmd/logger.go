package cmd

type loggerI interface {
	Infof(format string, args ...any)
	Infoln(args ...any)
	Warnf(format string, args ...any)
	Warnln(args ...any)
	Errorf(format string, args ...any)
	Errorln(args ...any)
	Debugf(format string, args ...any)
	Debugln(args ...any)
	PrintJSON(msg string, v any)
}
