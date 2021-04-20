package main

import (
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/digisan/gotk/io"
	proc "github.com/digisan/gotk/process"
	"github.com/digisan/gotk/slice/ts"
)

// table header order
const (
	iAPI = iota
	iExePath
	iArgs
	iRedir
	iMethod
	iEnable
)

var (
	qSvrExePath   = make([]string, 0) // server may be repeated
	qSvrExeArgs   = make([]string, 0)
	mutex         = &sync.Mutex{}
	qSvrPid       = make([]string, 0)
	mApiReDirGET  = make(map[string]string)
	mApiReDirPOST = make(map[string]string)
)

func at(items []string, i int) string {
	return sTrim(items[i], " \t")
}

func loadSvrTable(subSvrFile string) {

	_, err := scanLine(subSvrFile, func(ln string) (bool, string) {

		ss := sSplit(sTrim(ln, "|"), "|") // remove markdown table left & right '|', then split by '|'
		failOnErrWhen(len(ss) != 6, "%v", "services.md must be 6 columns, check it")
		api, exe, args, reDir, method, enable := at(ss, iAPI), at(ss, iExePath), at(ss, iArgs), at(ss, iRedir), at(ss, iMethod), at(ss, iEnable)

		if enable != "true" {
			return false, ""
		}

		if exe != "" {
			exePath, err := io.AbsPath(exe, true)
			failOnErr("%v", err)
			qSvrExePath = append(qSvrExePath, exePath) // same executable could be started multiple times // ts.MkSet(append(qSvrExePath, exePath)...)
			qSvrExeArgs = append(qSvrExeArgs, args)
		}

		if api != "" {
			if sHasPrefix(reDir, ":") {
				reDir = "http://localhost" + reDir
			}
			switch method {
			case "GET":
				mApiReDirGET[api] = reDir
			case "POST":
				mApiReDirPOST[api] = reDir
			default:
				panic("Only [GET POST] are Supported")
			}
		}

		return true, ""

	}, "")

	failOnErr("%v", err)
}

func launchServers(subSvrFile string, chkRunning bool, launched chan<- struct{}) {

	loadSvrTable(subSvrFile)

	for i, exePath := range qSvrExePath {
		time.Sleep(80 * time.Millisecond) // if no sleep, simultaneously start same executable may fail.

		ok := make(chan struct{})

		// start executable
		go func(i int, exePath string) {
			fPf("<%s> is starting...\n", exePath)

			// check existing running PS
			if chkRunning {
				if qPidRunning := proc.GetRunningPID(exePath); len(qPidRunning) > 0 {
					closed := make(chan struct{})
					go closeServers(false, closed)
					<-closed
					failOnErr("%v", fEf("%v exists", exePath))
				}
			}

			ok <- struct{}{}

			// start executable
			exeWithArgs := exePath + " " + qSvrExeArgs[i]
			cmd := fSf("cd %s && %s", filepath.Dir(exePath), exeWithArgs)
			_, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
			if err == nil {
				fPf("<%s> is shutting down...\n", exePath)
				return
			}
			msg := fSf("%v", err)
			switch msg {
			case "exit status 1", "exit status 143", "signal: interrupt":
				fPf("<%s> is shutting down...<%s>\n", exePath, msg)
			default:
				panic(fSf("<%s> cannot be started @Error: %v", exePath, err.Error()))
			}

		}(i, exePath)

		// collect PID
		go func(exePath string) {
			<-ok
			I := 0
			for {
				time.Sleep(loopInterval * time.Millisecond)
				if pidGrp := proc.GetRunningPID(exePath); pidGrp != nil {
					mutex.Lock()
					qSvrPid = ts.MkSet(append(qSvrPid, pidGrp...)...)
					mutex.Unlock()
					break
				}
				I++
				failOnErrWhen(I > loopLmtStart, "%v", fEf("Cannot start server @ <%s> in %d(s)", exePath, timeoutStart))
			}
		}(exePath)
	}

	go func() {
		I := 0
		for {
			time.Sleep(loopInterval * time.Millisecond)
			if len(qSvrExePath) == len(qSvrPid) {
				launched <- struct{}{}
				break
			}
			I++
			failOnErrWhen(I > loopLmtStartAll, "%v", fEf("Cannot start all servers in %d(s)", timeoutStartAll))
		}
	}()
}

func closeServers(check bool, closed chan<- struct{}) {
	defer func() {
		if check {

			go func() {
				I := 0
			LOOP:
				for {
					for _, exePath := range qSvrExePath {
						if proc.ExistRunningPS(exePath) {
							time.Sleep(loopInterval * time.Millisecond)
							I++
							failOnErrWhen(I > loopLmtCloseAll, "%v", fEf("Cannot close all servers in %d(s)", timeoutCloseAll))
							continue LOOP
						}
					}
					closed <- struct{}{}
					break
				}
			}()

		} else {
			closed <- struct{}{}
		}
	}()

	for _, pid := range qSvrPid {
		time.Sleep(20 * time.Millisecond)

		go func(pid string) {
			cmd := fSf("kill -15 %s", pid)
			err := exec.Command("/bin/sh", "-c", cmd).Run()
			if err == nil {
				fPf("<%s> is shutting down...\n", pid)
				return
			}
			msg := fSf("%v", err)
			switch msg {
			case "exit status 1":
				fPf("<%s> is shutting down...<%s>\n", pid, msg)
			default:
				panic(fSf("<%s> shutdown error @Error: %v", pid, err))
			}

		}(pid)
	}
}
