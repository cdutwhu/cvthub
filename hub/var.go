package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/net"
	"github.com/digisan/gotk/io"
)

var (
	fSf           = fmt.Sprintf
	fPln          = fmt.Println
	fEf           = fmt.Errorf
	sLastIndex    = strings.LastIndex
	sTrim         = strings.Trim
	sTrimLeft     = strings.TrimLeft
	sSplit        = strings.Split
	sHasPrefix    = strings.HasPrefix
	sHasSuffix    = strings.HasSuffix
	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	warnOnErr     = fn.WarnOnErr
	info          = fn.Logger
	l2c           = fn.EnableLog2C
	l2f           = fn.EnableLog2F
	scanLine      = io.FileLineScan
	scanStrLine   = io.StrLineScan
	urlParamStr   = net.URLParamStr
	localIP       = net.LocalIP
)

const (
	PORT            = 1423 // PORT : this server port
	loopInterval    = 200  // Millisecond
	timeoutStartOne = 6    // Second
	timeoutStartAll = 10   // Second
	timeoutCloseAll = 10   // Second
	monitorInterval = 30   // Second
)

var (
	loopLmtStartOne = timeoutStartOne * 1000 / loopInterval
	loopLmtStartAll = timeoutStartAll * 1000 / loopInterval
	loopLmtCloseAll = timeoutCloseAll * 1000 / loopInterval
	logpath         = "./services_log/"
	mtx4log         = &sync.Mutex{}
)

func init() {
	log.SetFlags(log.LstdFlags) // overwrite "info/warn/fail" print style
}
