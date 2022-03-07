module github.com/mwat56/nele

go 1.17

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.6.1
	github.com/mwat56/cssfs v0.2.6
	github.com/mwat56/errorhandler v1.1.9
	github.com/mwat56/hashtags v0.6.1
	github.com/mwat56/ini v1.5.2
	github.com/mwat56/jffs v0.1.3
	github.com/mwat56/passlist v1.3.7
	github.com/mwat56/screenshot v0.4.3
	github.com/mwat56/uploadhandler v1.1.9
	github.com/mwat56/whitespace v0.2.3
	github.com/russross/blackfriday/v2 v2.1.0
)

require (
	github.com/chromedp/cdproto v0.0.0-20220304215434-892afa710589 // indirect
	github.com/chromedp/chromedp v0.7.8 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
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
