package utils

import (
	"github.com/laizy/log"
	"github.com/laizy/log/ext"
)

func InitLog(fileName string) {
	multi := log.MultiHandler(ext.RollingFileHandler(fileName), log.StdoutHandler)
	log.Root().SetHandler(multi)
}
