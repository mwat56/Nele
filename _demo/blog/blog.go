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

	"github.com/mwat56/go-apachelogger"
	"github.com/mwat56/go-blog"
	"github.com/mwat56/go-errorhandler"
)

// `doConsole()` checks for the `add` commandline argument, adds the
// text from StdIn as a new post, and terminates the program.
func doConsole(aMe string) {
	var (
		err error
		i64 int64
		s   string
	)
	if s, err = blog.AppArguments.Get("pa"); nil != err {
		// we assume, an error means: no cmd line action
		return
	}
	if "true" == s {
		if i64, err = blog.AddConsolePost(); nil != err {
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
	if s, err = blog.AppArguments.Get("pf"); nil != err {
		// we assume, an error means: no cmd line action
		return
	}
	if 0 < len(s) {
		if i64, err = blog.AddFilePost(s); nil != err {
			log.Fatalf("%s: %v", aMe, err)
		}
		log.Printf("\n\t%s stored %d bytes in a new posting", aMe, i64)
		os.Exit(0)
	}
} // doFile()

// `setupSinals()` configures the capture of the interrupts `SIGINT`,
// `SIGKILL`, and `SIGTERM` to terminate the program gracefully.
func setupSinals(aMe string, aServer *http.Server) {
	// handle `CTRL-C` and `kill(9)` and `kill(15)`.
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		for signal := range c {
			log.Printf("%s captured '%v', stopping program and exiting ...", aMe, signal)
			if err := aServer.Shutdown(context.Background()); nil != err {
				log.Fatalf("%s: %v", aMe, err)
			}
		}
	}()
} // setupSinals()

// Actually run the program …
func main() {
	var (
		err       error
		handler   http.Handler
		ph        *blog.TPageHandler
		ck, cp, s string
	)
	Me, _ := filepath.Abs(os.Args[0])

	// Add a new posting via command line:
	doConsole(Me)

	// Read in a file:
	doFile(Me)

	if s, err = blog.AppArguments.Get("uf"); (nil == err) && (0 < len(s)) {
		fn := s
		if s, err = blog.AppArguments.Get("ua"); (nil == err) && (0 < len(s)) {
			blog.AddUser(s, fn)
		}
		if s, err = blog.AppArguments.Get("uc"); (nil == err) && (0 < len(s)) {
			blog.CheckUser(s, fn)
		}
		if s, err = blog.AppArguments.Get("ud"); (nil == err) && (0 < len(s)) {
			blog.DeleteUser(s, fn)
		}
		if s, err = blog.AppArguments.Get("ul"); (nil == err) && (0 < len(s)) {
			blog.ListUser(fn)
		}
		if s, err = blog.AppArguments.Get("uu"); (nil == err) && (0 < len(s)) {
			blog.UpdateUser(s, fn)
		}
	}

	if ph, err = blog.NewPageHandler(); nil != err {
		blog.ShowHelp()
		log.Fatalf("%s: %v", Me, err)
	}
	handler = errorhandler.Wrap(ph, ph)

	// inspect `logfile` commandline argument and setup the `ApacheLogger`
	if s, err = blog.AppArguments.Get("logfile"); (nil == err) && (0 < len(s)) {
		// we assume, an error means: no logfile
		handler = apachelogger.Wrap(handler, s)
	}
	// We need a `server` reference to use it in setupSinals() below
	server := &http.Server{Addr: ph.Address(), Handler: handler}

	ck, err = blog.AppArguments.Get("certKey")
	cp, err = blog.AppArguments.Get("certPem")

	setupSinals(Me, server)

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
