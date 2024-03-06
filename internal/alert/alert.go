package alert

import (
	lg "log"
	"os"
)

var (
	infoLog    = lg.New(os.Stdout, "[INFO] ", lg.Lshortfile)
	errorLog   = lg.New(os.Stderr, "[ERROR] ", lg.Lshortfile)
	warningLog = lg.New(os.Stderr, "[WARNING] ", lg.Lshortfile)
)

func Info(v ...any) {
	infoLog.Println(v...)
}

func Error(v ...any) {
	errorLog.Println(v...)
}

func Warning(v ...any) {
	warningLog.Println(v...)
}
