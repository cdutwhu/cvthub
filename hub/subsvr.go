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
	iService = iota
	iAPI
	iSvrPath
	iArgs
	iRedir
	iMethod
	iEnable
)

var (
	qSvrExePath  = make([]string, 0)
	mutex        = &sync.Mutex{}
	qSvrPid      = make([]string, 0)
	qSvrPidExist []string
	mSvrRedirect = make(map[string]string)
	mSvrGETPath  = make(map[string]string)
	mSvrPOSTPath = make(map[string]string)
)

func at(items []string, i int) string {
	return sTrim(items[i], " \t")
}

func loadSvrTable(subSvrFile string) {

	_, err := scanLine(subSvrFile, func(ln string) (bool, string) {

		ln = sTrim(ln, " \t|") // also remove markdown table left & right '|'
		ss := sSplit(ln, "|")
		if sContains(ln, "GET") || sContains(ln, "POST") {

			service, api, exe, reDir, enable := at(ss, iService), at(ss, iAPI), at(ss, iSvrPath), at(ss, iRedir), at(ss, iEnable)
			if enable != "true" && enable != "TRUE" {
				return true, ""
			}

			abspath, err := filepath.Abs(exe)
			failOnErr("%v", err)
			mSvrRedirect[service] = reDir
			qSvrExePath = ts.MkSet(append(qSvrExePath, abspath)...)

			switch {
			case sContains(ln, "GET"):
				mSvrGETPath[service] = api
			case sContains(ln, "POST"):
				mSvrPOSTPath[service] = api
			}
		}
		return true, ""

	}, "")
	failOnErr("%v", err)
}

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
			cmd := fSf("cd %s && %s", filepath.Dir(exePath), exePath)
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
