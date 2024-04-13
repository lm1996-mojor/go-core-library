package log

import (
	"fmt"
	"os"
	"runtime"

	"github.com/lm1996-mojor/go-core-library/config"
	"github.com/lm1996-mojor/go-core-library/global"

	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const runLevel = -5

func Init(app *iris.Application) {
	log.Info().Msg("日志组件初始化....")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	if config.Sysconfig.SystemEnv.Env == "prod" {
		output.NoColor = true
	} else {
		output.NoColor = false
	}
	log.Logger = log.Output(output)
	log.Info().Msg("日志组件初始化完成")
}

func init() {
	global.RegisterInit(global.Initiator{Action: Init, Level: runLevel})
}

// Debug log in DEBUG level
func Debug(msg string) {
	log.Debug().Msg(msg)
}

// Info log in INFO level
func Info(msg string) {
	log.Info().Msg(msg)
}

// Infof log in INFO level
func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Info().Msg(msg)
}

// Warn log in WARN level
func Warn(msg string) {
	file, line, funcName := logMasterSource()
	msg += ", On line [" + fmt.Sprint(line) + "] of the [" + file + "],\n the function name[" + funcName + "]"
	log.Warn().Msg(msg)
}

func WarnF(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	file, line, funcName := logMasterSource()
	msg += ", On line [" + fmt.Sprint(line) + "] of the [" + file + "],\n the function name[" + funcName + "]"
	log.Warn().Msgf(format, v)
}

// Error log in Error level
func Error(msg string) {
	file, line, funcName := logMasterSource()
	msg += ", On line [" + fmt.Sprint(line) + "] of the [" + file + "],\n the function name[" + funcName + "]"
	log.Error().Msg(msg)
}

// Errorf log in INFO level
func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	file, line, funcName := logMasterSource()
	msg += ", On line [" + fmt.Sprint(line) + "] of the [" + file + "],\n the function name[" + funcName + "]"
	log.Error().Msg(msg)
}

// Fatal log in Error level
func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

// Panic log in PANIC level
func Panic(msg string) {
	file, line, funcName := logMasterSource()
	msg += ", On line [" + fmt.Sprint(line) + "] of the [" + file + "],\n the function name[" + funcName + "]"
	log.Panic().Msg(msg)
}

// Err log an error and a message
func Err(err error, msg string) {
	file, line, funcName := logMasterSource()
	msg += ", On line [" + fmt.Sprint(line) + "] of the [" + file + "],\n the function name[" + funcName + "]"
	log.Err(err).Msg(msg)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func logMasterSource() (string, int, string) {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return file, line, funcName
}
