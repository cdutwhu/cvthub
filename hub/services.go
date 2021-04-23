package main

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/digisan/gotk/io"
	proc "github.com/digisan/gotk/process"
	"github.com/digisan/gotk/slice/ts"
)

// table header order
const (
	iExePath = iota
	iArgs
	iDelay
	iAPI
	iRedir
	iMethod
	iEnable
)

var (
	qExePath      = make([]string, 0) // server may be repeated
	qExeArgs      = make([]string, 0)
	qStartDelay   = make([]int, 0)
	mutex         = &sync.Mutex{}
	qPid          = make([]string, 0)
	mApiReDirGET  = make(map[string]string)
	mApiReDirPOST = make(map[string]string)
)

func at(items []string, i int) string {
	return sTrim(items[i], " \t")
}

func loadSvrTable(subSvrFile string) {

	_, err := scanLine(subSvrFile, func(ln string) (bool, string) {

		ss := sSplit(sTrim(ln, "|"), "|") // remove markdown table left & right '|', then split by '|'
		failOnErrWhen(len(ss) != 7, "%v", "services.md must be 7 columns, check it")

		var (
			exe    = at(ss, iExePath)
			args   = at(ss, iArgs)
			delay  = at(ss, iDelay)
			api    = at(ss, iAPI)
			reDir  = at(ss, iRedir)
			method = at(ss, iMethod)
			enable = at(ss, iEnable)
		)

		if enable != "true" {
			return false, ""
		}

		if exe != "" {
			exePath, err := io.AbsPath(exe, true)
			failOnErr("%v", err)
			qExePath = append(qExePath, exePath) // same executable could be started multiple times // ts.MkSet(append(qSvrExePath, exePath)...)
			qExeArgs = append(qExeArgs, args)
			nDelay, err := strconv.Atoi(delay)
			if err != nil {
				nDelay = 0
			}
			qStartDelay = append(qStartDelay, nDelay)
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

	for i, exePath := range qExePath {
		time.Sleep(80 * time.Millisecond) // if no sleep, simultaneously start same executable may fail.

		ok := make(chan struct{})

		// start executable
		go func(i int, exePath string) {
			time.Sleep(time.Duration(qStartDelay[i]) * time.Second)
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
			cmdstr := fSf("cd %s && %s %s", filepath.Dir(exePath), exePath, qExeArgs[i])
			cmd := exec.Command("/bin/sh", "-c", cmdstr)
			_, err := cmd.CombinedOutput()

			// fPln(cmd.Process.Pid)

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
				if pidGrp := proc.GetRunningPID(exePath); len(pidGrp) > 0 {
					mutex.Lock()
					qPid = ts.MkSet(append(qPid, pidGrp...)...)
					fPf("<%s> is running...\n", exePath)
					mutex.Unlock()
					break
				}
				I++
				failOnErrWhen(I > loopLmtStart, "%v", fEf("Cannot start <%s> as service in %d(s)", exePath, timeoutStart))
			}
		}(exePath)
	}

	go func() {
		I := 0
		for {
			time.Sleep(loopInterval * time.Millisecond)
			if len(qExePath) == len(qPid) {
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
					for _, exePath := range qExePath {
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

	for _, pid := range qPid {
		time.Sleep(20 * time.Millisecond)

		go func(pid string) {
			cmdstr := fSf("kill -15 %s", pid)
			err := exec.Command("/bin/sh", "-c", cmdstr).Run()
			if err == nil {
				fPf("PID<%s> is shutting down...\n", pid)
				return
			}
			msg := fSf("%v", err)
			switch msg {
			case "exit status 1":
				fPf("PID<%s> is shutting down...<%s>\n", pid, msg)
			default:
				panic(fSf("PID<%s> shutdown error @Error: %v", pid, err))
			}

		}(pid)
	}
}
