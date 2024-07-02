module github.com/mwat56/nele

go 1.22

toolchain go1.22.1

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.7.0
	github.com/mwat56/cssfs v0.2.7
	github.com/mwat56/errorhandler v1.1.11
	github.com/mwat56/hashtags v0.6.3
	github.com/mwat56/ini v1.9.0
	github.com/mwat56/jffs v0.1.4
	github.com/mwat56/passlist v1.3.11
	github.com/mwat56/screenshot v0.4.9
	github.com/mwat56/sourceerror v0.2.1
	github.com/mwat56/uploadhandler v1.1.11
	github.com/mwat56/whitespace v0.2.6
	github.com/russross/blackfriday/v2 v2.1.0
)

require (
	github.com/chromedp/cdproto v0.0.0-20240626232640-f933b107c653 // indirect
	github.com/chromedp/chromedp v0.9.5 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
)

replace (
	github.com/mwat56/apachelogger => ../apachelogger
	github.com/mwat56/cssfs => ../cssfs
	github.com/mwat56/errorhandler => ../errorhandler
	github.com/mwat56/hashtags => ../hashtags
	github.com/mwat56/ini => ../ini
	github.com/mwat56/jffs => ../jffs
	github.com/mwat56/passlist => ../passlist
	github.com/mwat56/screenshot => ../screenshot
	github.com/mwat56/sessions => ../sessions
	github.com/mwat56/uploadhandler => ../uploadhandler
	github.com/mwat56/whitespace => ../whitespace
)
