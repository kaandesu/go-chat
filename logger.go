package main

import (
	"log"
	"os"
)

func NewLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	return log.New(logFile, "[kaandesu/go-chat]", log.Ltime)
}
