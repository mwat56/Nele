/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package main

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	nele "github.com/mwat56/Nele"
	"github.com/mwat56/apachelogger"
	"github.com/mwat56/errorhandler"
)

// `doConsole()` checks for the `add` commandline argument, adds the
// text from StdIn as a new post, and terminates the program.
func doConsole(aMe string) {
	var (
		err error
		i64 int64
		s   string
	)
	if s, err = nele.AppArguments.Get("pa"); nil != err {
		// we assume, an error means: no cmd line action
		return
	}
	if "true" == s {
		if i64, err = nele.AddConsolePost(); nil != err {
			log.Fatalf("%s: %v", aMe, err)
		}
		log.Printf("\n\t%s wrote %d bytes in a new posting", aMe, i64)
		os.Exit(0)
	}
} // doConsole()

// `doFile()` checks for the `filename` commandline argument, adds the
// text from the file as a new post, and terminates the program.
func doFile(aMe string) {
	var (
		err error
		i64 int64
		s   string
	)
	if s, err = nele.AppArguments.Get("pf"); nil != err {
		// we assume, an error means: no cmd line action
		return
	}
	if 0 < len(s) {
		if i64, err = nele.AddFilePost(s); nil != err {
			log.Fatalf("%s: %v", aMe, err)
		}
		log.Printf("\n\t%s stored %d bytes in a new posting", aMe, i64)
		os.Exit(0)
	}
} // doFile()

// `fatal()` logs `aMessage` and terminates the program.
func fatal(aMessage string) {
	apachelogger.Err("Nele/main", aMessage)
	runtime.Gosched() // let the logger write
	apachelogger.Close()
	log.Fatalln(aMessage)
} // fatal()

// `setupSignals()` configures the capture of the interrupts `SIGINT`
// `and `SIGTERM` to terminate the program gracefully.
func setupSignals(aServer *http.Server) {
	// handle `CTRL-C` and `kill(15)`.
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for signal := range c {
			msg := fmt.Sprintf("%s captured '%v', stopping program and exiting ...", os.Args[0], signal)
			apachelogger.Err(`Nele/catchSignals`, msg)
			log.Println(msg)
			runtime.Gosched() // let the logger write
			if err := aServer.Shutdown(context.Background()); nil != err {
				fatal(fmt.Sprintf("%s: %v", os.Args[0], err))
			}
		}
	}()
} // setupSignals()

// `userCmdline()` checks for and executes password file commandline actions.
func userCmdline() {
	var (
		err   error
		fn, s string
	)
	if fn, err = nele.AppArguments.Get("uf"); (nil != err) || (0 == len(fn)) {
		return // without file no user handling
	}
	// All the following `nele.*` function calls will
	// terminate the program.
	if s, err = nele.AppArguments.Get("ua"); (nil == err) && (0 < len(s)) {
		nele.AddUser(s, fn)
	}
	if s, err = nele.AppArguments.Get("uc"); (nil == err) && (0 < len(s)) {
		nele.CheckUser(s, fn)
	}
	if s, err = nele.AppArguments.Get("ud"); (nil == err) && (0 < len(s)) {
		nele.DeleteUser(s, fn)
	}
	if s, err = nele.AppArguments.Get("ul"); (nil == err) && (0 < len(s)) {
		nele.ListUsers(fn)
	}
	if s, err = nele.AppArguments.Get("uu"); (nil == err) && (0 < len(s)) {
		nele.UpdateUser(s, fn)
	}
} // userCmdline()

// Actually run the program …
func main() {
	var (
		err       error
		handler   http.Handler
		ph        *nele.TPageHandler
		ck, cp, s string
	)
	Me, _ := filepath.Abs(os.Args[0])

	// Read INI files and commandline options
	nele.InitConfig()

	// Add a new posting via command line:
	doConsole(Me)

	// Read in a file as a new posting:
	doFile(Me)

	// Handle password file maintenance:
	userCmdline()

	if ph, err = nele.NewPageHandler(); nil != err {
		nele.ShowHelp()
		fatal(fmt.Sprintf("%s: %v", Me, err))
	}
	// Setup the errorpage handler:
	handler = errorhandler.Wrap(ph, ph)

	// Inspect `gzip` commandline argument and setup the Gzip handler:
	if s, err = nele.AppArguments.Get("gzip"); (nil == err) && ("true" == s) {
		handler = gziphandler.GzipHandler(handler)
	}

	// Inspect logging commandline arguments and setup the `ApacheLogger`:
	if s, err = nele.AppArguments.Get("accessLog"); (nil == err) && (0 < len(s)) {
		// we assume, an error means: no logfile
		if s2, err := nele.AppArguments.Get("errorLog"); (nil == err) && (0 < len(s2)) {
			handler = apachelogger.Wrap(handler, s, s2)
		} else {
			handler = apachelogger.Wrap(handler, s, "")
		}
		err = nil // for use by test for `apachelogger.SetErrLog()` (below)
	} else if s, err = nele.AppArguments.Get("errorLog"); (nil == err) && (0 < len(s)) {
		handler = apachelogger.Wrap(handler, "", s)
	}

	// We need a `server` reference to use it in `setupSinals()`
	// and to set some reasonable timeouts:
	server := &http.Server{
		Addr:              ph.Address(),
		Handler:           handler,
		IdleTimeout:       2 * time.Minute,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      5 * time.Minute,
	}
	setupSignals(server)
	if (nil == err) && (0 < len(s)) { // values from logfile test
		apachelogger.SetErrLog(server)
	}

	ck, _ = nele.AppArguments.Get("certKey")
	cp, _ = nele.AppArguments.Get("certPem")
	if 0 < len(ck) && (0 < len(cp)) {
		s = fmt.Sprintf("%s listening HTTPS at: %s", Me, ph.Address())
		log.Println(s)
		apachelogger.Log("Nele/main", s)
		if err = server.ListenAndServeTLS(cp, ck); nil != err {
			fatal(fmt.Sprintf("%s: %v", Me, err))
		}
		return
	}

	s = fmt.Sprintf("%s listening HTTP at: %s", Me, ph.Address())
	log.Println(s)
	apachelogger.Log("Nele/main", s)
	if err = server.ListenAndServe(); nil != err {
		fatal(fmt.Sprintf("%s: %v", Me, err))
	}
} // main()

/* _EoF_ */
