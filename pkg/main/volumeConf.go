package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var sourceVolPath string
var logPath string
var logFilename string

func init() {
	flag.StringVar(&logPath, "log-path", "/tmp/uruk/log", "`location` where uruk logs are stored")
	flag.StringVar(&logFilename, "log-filename", "uruk.log", "`filename` for logs")
}

func getLogfileName() string {
	return filepath.Join(logPath, logFilename)
}

func getLogfile() *os.File {
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(getLogfileName(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file
}
