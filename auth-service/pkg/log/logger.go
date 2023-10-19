package log

import (
	"os"

	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/apperrors"
	"github.com/rs/zerolog"
)

var (
	zlog zerolog.Logger
	logW *os.File
)

func init() {
	initLogger()
}

type multiWritter struct {
	lw    zerolog.LevelWriter
	level zerolog.Level
}

func (m multiWritter) Write(p []byte) (n int, err error) {
	return m.lw.Write(p)
}

func (m multiWritter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {

	if level == m.level {
		return m.lw.WriteLevel(level, p)
	}
	return len(p), nil
}

func initLogger() {
	// success and fail log are preferred to be wrote to different file to make tracing is easier
	// os.stdout/os.stderr is preferred to accept all level of logs
	// hence 3 writer is used (2 file and 1 console writer)

	// fInfo is used for info level log
	fInfo, err := os.OpenFile("./etc/log/success.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	// fFail is for warn, error and fatal
	fFail, err := os.OpenFile("./etc/log/fail.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	writterSuccess := zerolog.MultiLevelWriter(fInfo)
	writterFail := zerolog.MultiLevelWriter(fFail)

	mwInfo := multiWritter{lw: writterSuccess, level: zerolog.InfoLevel}
	mwError := multiWritter{lw: writterFail, level: zerolog.ErrorLevel}
	mwWarn := multiWritter{lw: writterFail, level: zerolog.WarnLevel}
	mwFatal := multiWritter{lw: writterFail, level: zerolog.FatalLevel}
	consoleWritter := zerolog.ConsoleWriter{Out: os.Stdout}
	logW = os.Stdout

	zlog = zerolog.New(zerolog.MultiLevelWriter(mwInfo, mwError, mwWarn, mwFatal, consoleWritter)).With().Str("service_name", "auth-service").Timestamp().Logger()

}
func msg(le *zerolog.Event, message string, args ...interface{}) {
	if len(args) == 0 {
		le.Msg(message)
	} else {
		le.Msgf(message, args...)
	}
}

func Debug(message string, args ...interface{}) {
	le := zlog.Debug()
	msg(le, message, args...)
}

func Info(message string, args ...interface{}) {
	le := zlog.Info()
	msg(le, message, args...)
}

func Warn(message string, args ...interface{}) {
	le := zlog.Warn()
	msg(le, message, args...)
}

func ErrorWithCause(err, errCause error, message string, args ...interface{}) {
	le := zlog.Error().Err(err).AnErr("cause", errCause)
	msg(le, message, args...)
}

func Error(err error, message string, args ...interface{}) {
	le := zlog.Error().Err(err)
	msg(le, message, args...)
	if errStack, ok := err.(apperrors.WrappedError); ok && zerolog.GlobalLevel() == zerolog.DebugLevel {
		errStack.PrintStack(logW)
	}
}

func Fatal(err error, message string, args ...interface{}) {
	le := zlog.Fatal().Err(err)
	msg(le, message, args...)
}

func Log() zerolog.Logger {
	return zlog
}
