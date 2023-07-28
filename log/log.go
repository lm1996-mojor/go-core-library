package log

import (
	"fmt"
	"os"

	"mojor/go-core-library/global"

	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const runLevel = -5

func Init(app *iris.Application) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: "2006/01/02 15:04:05"}
	log.Logger = log.Output(output)
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
	log.Warn().Msg(msg)
}

// Error log in Error level
func Error(msg string) {
	log.Error().Msg(msg)
}

// Errorf log in INFO level
func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Error().Msg(msg)
}

// Fatal log in Error level
func Fatal(msg string) {
	log.Fatal().Msg(msg)
}

// Panic log in PANIC level
func Panic(msg string) {
	log.Panic().Msg(msg)
}

// Err log an error and a message
func Err(err error, msg string) {
	log.Err(err).Msg(msg)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
