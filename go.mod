module github.com/mwat56/nele

go 1.13

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.4.5
	github.com/mwat56/cssfs v0.1.2
	github.com/mwat56/errorhandler v1.1.5
	github.com/mwat56/hashtags v0.4.15
	github.com/mwat56/ini v1.3.8
	github.com/mwat56/jffs v0.0.5
	github.com/mwat56/pageview v0.4.2
	github.com/mwat56/passlist v1.2.1
	github.com/mwat56/uploadhandler v1.1.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
