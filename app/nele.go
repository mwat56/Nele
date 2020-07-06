/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package main

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
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
		i   int
		s   string
	)
	if !nele.AppArgs.PostAdd {
		// no cmd line action
		return
	}
	if "true" == s {
		if i, err = nele.AddConsolePost(); nil != err {
			log.Fatalf("%s: %v", aMe, err)
		}
		log.Printf("\n\t%s wrote %d bytes in a new posting", aMe, i)
		os.Exit(0)
	}
} // doConsole()

// `doFile()` checks for the `filename` commandline argument, adds the
// text from the file as a new post, and terminates the program.
func doFile(aMe string) {
	var (
		err error
		i   int
		s   string
	)
	if 0 == len(nele.AppArgs.PostFile) {
		// no posting file
		return
	}
	if i, err = nele.AddFilePost(s); nil != err {
		log.Fatalf("%s: %v", aMe, err)
	}
	log.Printf("\n\t%s stored %d bytes in a new posting", aMe, i)
	os.Exit(0)
} // doFile()

// `exit()` Log `aMessage` and terminate the program.
func exit(aMessage string) {
	apachelogger.Err("Nele/main", aMessage)
	runtime.Gosched() // let the logger write
	log.Fatalln(aMessage)
} // exit()

// `redirHTTP()` Send HTTP clients to HTTPS server.
//
// see: https://gist.github.com/d-schmidt/587ceec34ce1334a5e60
func redirHTTP(aWriter http.ResponseWriter, aRequest *http.Request) {
	// Copy the original URL and replace the scheme:
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
	http.Redirect(aWriter, aRequest, target, http.StatusTemporaryRedirect)
} // redirHTTP()

// `setupSignals()` configures the capture of the interrupts `SIGINT`
// `and `SIGTERM` to terminate the program gracefully.
func setupSignals(aServer *http.Server) {
	// handle `CTRL-C` and `kill(15)`.
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for signal := range c {
			msg := fmt.Sprintf("%s captured '%v', stopping program and exiting ...", filepath.Base(os.Args[0]), signal)
			apachelogger.Err(`Nele/catchSignals`, msg)
			log.Println(msg)
			break
		} // for

		ctx, cancel := context.WithCancel(context.Background())
		aServer.BaseContext = func(net.Listener) context.Context {
			return ctx
		}
		aServer.RegisterOnShutdown(cancel)

		ctxTimeout, cancelTimeout := context.WithTimeout(
			context.Background(), time.Second*10)
		defer cancelTimeout()
		if err := aServer.Shutdown(ctxTimeout); nil != err {
			exit(fmt.Sprintf("%s: %v", filepath.Base(os.Args[0]), err))
		}
	}()
} // setupSignals()

// `userCmdline()` checks for and executes password file commandline actions.
func userCmdline() {
	if 0 == len(nele.AppArgs.UserFile) {
		return // without file no user handling
	}

	// All the following `nele.UserXxx()` function calls will
	// terminate the program.
	if 0 < len(nele.AppArgs.UserAdd) {
		nele.UserAdd(nele.AppArgs.UserAdd, nele.AppArgs.UserFile)
	}
	if 0 < len(nele.AppArgs.UserCheck) {
		nele.UserCheck(nele.AppArgs.UserCheck, nele.AppArgs.UserFile)
	}
	if 0 < len(nele.AppArgs.UserDelete) {
		nele.UserDelete(nele.AppArgs.UserDelete, nele.AppArgs.UserFile)
	}
	if nele.AppArgs.UserList {
		nele.UserList(nele.AppArgs.UserFile)
	}
	if 0 < len(nele.AppArgs.UserUpdate) {
		nele.UserUpdate(nele.AppArgs.UserUpdate, nele.AppArgs.UserFile)
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
		exit(fmt.Sprintf("%s: %v", Me, err))
	}
	// Setup the errorpage handler:
	handler = errorhandler.Wrap(ph, ph)

	// Inspect `gzip` commandline argument and setup the Gzip handler:
	if nele.AppArgs.GZip {
		handler = gziphandler.GzipHandler(handler)
	}

	// Inspect logging commandline arguments and setup the `ApacheLogger`:
	handler = apachelogger.Wrap(handler, nele.AppArgs.AccessLog, nele.AppArgs.ErrorLog)

	ctxTimeout, cancelTimeout := context.WithTimeout(
		context.Background(), time.Second*10)
	defer cancelTimeout()

	// We need a `server` reference to use it in `setupSignals()`
	// and to set some reasonable timeouts:
	server := &http.Server{
		// The TCP address for the server to listen on:
		Addr: nele.AppArgs.Addr,
		// Return the base context for incoming requests on this server:
		BaseContext: func(net.Listener) context.Context {
			return ctxTimeout
		},
		// Request handler to invoke:
		Handler: handler,
		// Set timeouts so that a slow or malicious client
		// doesn't hold resources forever
		//
		// The maximum amount of time to wait for the next request;
		// if IdleTimeout is zero, the value of ReadTimeout is used:
		IdleTimeout: 0,
		// The amount of time allowed to read request headers:
		ReadHeaderTimeout: 10 * time.Second,
		// The maximum duration for reading the entire request,
		// including the body:
		ReadTimeout: 10 * time.Second,
		// The maximum duration before timing out writes of the response:
		WriteTimeout: 10 * time.Second,
	}
	if 0 < len(nele.AppArgs.ErrorLog) {
		apachelogger.SetErrLog(server)
	}
	setupSignals(server)

	if 0 < len(nele.AppArgs.CertKey) && (0 < len(nele.AppArgs.CertPem)) {
		// start the HTTP to HTTPS redirector:
		go http.ListenAndServe(nele.AppArgs.Addr, http.HandlerFunc(redirHTTP))

		// see:
		// https://ssl-config.mozilla.org/#server=golang&version=1.14.1&config=old&guideline=5.4
		server.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS10,
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, // #nosec G402
			},
		} // #nosec G402
		// server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))

		s := fmt.Sprintf("%s listening HTTPS at: %s", Me, nele.AppArgs.Addr)
		log.Println(s)
		apachelogger.Log("Nele/main", s)
		exit(fmt.Sprintf("%s: %v", Me,
			server.ListenAndServeTLS(nele.AppArgs.CertPem, nele.AppArgs.CertKey)))
		return
	}

	s := fmt.Sprintf("%s listening HTTP at: %s", Me, nele.AppArgs.Addr)
	log.Println(s)
	apachelogger.Log("Nele/main", s)
	exit(fmt.Sprintf("%s: %v", Me, server.ListenAndServe()))
} // main()

/* _EoF_ */
