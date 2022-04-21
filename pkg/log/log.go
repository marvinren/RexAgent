package log

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func InitLogger(logpath string){
	log_file_name := logpath + "/rexagent.log"
	if err2 := os.MkdirAll(logpath, 766);err2 != nil{
		panic(errors.New("log path can't create"))
	}
	file, err := os.OpenFile(log_file_name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	if err == nil {
		logger.SetOutput(multiWriter)
	} else {
		logger.Printf("can not open log file %s", log_file_name)
	}
}



func Log(topic string, level string, format string, v ...interface{}) {
	prefix := fmt.Sprintf("%s [%s] ", level, topic)
	if len(v) == 0 {
		logger.Println(prefix + format)
	} else {
		logger.Printf(prefix+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	Log("rexagent", "INFO", format, v...)
}

func Warn(format string, v ...interface{}) {
	Log("rexagent", "WARN", format, v...)
}

func Error(format string, v ...interface{}) {
	Log("rexagent", "ERROR", format, v...)
}


