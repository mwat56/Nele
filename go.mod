module github.com/mwat56/nele

go 1.14

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.4.5
	github.com/mwat56/cssfs v0.2.1
	github.com/mwat56/errorhandler v1.1.5
	github.com/mwat56/hashtags v0.5.1
	github.com/mwat56/ini v1.3.8
	github.com/mwat56/jffs v0.1.0
	github.com/mwat56/pageview v0.4.3
	github.com/mwat56/passlist v1.3.1
	github.com/mwat56/uploadhandler v1.1.2
	github.com/mwat56/whitespace v0.2.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 // indirect
	golang.org/x/sys v0.0.0-20200331124033-c3d80250170d // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
