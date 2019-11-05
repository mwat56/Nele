module github.com/mwat56/Nele

go 1.13

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/mwat56/apachelogger v1.2.9
	github.com/mwat56/errorhandler v1.1.1
	github.com/mwat56/hashtags v0.4.9
	github.com/mwat56/ini v1.3.6
	github.com/mwat56/jffs v0.0.3
	github.com/mwat56/passlist v1.2.0
	github.com/mwat56/uploadhandler v1.0.5
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20191105034135-c7e5f84aec59 // indirect
	golang.org/x/sys v0.0.0-20191104094858-e8c54fb511f6 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
