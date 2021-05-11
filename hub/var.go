package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/net"
	"github.com/digisan/gotk/io"
	"github.com/digisan/gotk/mapslice"
)

var (
	fSf           = fmt.Sprintf
	fPln          = fmt.Println
	fEf           = fmt.Errorf
	sContains     = strings.Contains
	sLastIndex    = strings.LastIndex
	sTrim         = strings.Trim
	sTrimLeft     = strings.TrimLeft
	sSplit        = strings.Split
	sJoin         = strings.Join
	sHasPrefix    = strings.HasPrefix
	sHasSuffix    = strings.HasSuffix
	sReplaceAll   = strings.ReplaceAll
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
	ksvs2slc      = mapslice.KsVs2Slc
)

const (
	PORT            = 1423 // PORT : this server port
	loopInterval    = 200  // Millisecond
	timeoutStartOne = 6    // Second
	timeoutStartAll = 10   // Second
	timeoutCloseAll = 10   // Second
	monitorInterval = 300  // Second
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

func chunk2map(filepath, markstart, markend, sep, keyprefix string) map[string]string {
	m := make(map[string]string)
	chunkproc := false
	_, err := scanLine(filepath, func(ln string) (bool, string) {
		if sHasPrefix(ln, markstart) && !chunkproc {
			chunkproc = true
			return false, ""
		}
		if sHasPrefix(ln, markend) && chunkproc {
			chunkproc = false
			return false, ""
		}
		if chunkproc && sContains(ln, sep) {
			ss := sSplit(ln, sep)
			m[keyprefix+ss[0]] = sJoin(ss[1:], sep)
		}
		return false, ""
	}, "")
	failOnErr("%v", err)
	return m
}
