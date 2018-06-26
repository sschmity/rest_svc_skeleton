package main

import (
	"cs.as_rest_svc/event"
	"cs.as_rest_svc/sqlSvc"
	"cs.as_rest_svc/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"cs.as_rest_svc/cucumberSvc"
	"io"
	"flag"
	"io/ioutil"
)

var wg sync.WaitGroup
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var (
	argConfigFileNamePtr  *string
	argConfigTraceModePtr *bool
)

func main() {
	//read command line arguments
	readCommandLineArgs()
	if *argConfigTraceModePtr {
		Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatal(err)
	}
	Trace.Println(dir)

	//load json configuration file
	config.LoadConfigFile(*argConfigFileNamePtr)

	// initialize Log
	//	initializeLogConfig()
	// initialize server port
	initializeServerConfig()

	//register url services
	wg.Add(3)
	go func() {

		defer wg.Done()
		Info.Println("->event.InitRestAPI function started")
		event.InitRestAPI()
	}()
	go func() {
		defer wg.Done()
		Info.Println("->cucumber.InitRestAPI function started")
		cucumberSvc.InitRestAPI()
	}()
	go func() {
		defer wg.Done()
		Info.Println("->sqlSvc.InitRestAPI function started")
		sqlSvc.InitRestAPI()
	}()

	wg.Wait()

	//start server
	Error.Fatal(http.ListenAndServe(config.GetServerConfigInstance().GetServerPort(), nil))
	Info.Println("Http as_rest_svc service started at : " + config.GetServerConfigInstance().GetServerPort())
}

func initializeServerConfig() {
	a := config.GetServerConfigInstance()
	a.SetServerPort(config.ConfigFileInstance.HttpServerConfig.Port)
	Info.Println("Registered server port at : " + a.GetServerPort())
}

func initializeLogConfig() {
	a := config.GetlogConfigInstance()
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	a.SetAbsLogPath(dir+config.ConfigFileInstance.LogConfig.FilePath, config.ConfigFileInstance.LogConfig.FileName)
	config.SetLogger()
	Info.Println("Registered LogFileName : " + a.GetAbsLogPath())
}

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func readCommandLineArgs() {
	argConfigFileNamePtr = flag.String("configFile", "config.json", "json configuration file name")
	argConfigTraceModePtr = flag.Bool("traceModeOn", false, "trace mode is desactived by default")
	flag.Parse()
}
