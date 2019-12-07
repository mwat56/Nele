module github.com/mwat56/Nele

go 1.13

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.4.2
	github.com/mwat56/errorhandler v1.1.3
	github.com/mwat56/hashtags v0.4.12
	github.com/mwat56/ini v1.3.7
	github.com/mwat56/jffs v0.0.5
	github.com/mwat56/pageview v0.3.1
	github.com/mwat56/passlist v1.2.1
	github.com/mwat56/uploadhandler v1.0.8
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/sys v0.0.0-20191206220618-eeba5f6aabab // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1

replace (
	github.com/mwat56/apachelogger => ../apachelogger
	github.com/mwat56/errorhandler => ../errorhandler
	github.com/mwat56/hashtags => ../hashtags
	github.com/mwat56/ini => ../ini
	github.com/mwat56/jffs => ../jffs
	github.com/mwat56/pageview => ../pageview
	github.com/mwat56/passlist => ../passlist
	github.com/mwat56/sessions => ../sessions
	github.com/mwat56/uploadhandler => ../uploadhandler
)
