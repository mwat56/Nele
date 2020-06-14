/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
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
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/mwat56/apachelogger"
	"github.com/mwat56/errorhandler"
	"github.com/mwat56/nele"
)

// `doConsole()` checks for the `add` commandline argument, adds the
// text from StdIn as a new post, and terminates the program.
func doConsole(aMe string) {
	var (
		err error
		i64 int64
		s   string
	)
	if !nele.AppArgs.PostAdd {
		// no cmd line action
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
	if 0 == len(nele.AppArgs.PostFile) {
		// no posting file
		return
	}
	if i64, err = nele.AddFilePost(s); nil != err {
		log.Fatalf("%s: %v", aMe, err)
	}
	log.Printf("\n\t%s stored %d bytes in a new posting", aMe, i64)
	os.Exit(0)
} // doFile()

// `fatal()` logs `aMessage` and terminates the program.
func fatal(aMessage string) {
	apachelogger.Err("Nele/main", aMessage)
	runtime.Gosched() // let the logger write
	apachelogger.Close()
	log.Fatalln(aMessage)
} // fatal()

// `redirHTTP()` sends HTTP clients to HTTPS server.
//
// see: https://gist.github.com/d-schmidt/587ceec34ce1334a5e60
func redirHTTP(aWriter http.ResponseWriter, aRequest *http.Request) {
	// copy the original URL and replace the scheme:
	targetURL := url.URL{
		Scheme:     `https`,
		Opaque:     aRequest.URL.Opaque,
		User:       aRequest.URL.User,
		Host:       aRequest.URL.Host,
		Path:       aRequest.URL.Path,
		RawPath:    aRequest.URL.RawPath,
		ForceQuery: aRequest.URL.ForceQuery,
		RawQuery:   aRequest.URL.RawQuery,
		Fragment:   aRequest.URL.Fragment,
	}
	target := targetURL.String()

	apachelogger.Err(`Nele/main`, `redirecting to: `+target)
	http.Redirect(aWriter, aRequest, target, http.StatusMovedPermanently)
} // redirHTTP()

// `setupSignals()` configures the capture of the interrupts `SIGINT`
// `and `SIGTERM` to terminate the program gracefully.
func setupSignals(aServer *http.Server) {
	// handle `CTRL-C` and `kill(15)`.
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	catcher := func() {
		for signal := range c {
			msg := fmt.Sprintf("%s captured '%v', stopping program and exiting ...", os.Args[0], signal)
			apachelogger.Err(`Nele/catchSignals`, msg)
			log.Println(msg)
			runtime.Gosched() // let the logger write
			if err := aServer.Shutdown(context.Background()); nil != err {
				fatal(fmt.Sprintf("%s: %v", os.Args[0], err))
			}
		}
	} // catcher()

	go catcher()
} // setupSignals()

// `userCmdline()` checks for and executes password file commandline actions.
func userCmdline() {
	if 0 == len(nele.AppArgs.UserFile) {
		return // without file no user handling
	}
	// All the following `nele.UserXxx` function calls will
	// terminate the program.
	if 0 < len(nele.AppArgs.UserAdd) {
		nele.AddUser(nele.AppArgs.UserAdd, nele.AppArgs.UserFile)
	}
	if 0 < len(nele.AppArgs.UserCheck) {
		nele.CheckUser(nele.AppArgs.UserCheck, nele.AppArgs.UserFile)
	}
	if 0 < len(nele.AppArgs.UserDelete) {
		nele.DeleteUser(nele.AppArgs.UserDelete, nele.AppArgs.UserFile)
	}
	if nele.AppArgs.UserList {
		nele.ListUsers(nele.AppArgs.UserFile)
	}
	if 0 < len(nele.AppArgs.UserUpdate) {
		nele.UpdateUser(nele.AppArgs.UserUpdate, nele.AppArgs.UserFile)
	}
} // userCmdline()

// Actually run the program …
func main() {
	var (
		err     error
		handler http.Handler
		ph      *nele.TPageHandler
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
	if nele.AppArgs.GZip {
		handler = gziphandler.GzipHandler(handler)
	}

	// Inspect logging commandline arguments and setup the `ApacheLogger`:
	if 0 < len(nele.AppArgs.AccessLog) {
		// we assume, an error means: no logfile
		if 0 < len(nele.AppArgs.ErrorLog) {
			handler = apachelogger.Wrap(handler, nele.AppArgs.AccessLog, nele.AppArgs.ErrorLog)
		} else {
			handler = apachelogger.Wrap(handler, nele.AppArgs.AccessLog, ``)
		}
		// err = nil // for use by test for `apachelogger.SetErrLog()` (below)
	} else if 0 < len(nele.AppArgs.ErrorLog) {
		handler = apachelogger.Wrap(handler, ``, nele.AppArgs.ErrorLog)
	}

	// We need a `server` reference to use it in `setupSignals()`
	// and to set some reasonable timeouts:
	server := &http.Server{
		Addr:    nele.AppArgs.Addr,
		Handler: handler,
		// Set timeouts so that a slow or malicious client
		// doesn't hold resources forever
		IdleTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	if 0 < len(nele.AppArgs.ErrorLog) {
		apachelogger.SetErrLog(server)
	}
	setupSignals(server)

	if 0 < len(nele.AppArgs.CertKey) && (0 < len(nele.AppArgs.CertPem)) {
		// start the HTTP to HTTPS redirector:
		go http.ListenAndServe(nele.AppArgs.Addr, http.HandlerFunc(redirHTTP))

		s := fmt.Sprintf("%s listening HTTPS at: %s", Me, nele.AppArgs.Addr)
		log.Println(s)
		apachelogger.Log("Nele/main", s)
		fatal(fmt.Sprintf("%s: %v", Me,
			server.ListenAndServeTLS(nele.AppArgs.CertPem, nele.AppArgs.CertKey)))
		return
	}

	s := fmt.Sprintf("%s listening HTTP at: %s", Me, nele.AppArgs.Addr)
	log.Println(s)
	apachelogger.Log("Nele/main", s)
	fatal(fmt.Sprintf("%s: %v", Me, server.ListenAndServe()))
} // main()

/* _EoF_ */
