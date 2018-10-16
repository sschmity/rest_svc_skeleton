package main

import (
	"flag"
	"strconv"
	"strings"
	"os"
	"net/http"
)

var (
	argConfigFileName string
	argLogLevel       int
	argLogFile        string
	argLogMultiWriter bool
	log               = logger.GetLoggerInstance()
)

func main() {

	parseCmdLine()

	log.SetLogFile(argLogFile, argLogMultiWriter)
	log.SetLogLevel(argLogLevel)

	log.Log(logger.INFO, "Service started using logfile: " + argLogFile)
	log.Log(logger.INFO, "Log Level: "+strconv.Itoa(argLogLevel)+"(NO_DEBUG=0, FATAL=1, ERROR= 2, WARNING=3, INFO=4, DEBUG=5)")
	if len(flag.Args()) != 0 {
		log.Log(logger.ERROR, "Extra Command line arguments not taken into account: "+strings.Join(flag.Args(), " + "))
	}

}

func parseCmdLine() {
	flag.StringVar(&argConfigFileName, "config", "", "Name of config file")
	flag.StringVar(&argLogFile, "log", "./AsBddRestSvc.log", "Name of log file")
	flag.IntVar(&argLogLevel, "loglevel", 4, "log level is 4 by default (NO_DEBUG=0, FATAL=1, ERROR= 2, WARNING=3, INFO=4, DEBUG=5)")
	flag.BoolVar(&argLogMultiWriter, "multiWriter", false, "multi writer is false by default (true allows file and Stdout logging)")

	flag.Parse()
}
