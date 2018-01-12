package main

import (
	"flag"
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	log "github.com/kdar/factorlog"
	"math/rand"
	"os"
	"runtime"
	"time"
)

var (
	logFlag  = flag.String("log", "", "set log path")
	portFlag = flag.Int("port", 9891, "set port")
	logger   *log.FactorLog
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func SetGlobalLogger(logPath string) *log.FactorLog {
	sfmt := `%{Color "red:white" "CRITICAL"}%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "blue" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}:%{ShortFile}:%{Line}] %{Message}%{Color "reset"}`
	logger := log.New(os.Stdout, log.NewStdFormatter(sfmt))
	if len(logPath) > 0 {
		logf, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			return logger
		}
		logger = log.New(logf, log.NewStdFormatter(sfmt))
	}
	logger.SetSeverities(log.INFO | log.WARN | log.ERROR | log.FATAL | log.CRITICAL)
	return logger
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	logger = SetGlobalLogger(*logFlag)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	requestGroup := router.Group("/urlproxy")
	{
		requestGroup.GET("/", UrlProxyHandler)
	}

	logger.Infof("URL Proxy Server started at:0.0.0.0:%d", *portFlag)
	defer func() {
		logger.Infof("URL Proxy Server exit from:0.0.0.0:%d", *portFlag)
	}()
	endless.ListenAndServe(fmt.Sprintf(":%d", *portFlag), router)
}
