package main

import (
	"os/exec"
	"path/filepath"
	"time"

	"github.com/digisan/gotk/slice/ts"
)

// table header order
const (
	iSvr = iota
	iAPI
	iSubSvrDir
	iExe
	iRedir
	iMethod
	iEnable
)

var (
	mSvrGETPath  = make(map[string]string)
	mSvrPOSTPath = make(map[string]string)
	mSvrPkgPath  = make(map[string]string)
	mSvrExeName  = make(map[string]string)
	mSvrRedirect = make(map[string]string)
)

func at(items []string, i int) string {
	return sTrim(items[i], " \t")
}

func initSubSvr(subSvrFile string) {
	_, err := scanLine(subSvrFile, func(ln string) (bool, string) {
		ln = sTrim(ln, " \t|") // also remove markdown table left & right '|'
		ss := sSplit(ln, "|")
		svr, api, ssDir, exe, reDir, enable := "", "", "", "", "", ""
		if sContains(ln, "GET") || sContains(ln, "POST") {
			svr, api, ssDir, exe, reDir, enable = at(ss, iSvr), at(ss, iAPI), at(ss, iSubSvrDir), at(ss, iExe), at(ss, iRedir), at(ss, iEnable)
			if enable != "true" && enable != "TRUE" {
				return true, ""
			}
			abspath, err := filepath.Abs(ssDir)
			failOnErr("%v", err)
			mSvrPkgPath[svr] = "\"" + abspath + "\""
			mSvrExeName[svr] = exe
			mSvrRedirect[svr] = reDir
		}
		switch {
		case sContains(ln, "GET"):
			mSvrGETPath[svr] = api
		case sContains(ln, "POST"):
			mSvrPOSTPath[svr] = api
		}
		return true, ""
	}, "")
	failOnErr("%v", err)
}

func startSubServers(subSvrFile string) {
	initSubSvr(subSvrFile)
	for svr, exeDir := range mSvrPkgPath {
		go func(svr, wd, exe string) {
			fPln(svr, "is starting...")
			// failOnErr("%v @ %v", exec.Command("/bin/sh", "-c", "cd "+wd+" && ./"+exe).Run(), svr)
			_, err := exec.Command("/bin/sh", "-c", "cd "+wd+" && ./"+exe).CombinedOutput()
			switch {
			case fSf("%v", err) == "exit status 143":
				fPln(svr, "is shutting down...")
			case fSf("%v", err) == "signal: interrupt":
				fPln(svr, "is shutting down...")
			default:
				panic("NOT BE HERE! @ " + err.Error())
			}
		}(svr, exeDir, mSvrExeName[svr])
	}
}

func pidSubServers() (pidGrp []string) {
	for _, name := range mSvrExeName {
		cmd := exec.Command("/bin/sh", "-c", "pgrep "+name)
		out, err := cmd.CombinedOutput()
		failOnErr("%v", err)
		pidGrp = append(pidGrp, sSplit(sTrim(string(out), " \t\r\n"), "\n")...)
	}
	return ts.MkSet(pidGrp...)
}

func closeSubServers() {
	for _, pid := range pidSubServers() {
		go func(pid string) {
			failOnErr("%v @ %v", exec.Command("/bin/sh", "-c", "kill -15 "+pid).Run(), pid)
		}(pid)
	}
	time.Sleep(1 * time.Second)
}
