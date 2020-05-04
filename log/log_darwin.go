package log

import (
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
	"syscall"
)

func InitLog(runMode string, dirs ...string) *logs.BeeLogger {
	DefaultLog = logs.NewLogger()
	DefaultLog.EnableFuncCallDepth(true)
	DefaultLog.SetLogFuncCallDepth(3)
	if runMode == "dev" {
		DefaultLog.SetLogger(logs.AdapterConsole)
	} else {
		workPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		logDir := filepath.Join(workPath)
		for _, dir := range dirs {
			logDir = filepath.Join(logDir, dir)
			if !isPathExist(logDir) {
				oldMask := syscall.Umask(0)
				err = os.Mkdir(logDir, 0755)
				if err != nil {
					println(err.Error(), logDir)
					panic(err)
				}
				syscall.Umask(oldMask)
			}
		}
		logFile := filepath.Join(logDir, "log.txt")
		if !isPathExist(logDir) {
			oldMask := syscall.Umask(0)
			err = os.Mkdir(logDir, 0755)
			if err != nil {
				println(err.Error(), logDir)
				panic(err)
			}
			f, err := os.Create(logFile)
			if err != nil {
				println(err.Error(), logFile)
				panic(err)
			}
			f.Close()
			syscall.Umask(oldMask)
		} else {
			if !isPathExist(logFile) {
				f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
				if err != nil {
					println(err.Error(), logFile)
					panic(err)
				}
				f.Close()
			}
		}
		configs := `{"filename":"` + logFile + `","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`
		err = DefaultLog.SetLogger(logs.AdapterMultiFile, configs)
		if err != nil {
			println(err.Error())
			panic(err)
		}
	}
	return DefaultLog
}
