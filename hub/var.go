package main

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
)

var (
	fEf         = fmt.Errorf
	fSf         = fmt.Sprintf
	fPln        = fmt.Println
	sIndex      = strings.Index
	sLastIndex  = strings.LastIndex
	sTrim       = strings.Trim
	sTrimLeft   = strings.TrimLeft
	sTrimRight  = strings.TrimRight
	sSplit      = strings.Split
	sHasPrefix  = strings.HasPrefix
	sHasSuffix  = strings.HasSuffix
	sContains   = strings.Contains
	sReplace    = strings.Replace
	failOnErr   = fn.FailOnErr
	readLine    = io.EditFileByLine
	readStrLine = io.EditStrByLine
	mapKeys     = rflx.MapKeys
	urlParamStr = net.URLParamStr
	localIP     = net.LocalIP
)

const (
	// PORT : this server port
	PORT = 1323
)
