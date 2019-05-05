module github.com/mwat56/go-blog

go 1.12

require (
	github.com/mwat56/go-apachelogger v0.0.0-20190504155952-32112b38fab9
	github.com/mwat56/go-errorhandler v0.0.0-20190504160024-b256f309d41b
	github.com/mwat56/go-ini v0.0.0-20190504160039-e4168b397c01
	github.com/mwat56/go-passlist v0.0.0-20190505110607-f1be27ce6ce3
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000
)

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
