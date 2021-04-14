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

var (
	errDelay   time.Duration = 600 // Millisecond
	startDelay time.Duration = 400
	exitDelay  time.Duration = 400
	maxDelay   time.Duration = 1200
)

func launchServers(subSvrFile string) time.Duration {

	loadSvrTable(subSvrFile)

	for _, exePath := range qSvrExePath {

		ok := make(chan struct{})

		// start executable
		go func(exePath string) {
			fPf("<%s> is starting...\n", exePath)

			// check existing PS
			if qSvrPidExist = GetRunningPID(exePath); len(qSvrPidExist) > 0 {
				time.Sleep(errDelay * time.Millisecond)
				closeServers(false)
				failOnErr("%v", fEf("%v exists", exePath))
			}

			ok <- struct{}{}

			// start executable
			cmd := fSf("cd %s && %s", filepath.Dir(exePath), exePath)
			_, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
			switch {
			case fSf("%v", err) == "exit status 143":
				fPf("<%s> is shutting down...(143)\n", exePath)
			case fSf("%v", err) == "signal: interrupt":
				fPf("<%s> is shutting down...(int)\n", exePath)
			default:
				panic("NOT BE HERE! @ " + err.Error())
			}

		}(exePath)

		// collect PID
		go func(exePath string) {

			<-ok
			time.Sleep(startDelay * time.Millisecond)
			pidGrp := GetRunningPID(exePath)
			mutex.Lock()
			qSvrPid = append(qSvrPid, pidGrp...)
			mutex.Unlock()

		}(exePath)
	}

	return maxDelay
}

func closeServers(check bool) {
	defer func() {
		if check {
			for _, exePath := range qSvrExePath {
				failOnErrWhen(ExistRunningPS(exePath), "%v", fEf("%v is still running", exePath))
			}
		}
	}()

	for _, pid := range qSvrPid {
		go func(pid string) {
			failOnErr("%v @ %v", exec.Command("/bin/sh", "-c", "kill -15 "+pid).Run(), pid)
		}(pid)
	}
	time.Sleep(exitDelay * time.Millisecond)
}
