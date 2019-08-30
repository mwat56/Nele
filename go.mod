module github.com/mwat56/Nele

go 1.12

require (
	github.com/mwat56/apachelogger v1.2.5
	github.com/mwat56/errorhandler v1.1.0
	github.com/mwat56/hashtags v0.4.7
	github.com/mwat56/ini v1.3.4
	github.com/mwat56/passlist v1.1.2
	github.com/mwat56/uploadhandler v1.0.3
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sys v0.0.0-20190830142957-1e83adbbebd0 // indirect
	golang.org/x/tools v0.0.0-20190830082254-f340ed3ae274 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
