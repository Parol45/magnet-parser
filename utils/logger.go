package utils

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)


func SetupLogger(logFileName string) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "./logs/" + logFileName + ".log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     31, //days
	}))
}
