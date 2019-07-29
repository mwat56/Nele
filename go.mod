module github.com/mwat56/Nele

go 1.12

require (
	github.com/mwat56/apachelogger v1.2.2
	github.com/mwat56/errorhandler v1.0.3
	github.com/mwat56/hashtags v0.4.4
	github.com/mwat56/ini v1.3.1
	github.com/mwat56/passlist v1.1.1
	github.com/mwat56/uploadhandler v1.0.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/sys v0.0.0-20190726091711-fc99dfbffb4e // indirect
	golang.org/x/tools v0.0.0-20190729092621-ff9f1409240a // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
