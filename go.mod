module github.com/mwat56/Nele

go 1.12

require (
	github.com/mwat56/apachelogger v1.1.1
	github.com/mwat56/errorhandler v1.0.2
	github.com/mwat56/hashtags v0.2.4
	github.com/mwat56/ini v1.0.1
	github.com/mwat56/passlist v1.0.1
	github.com/mwat56/uploadhandler v1.0.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
