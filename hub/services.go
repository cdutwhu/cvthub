package main

import (
	"os/exec"
	"path/filepath"
	"sync"
	"time"

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
	qSvrExePath   = make([]string, 0)
	mSvrExeArgs   = make(map[string]string)
	mutex         = &sync.Mutex{}
	qSvrPid       = make([]string, 0)
	qSvrPidExist  []string
	mApiReDirGET  = make(map[string]string)
	mApiReDirPOST = make(map[string]string)
)

func at(items []string, i int) string {
	return sTrim(items[i], " \t")
}

func loadSvrTable(subSvrFile string) {

	_, err := scanLine(subSvrFile, func(ln string) (bool, string) {

		ln = sTrim(ln, " \t|") // also remove markdown table left & right '|'
		ss := sSplit(ln, "|")
		api, exe, args, reDir, method, enable := at(ss, iAPI), at(ss, iExePath), at(ss, iArgs), at(ss, iRedir), at(ss, iMethod), at(ss, iEnable)

		if enable != "true" {
			return false, ""
		}

		if exe != "" {
			exePath, err := AbsPath(exe, true)
			failOnErr("%v", err)
			qSvrExePath = ts.MkSet(append(qSvrExePath, exePath)...)
			mSvrExeArgs[exePath] = args
		}

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

		return true, ""

	}, "")

	failOnErr("%v", err)
}

// TODO: put it into github.com/digisan/gotk/
func GetRunningPID(pathOfExe string) (pidGrp []string) {
	abspath, err := filepath.Abs(pathOfExe)
	failOnErr("%v", err)
	dir, exe := filepath.Dir(abspath), filepath.Base(abspath)
	out, err := exec.Command("/bin/sh", "-c", "pgrep "+exe).CombinedOutput()
	if fSf("%v", err) == "exit status 1" {
		return
	}
	failOnErr("%v", err)

	outstr := sTrim(string(out), " \t\r\n")
	for _, pid := range sSplit(outstr, "\n") {
		out, err := exec.Command("/bin/sh", "-c", "pwdx "+pid).CombinedOutput()
		if fSf("%v", err) == "exit status 1" {
			return
		}
		failOnErr("%v", err)

		outstr := sTrim(string(out), " \t\r\n")
		procpath := sSplit(outstr, ": ")[1]
		if dir == procpath {
			pidGrp = append(pidGrp, pid)
		}
	}
	return
}

// TODO: put it into github.com/digisan/gotk/
func ExistRunningPS(pathOfExe string) bool {
	return len(GetRunningPID(pathOfExe)) > 0
}

func launchServers(subSvrFile string, launched chan<- struct{}) {

	loadSvrTable(subSvrFile)

	for _, exePath := range qSvrExePath {

		ok := make(chan struct{})

		// start executable
		go func(exePath string) {
			fPf("<%s> is starting...\n", exePath)

			// check existing PS
			if qSvrPidExist = GetRunningPID(exePath); len(qSvrPidExist) > 0 {
				closed := make(chan struct{})
				go closeServers(false, closed)
				<-closed
				failOnErr("%v", fEf("%v exists", exePath))
			}

			ok <- struct{}{}

			// start executable
			exeWithArgs := exePath + " " + mSvrExeArgs[exePath]
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

		}(exePath)

		// collect PID
		go func(exePath string) {
			<-ok
			I := 0
			for {
				time.Sleep(loopInterval * time.Millisecond)
				if pidGrp := GetRunningPID(exePath); pidGrp != nil {
					mutex.Lock()
					qSvrPid = append(qSvrPid, pidGrp...)
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
						if ExistRunningPS(exePath) {
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
