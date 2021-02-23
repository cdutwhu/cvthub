package main

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/rflx"
)

var (
	fEf        = fmt.Errorf
	fSf        = fmt.Sprintf
	fPln       = fmt.Println
	sTrim      = strings.Trim
	sSplit     = strings.Split
	sHasPrefix = strings.HasPrefix
	failOnErr  = fn.FailOnErr
	readLine   = io.EditFileByLine
	mapKeys    = rflx.MapKeys
)
