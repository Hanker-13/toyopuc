package log

import (
	"strings"

	"github.com/astaxie/beego/logs"
)

// Log levels to control the logging output.
// Deprecated: use github.com/astaxie/beego/logs instead.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

func init() {
	SetLogFuncCall(true)
}

// BeeLogger references the used application logger.
// Deprecated: use github.com/astaxie/beego/logs instead.
var BeeLogger = logs.GetBeeLogger()

func EnableFuncCallDepth(b bool) {
	logs.EnableFuncCallDepth(b)
}

func SetLogFuncCallDepth(d int) {
	logs.SetLogFuncCallDepth(d)

}

// SetLevel sets the global log level used by the simple logger.
// Deprecated: use github.com/astaxie/beego/logs instead.
func SetLevel(l int) {
	logs.SetLevel(l)
}

// SetLogFuncCall set the CallDepth, default is 3
// Deprecated: use github.com/astaxie/beego/logs instead.
func SetLogFuncCall(b bool) {
	logs.SetLogFuncCall(b)
}

// SetLogger sets a new logger.
// Deprecated: use github.com/astaxie/beego/logs instead.
func SetLogger(adaptername string, config string) error {
	return logs.SetLogger(adaptername, config)
}

// Emergency logs a message at emergency level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Emergency(v ...interface{}) {
	logs.Emergency(generateFmtStr(len(v)), v...)
}

// Alert logs a message at alert level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Alert(v ...interface{}) {
	logs.Alert(generateFmtStr(len(v)), v...)
}

// Critical logs a message at critical level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Critical(v ...interface{}) {
	logs.Critical(generateFmtStr(len(v)), v...)
}

// Error logs a message at error level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Error(v ...interface{}) {
	logs.Error(generateFmtStr(len(v)), v...)
}

// Warning logs a message at warning level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Warning(v ...interface{}) {
	logs.Warning(generateFmtStr(len(v)), v...)
}

// Warn compatibility alias for Warning()
// Deprecated: use github.com/astaxie/beego/logs instead.
func Warn(v ...interface{}) {
	logs.Warn(generateFmtStr(len(v)), v...)
}

// Notice logs a message at notice level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Notice(v ...interface{}) {
	logs.Notice(generateFmtStr(len(v)), v...)
}

// Informational logs a message at info level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Informational(v ...interface{}) {
	logs.Informational(generateFmtStr(len(v)), v...)
}

// Info compatibility alias for Warning()
// Deprecated: use github.com/astaxie/beego/logs instead.
func Info(v ...interface{}) {
	logs.Info(generateFmtStr(len(v)), v...)
}

// Debug logs a message at debug level.
// Deprecated: use github.com/astaxie/beego/logs instead.
func Debug(v ...interface{}) {
	logs.Debug(generateFmtStr(len(v)), v...)
}

// Trace logs a message at trace level.
// compatibility alias for Warning()
// Deprecated: use github.com/astaxie/beego/logs instead.
func Trace(v ...interface{}) {
	logs.Trace(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
