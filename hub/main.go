package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/postfinance/single"
)

func main() {
	one, err := single.New("cvthub", single.WithLockPath("/tmp"))
	failOnErr("%v", err)
	failOnErr("%v", one.Lock())
	defer func() {
		closeSubServers()
		failOnErr("%v", one.Unlock())
		fPln("hub exit")
	}()

	startSubServers("./subsvr.md")
	time.Sleep(1 * time.Second)

	initRtModifier()

	// Start Service
	done := make(chan string)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go HostHTTPAsync(c, done)
	<-done
	// logGrp.Do(<-done)
}

func shutdownAsync(e *echo.Echo, sig <-chan os.Signal, done chan<- string) {
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	failOnErr("%v", e.Shutdown(ctx))
	time.Sleep(20 * time.Millisecond)
	done <- "Shutdown Successfully"
}

// HostHTTPAsync : Host a HTTP Server for XML to JSON
func HostHTTPAsync(sig <-chan os.Signal, done chan<- string) {
	// defer logGrp.Do("HostHTTPAsync Exit")

	e := echo.New()
	defer e.Close()

	// waiting for shutdown
	go shutdownAsync(e, sig, done)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2G"))
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST},
		AllowCredentials: true,
	}))

	e.Logger.SetOutput(os.Stdout)
	e.Logger.Infof(" ------------------------ e.Logger.Infof ------------------------ ")

	defer e.Start(fSf(":%d", PORT))
	// logGrp.Do("Echo Service is Starting ...")

	// ------------------------------------------------------------------------------------ //

	routeFun := func(method, svr string) func(c echo.Context) error {
		return func(c echo.Context) (err error) {
			var (
				status = http.StatusOK
				resp   *http.Response
				ret    []byte
				url    = mSvrRedirect[svr]
			)
			if ok, paramstr := urlParamStr(c.QueryParams()); ok {
				url += "?" + paramstr
			}

			switch method {
			case "GET":
				resp, err = http.Get(url)
			case "POST":
				resp, err = http.Post(url, "application/json", c.Request().Body)
			default:
				panic("Only Support [GET POST]")
			}

			if err != nil {
				ret = []byte(err.Error())
				status = http.StatusInternalServerError
				goto ERR_RET
			}
			if ret, err = io.ReadAll(resp.Body); err != nil {
				ret = []byte(err.Error())
				status = http.StatusInternalServerError
				goto ERR_RET
			}

		ERR_RET:
			retstr := ""
			for _, m := range modifiers {
				retstr = m.ModifyRet(svr, string(ret))
			}

			return c.String(status, retstr) // If already JSON String, so return String
		}
	}

	for svr, path := range mSvrGETPath {
		e.GET(path, routeFun("GET", svr))
	}

	for svr, path := range mSvrPOSTPath {
		e.POST(path, routeFun("POST", svr))
	}
}
