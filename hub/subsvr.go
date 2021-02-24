package main

import (
	"os/exec"
	"path/filepath"
	"time"
)

var (
	mSvrGETPath  = make(map[string]string)
	mSvrPOSTPath = make(map[string]string)
	mSvrPkgPath  = make(map[string]string)
	mSvrExeName  = make(map[string]string)
	mSvrRedirect = make(map[string]string)
)

func initSubSvr(subSvrFile string) {
	readLine(subSvrFile, func(ln string) (bool, string) {
		ln = sTrim(ln, " \t")
		ss := sSplit(ln, "|")
		svr, api, exeDir, exeName, reDir := "", "", "", "", ""
		if sContains(ln, "GET") || sContains(ln, "POST") {
			svr, api, exeDir, exeName, reDir = sTrim(ss[0], " \t"), sTrim(ss[1], " \t"), sTrim(ss[2], " \t"), sTrim(ss[3], " \t"), sTrim(ss[4], " \t")
			abspath, err := filepath.Abs(exeDir)
			failOnErr("%v", err)
			mSvrPkgPath[svr] = "\"" + abspath + "\""
			mSvrExeName[svr] = exeName
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
	// remove duplicated pid
	m := make(map[string]struct{})
	for _, pid := range pidGrp {
		m[pid] = struct{}{}
	}
	pidGrp = mapKeys(m).([]string)
	return
}

func closeSubServers() {
	for _, pid := range pidSubServers() {
		go func(pid string) {
			failOnErr("%v @ %v", exec.Command("/bin/sh", "-c", "kill -15 "+pid).Run(), pid)
		}(pid)
	}
	time.Sleep(1 * time.Second)
}
