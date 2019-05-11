module github.com/mwat56/go-blog

go 1.12

require (
	github.com/mwat56/go-apachelogger v1.0.0
	github.com/mwat56/go-errorhandler v1.0.0
	github.com/mwat56/go-ini v1.0.0
	github.com/mwat56/go-passlist v1.0.0
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190506204251-e1dfcc566284
	golang.org/x/tools v0.0.0-20190508150211-cf84161cff3f // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
