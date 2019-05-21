module github.com/mwat56/go-blog

go 1.12

require (
	github.com/mwat56/go-apachelogger v1.1.0
	github.com/mwat56/go-errorhandler v1.0.0
	github.com/mwat56/go-hashtags v0.0.6
	github.com/mwat56/go-ini v1.0.0
	github.com/mwat56/go-passlist v1.0.0
	github.com/mwat56/go-uploadhandler v0.2.1
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
