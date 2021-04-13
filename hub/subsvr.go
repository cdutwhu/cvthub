package main

import (
	"os/exec"
	"path/filepath"
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

func launchServers(subSvrFile string) {

	loadSvrTable(subSvrFile)

	for _, exe := range qSvrExePath {
		go func(exe string) {
			fPf("<%s> is starting...\n", exe)
			_, err := exec.Command("/bin/sh", "-c", exe).CombinedOutput()
			switch {
			case fSf("%v", err) == "exit status 143":
				fPf("<%s> is shutting down...(143)\n", exe)
			case fSf("%v", err) == "signal: interrupt":
				fPf("<%s> is shutting down...(int)\n", exe)
			default:
				panic("NOT BE HERE! @ " + err.Error())
			}
		}(exe)
	}
}

func pidServers() (pidGrp []string) {
	for _, path := range qSvrExePath {
		name := filepath.Base(path)
		cmd := exec.Command("/bin/sh", "-c", "pgrep "+name)
		out, err := cmd.CombinedOutput()
		failOnErr("%v", err)
		fPln(string(out))
		pidGrp = append(pidGrp, sSplit(sTrim(string(out), " \t\r\n"), "\n")...)
	}
	return ts.MkSet(pidGrp...)
}

func closeServers() {
	for _, pid := range pidServers() {
		go func(pid string) {
			failOnErr("%v @ %v", exec.Command("/bin/sh", "-c", "kill -15 "+pid).Run(), pid)
		}(pid)
	}
	time.Sleep(500 * time.Millisecond)
}
