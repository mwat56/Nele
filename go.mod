module github.com/mwat56/nele

go 1.17

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.6.0
	github.com/mwat56/cssfs v0.2.6
	github.com/mwat56/errorhandler v1.1.9
	github.com/mwat56/hashtags v0.6.1
	github.com/mwat56/ini v1.5.2
	github.com/mwat56/jffs v0.1.3
	github.com/mwat56/pageview v0.4.6
	github.com/mwat56/passlist v1.3.6
	github.com/mwat56/uploadhandler v1.1.8
	github.com/mwat56/whitespace v0.2.3
	github.com/russross/blackfriday/v2 v2.1.0
)

require (
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
)

replace (
	github.com/mwat56/apachelogger => ../apachelogger
	github.com/mwat56/cssfs => ../cssfs
	github.com/mwat56/errorhandler => ../errorhandler
	github.com/mwat56/hashtags => ../hashtags
	github.com/mwat56/ini => ../ini
	github.com/mwat56/jffs => ../jffs
	github.com/mwat56/pageview => ../pageview
	github.com/mwat56/passlist => ../passlist
	github.com/mwat56/screenshot => ../screenshot
	github.com/mwat56/sessions => ../sessions
	github.com/mwat56/uploadhandler => ../uploadhandler
	github.com/mwat56/whitespace => ../whitespace
)
