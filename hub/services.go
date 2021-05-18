package main

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/digisan/gotk/io"
	proc "github.com/digisan/gotk/process"
	"github.com/digisan/gotk/slice/ti"
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
	qDelay        = make([]struct{ s, e int }, 0)
	mutex         = &sync.Mutex{}
	qPid          = make([]int, 0)
	mApiReDirGET  = make(map[string]string)
	mApiReDirPOST = make(map[string]string)
)

func at(items []string, i int) string {
	return envValued(sTrim(items[i], " \t"), nil) // use environment variables
}

func str2delay(s string) struct{ s, e int } {
	se := sSplitN(s, ",", 2)
	switch len(se) {
	case 2:
		ns, err := strconv.Atoi(se[0])
		if err != nil {
			ns = 0
		}
		ne, err := strconv.Atoi(se[1])
		if err != nil {
			ne = 0
		}
		return struct{ s, e int }{ns, ne}
	case 1:
		ns, err := strconv.Atoi(se[0])
		if err != nil {
			ns = 0
		}
		return struct{ s, e int }{ns, 0}
	default:
		return struct{ s, e int }{0, 0}
	}
}

func loadSvrTable(subSvrFile string) {

	// get variables defined in services.md & set them to environment variables
	chunk2map(subSvrFile, "```export", "```", "=", true, true)

	_, err := scanLine(subSvrFile, func(ln string) (bool, string) {

		ln = sTrim(ln, " \t")

		// only deal with table rows
		if !(sHasPrefix(ln, "|") && sHasSuffix(ln, "|")) {
			return false, ""
		}

		ss := sSplit(sTrim(ln, "|"), "|") // remove markdown table left & right '|', then split by '|'
		failOnErrWhen(len(ss) != 7, "%v", "services.md table must have 7 columns, check it")

		var (
			exe    = at(ss, iExePath)
			args   = at(ss, iArgs)
			delay  = at(ss, iDelay)
			api    = at(ss, iAPI)
			reDir  = at(ss, iRedir)
			method = at(ss, iMethod)
			enable = at(ss, iEnable)
		)

		// only care about [ENABLE-true] rows
		if enable != "true" {
			return false, ""
		}

		if exe != "" {
			exePath, err := io.AbsPath(exe, true) // validate each executable
			failOnErr("%v", err)
			qExePath = append(qExePath, exePath) // same executable could be invoked multiple times // ts.MkSet(append(qSvrExePath, exePath)...)
			qExeArgs = append(qExeArgs, args)
			qDelay = append(qDelay, str2delay(delay))
		}

		// validate qExePath (already done)
		// failOnErrWhen(!io.FilesAllExist(qExePath), "%v", fEf("Not All Executables Are In Valid Path"))

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
				failOnErr("%v", fEf("At present, only [GET POST] are supported, check mark-down service table"))
			}
		}

		return true, ""

	}, "")

	failOnErr("%v", err)
}

func StartExe(path, arg string) (string, error) {
	cmdstr := fSf("cd %s && %s %s", filepath.Dir(path), path, arg)
	cmd := exec.Command("/bin/sh", "-c", cmdstr)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func StopExe(pid int) error {
	cmdstr := fSf("kill -15 %d", pid)
	return exec.Command("/bin/sh", "-c", cmdstr).Run()
}

func LaunchServices(subSvrFile string, chkRunning bool, launched chan<- struct{}) {

	loadSvrTable(subSvrFile)

	chStartErr := make(chan error, len(qExePath))

	for i, exePath := range qExePath {
		time.Sleep(80 * time.Millisecond) // if no sleep, simultaneously start same executable may fail.

		ok := make(chan struct{})

		// start executable
		go func(i int, exePath string) {
			time.Sleep(time.Duration(qDelay[i].s) * time.Second)
			info("<%s> is starting...", exePath)

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
			output, err := StartExe(exePath, qExeArgs[i])

			// log each service's output to file, begin with command line string
			mtx4log.Lock()
			l2c(false)
			l2f(true, fSf("%s%02d#%s.log", logpath, i, filepath.Base(exePath)))
			info("%s %s\n%s", exePath, qExeArgs[i], output)
			l2c(true)
			l2f(false, "")
			mtx4log.Unlock()

			// exitSHPid := fSf("%d", cmd.Process.Pid)

			// check exited status
			if err == nil {
				info("<%s> is shutting down...", exePath)
				return
			}
			msg := fSf("%v", err)
			switch msg {
			case "exit status 1", "exit status 143", "signal: interrupt":
				info("<%s> is shutting down...<%s>", exePath, msg)
			default:
				chStartErr <- fEf("<%s> cannot be started @error: %v", exePath, err)
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
					qPid = ti.MkSet(append(qPid, pidGrp...)...)
					info("<%s> is running...", exePath)
					mutex.Unlock()
					break
				}
				I++
				if I > loopLmtStartOne {
					chStartErr <- fEf("Cannot start <%s> as service in %d(s)", exePath, timeoutStartOne)
				}
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
			if I > loopLmtStartAll {
				chStartErr <- fEf("Cannot successfully start all services in %d(s)", timeoutStartAll)
			}
		}
	}()

	// check services starting status
	time.Sleep(1 * time.Second)
	select {
	case msg := <-chStartErr:
		warnOnErr("%v", msg)
		closed := make(chan struct{})
		go closeServers(false, closed)
		<-closed
		failOnErr("Hub Abort as: %v", msg)

	case <-time.After(timeoutCloseAll * time.Second):
		info("No Services Starting Errors Detected in %d(s)", timeoutCloseAll)
	}

	// monitor services status
	chMStop := make(chan bool)
	chMMsg := make(chan string)
	go monitorServices(chMMsg, chMStop)
	for msg := range chMMsg {
		info(msg)
	}
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

	for i, pid := range qPid {
		time.Sleep(20 * time.Millisecond)

		go func(i, pid int) {
			time.Sleep(time.Duration(qDelay[i].e) * time.Second)

			err := StopExe(pid)
			if err == nil {
				info("PID<%d> is shutting down...", pid)
				return
			}
			switch fSf("%v", err) {
			case "exit status 1":
				info("PID<%d> is shutting down...<%v>", pid, err)
			default:
				failOnErr("PID<%d> shutdown error @Error: %v", pid, err)
			}

		}(i, pid)
	}
}

func monitorServices(msg chan<- string, stop <-chan bool) {
	ticker := time.NewTicker(monitorInterval * time.Second)
	for {
		select {
		case <-stop:
			ticker.Stop()
			return
		case <-ticker.C:
			for i, path := range qExePath {
				if !proc.ExistRunningPS(path) {
					msg <- fSf("<%s> process @ <%s> exited", path, qPid[i])
				}
			}
		}
	}
}
