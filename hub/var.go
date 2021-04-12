package main

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/net"
	"github.com/digisan/gotk/io"
)

var (
	fSf         = fmt.Sprintf
	fPln        = fmt.Println
	sLastIndex  = strings.LastIndex
	sTrim       = strings.Trim
	sTrimLeft   = strings.TrimLeft
	sSplit      = strings.Split
	sHasPrefix  = strings.HasPrefix
	sContains   = strings.Contains
	failOnErr   = fn.FailOnErr
	scanLine    = io.FileLineScan
	scanStrLine = io.StrLineScan
	urlParamStr = net.URLParamStr
	localIP     = net.LocalIP
)

const (
	// PORT : this server port
	PORT = 1323
)
