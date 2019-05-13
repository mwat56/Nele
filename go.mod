module github.com/mwat56/go-blog

go 1.12

require (
	github.com/mwat56/go-apachelogger v1.0.0
	github.com/mwat56/go-errorhandler v1.0.0
	github.com/mwat56/go-hashtags v0.0.0-20190513192225-ce96876e346b
	github.com/mwat56/go-ini v1.0.0
	github.com/mwat56/go-passlist v1.0.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	golang.org/x/tools v0.0.0-20190513184735-d81a07b7e584 // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
