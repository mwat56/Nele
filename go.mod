module github.com/mwat56/nele

go 1.13

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.4.4
	github.com/mwat56/errorhandler v1.1.4
	github.com/mwat56/hashtags v0.4.13
	github.com/mwat56/ini v1.3.7
	github.com/mwat56/jffs v0.0.5
	github.com/mwat56/pageview v0.3.2
	github.com/mwat56/passlist v1.2.1
	github.com/mwat56/uploadhandler v1.0.9
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
