/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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
} // doAdd()

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

// `setupSinals()` configures the capture of the interrupts `SIGINT`,
// `SIGKILL`, and `SIGTERM` to terminate the program gracefully.
func setupSinals(aServer *http.Server) {
	// handle `CTRL-C` and `kill(9)` and `kill(15)`.
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		for signal := range c {
			log.Printf("%s captured '%v', stopping program and exiting ...", os.Args[0], signal)
			if err := aServer.Shutdown(context.Background()); nil != err {
				log.Fatalf("%s: %v", os.Args[0], err)
			}
		}
	}()
} // setupSinals()

// Actually run the program …
func main() {
	var (
		err       error
		handler   http.Handler
		ph        *nele.TPageHandler
		ck, cp, s string
	)
	Me, _ := filepath.Abs(os.Args[0])

	// Add a new posting via command line:
	doConsole(Me)

	// Read in a file:
	doFile(Me)

	if s, err = nele.AppArguments.Get("uf"); (nil == err) && (0 < len(s)) {
		fn := s
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
			nele.ListUser(fn)
		}
		if s, err = nele.AppArguments.Get("uu"); (nil == err) && (0 < len(s)) {
			nele.UpdateUser(s, fn)
		}
	}

	if ph, err = nele.NewPageHandler(); nil != err {
		nele.ShowHelp()
		log.Fatalf("%s: %v", Me, err)
	}
	handler = errorhandler.Wrap(ph, ph)

	// inspect `logfile` commandline argument and setup the `ApacheLogger`
	if s, err = nele.AppArguments.Get("logfile"); (nil == err) && (0 < len(s)) {
		// we assume, an error means: no logfile
		handler = apachelogger.Wrap(handler, s)
	}
	// We need a `server` reference to use it in setupSinals() below
	server := &http.Server{Addr: ph.Address(), Handler: handler}
	setupSinals(server)

	ck, _ = nele.AppArguments.Get("certKey")
	cp, _ = nele.AppArguments.Get("certPem")

	if 0 < len(ck) && (0 < len(cp)) {
		log.Printf("%s listening HTTPS at: %s", Me, ph.Address())
		if err = server.ListenAndServeTLS(cp, ck); nil != err {
			log.Fatalf("%s: %v", Me, err)
		}
		return
	}

	log.Printf("%s listening HTTP at: %s", Me, ph.Address())
	if err = server.ListenAndServe(); nil != err {
		log.Fatalf("%s: %v", Me, err)
	}
} // main()

/* _EoF_ */