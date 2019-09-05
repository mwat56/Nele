module github.com/mwat56/Nele

go 1.13

require (
	github.com/mwat56/apachelogger v1.2.6
	github.com/mwat56/errorhandler v1.1.1
	github.com/mwat56/hashtags v0.4.8
	github.com/mwat56/ini v1.3.5
	github.com/mwat56/passlist v1.1.3
	github.com/mwat56/uploadhandler v1.0.4
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sys v0.0.0-20190904154756-749cb33beabd // indirect
	golang.org/x/tools v0.0.0-20190903025054-afe7f8212f0d // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
