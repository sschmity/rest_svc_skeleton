package logger

import (
	"sync"
	"log"
	"os"
	"io"
	"fmt"
	"runtime"
)

const (
	NO_DEBUG int = 0
	FATAL    int = 1
	ERROR    int = 2
	WARNING  int = 3
	INFO     int = 4
	DEBUG    int = 5

	START_METHOD = "[ START ]"
	END_METHOD   = "[ END ]"

)

type Logger struct {
	fileLoadedFlg bool
	fileName      string
	currLevel     int
	debug         *log.Logger
	info          *log.Logger
	warning       *log.Logger
	error         *log.Logger
	fatal         *log.Logger
}

type FuncIntString func(int, string)
type FuncIntStringMultiArgs func(int, string, ...interface{})

var (
	instanceLogger *Logger
	onceLogger     sync.Once
	stdout         io.Writer
	stderr         io.Writer
	Log            FuncIntString
	Logf           FuncIntStringMultiArgs
)


// init is used to initialized the package variables
// This function is called directly by golang.
// Please don't call this method in your code
func init() {
	Log = GetLoggerInstance().print
	Logf = GetLoggerInstance().printf
}



func GetLoggerInstance() *Logger {
	onceLogger.Do(func() {
		instanceLogger = &Logger{}
	})
	// set default leval
	instanceLogger.currLevel = DEBUG
	instanceLogger.fileLoadedFlg = false
	return instanceLogger
}

func (logPt *Logger) SetLogFile(aName string, aMultiWriterFlag ...bool) {

	f, err := os.OpenFile(aName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	//defer f.Close()

	instanceLogger.fileLoadedFlg = true

	if aMultiWriterFlag != nil && aMultiWriterFlag[0] == true {
		stdout, stderr = io.MultiWriter(f, os.Stdout), io.MultiWriter(f, os.Stderr)
	} else {
		stdout, stderr = io.MultiWriter(f), io.MultiWriter(f)
	}

	logPt.debug = log.New(stdout,
		"DEBUG:   ",
		//log.Ldate|log.Ltime|log.Lshortfile)
		log.Ldate|log.Ltime)

	logPt.info = log.New(stdout,
		"INFO:    ",
		//log.Ldate|log.Ltime|log.Lshortfile)
		log.Ldate|log.Ltime)

	logPt.warning = log.New(stdout,
		"WARNING: ",
		//log.Ldate|log.Ltime|log.Lshortfile)
		log.Ldate|log.Ltime)

	logPt.error = log.New(stderr,
		"ERROR:   ",
		//log.Ldate|log.Ltime|log.Lshortfile)
		log.Ldate|log.Ltime)

	logPt.fatal = log.New(stderr,
		"FATAL:   ",
		//log.Ldate|log.Ltime|log.Lshortfile)
		log.Ldate|log.Ltime)
}

func (logPt *Logger) SetLogLevel(aValue int) {
	logPt.currLevel = aValue
}

func (pt *Logger) IsActive() (bool) {
	return pt.currLevel > NO_DEBUG
}

func (pt *Logger) Log(aLevel int, aMsg string) {
	if aLevel <= pt.currLevel {
		pt.log(aLevel, fmt.Sprintf("%s :: %s", callFuncName(2), aMsg))
	}
}

func (pt *Logger) Logf(aLevel int, format string, a ...interface{}) {

	if aLevel <= pt.currLevel {
		pt.log(aLevel, fmt.Sprintf(fmt.Sprintf("%s,%s", callFuncName(2), format), a...))
	}
}


func (pt *Logger) print(aLevel int, aMsg string) {
	if aLevel <= pt.currLevel {
		pt.log(aLevel, fmt.Sprintf("%s :: %s", callFuncName(3), aMsg))
	}
}

func (pt *Logger) printf(aLevel int, format string, a ...interface{}) {

	if aLevel <= pt.currLevel {
		pt.log(aLevel, fmt.Sprintf(fmt.Sprintf("%s,%s", callFuncName(3), format), a...))
	}
}

func (pt *Logger) log(aLevel int, aMsg string) {

	if !pt.fileLoadedFlg {
		pt.SetLogFile("./Log.txt", false)
	}

	if aLevel <= pt.currLevel {

		switch aLevel {
		case INFO:
			pt.info.Print(aMsg)

		case DEBUG:
			pt.debug.Print(aMsg)

		case WARNING:
			pt.warning.Print(aMsg)

		case ERROR:
			pt.error.Print(aMsg)

		case FATAL:
			pt.fatal.Print(aMsg)
			//os.Exit(-1)

		default:
			pt.info.Print(aMsg)
		}
	}
}

func callFuncName(skip int) string {
	pc, _, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s(%d)", runtime.FuncForPC(pc).Name(), line)
}
