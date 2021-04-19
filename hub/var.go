package main

import (
	"fmt"
	"os/user"
	"path/filepath"
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
	sContains     = strings.Contains
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
	timeoutStart    = 1    // Second
	timeoutStartAll = 1    // Second
	timeoutCloseAll = 10   // Second
)

var (
	loopLmtStart    = timeoutStart * 1000 / loopInterval
	loopLmtStartAll = timeoutStartAll * 1000 / loopInterval
	loopLmtCloseAll = timeoutCloseAll * 1000 / loopInterval
)

// TODO: put it into github.com/digisan/gotk/io
var (
	AbsPath = func(path string, check bool) (string, error) {
		if sHasPrefix(path, "~/") {
			user, err := user.Current()
			failOnErr("%v", err)
			path = user.HomeDir + path[1:]
		}
		abspath, err := filepath.Abs(path)
		failOnErr("%v", err)

		if check && (!io.DirExists(abspath) && !io.FileExists(abspath)) {
			return abspath, fEf("%s doesn't exist", abspath)
		}
		return abspath, nil
	}
)
