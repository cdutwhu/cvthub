package main

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/net"
	"github.com/digisan/gotk/io"
)

var (
	fSf           = fmt.Sprintf
	fPln          = fmt.Println
	fPf           = fmt.Printf
	fEf           = fmt.Errorf
	sLastIndex    = strings.LastIndex
	sTrim         = strings.Trim
	sTrimLeft     = strings.TrimLeft
	sSplit        = strings.Split
	sHasPrefix    = strings.HasPrefix
	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	scanLine      = io.FileLineScan
	scanStrLine   = io.StrLineScan
	urlParamStr   = net.URLParamStr
	localIP       = net.LocalIP
)

const (
	PORT            = 1423 // PORT : this server port
	loopInterval    = 200  // Millisecond
	timeoutStart    = 6    // Second
	timeoutStartAll = 10   // Second
	timeoutCloseAll = 10   // Second
)

var (
	loopLmtStart    = timeoutStart * 1000 / loopInterval
	loopLmtStartAll = timeoutStartAll * 1000 / loopInterval
	loopLmtCloseAll = timeoutCloseAll * 1000 / loopInterval
)
