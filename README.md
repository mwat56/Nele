# Nele Blog

[![Golang](https://img.shields.io/badge/Language-Go-green.svg)](https://golang.org/)
[![GoDoc](https://godoc.org/github.com/mwat56/nele?status.svg)](https://godoc.org/github.com/mwat56/nele/)
[![Go Report](https://goreportcard.com/badge/github.com/mwat56/nele)](https://goreportcard.com/report/github.com/mwat56/nele)
[![Issues](https://img.shields.io/github/issues/mwat56/nele.svg)](https://github.com/mwat56/nele/issues?q=is%3Aopen+is%3Aissue)
[![Size](https://img.shields.io/github/repo-size/mwat56/nele.svg)](https://github.com/mwat56/nele/)
[![Tag](https://img.shields.io/github/tag/mwat56/nele.svg)](https://github.com/mwat56/nele/tags)
[![License](https://img.shields.io/github/license/mwat56/nele.svg)](https://github.com/mwat56/nele/blob/main/LICENSE)
[![View examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg)](https://github.com/mwat56/nele/blob/main/app/nele.go)

- [Nele Blog](#nele-blog)
	- [Purpose](#purpose)
	- [Features](#features)
	- [Installation](#installation)
	- [Usage](#usage)
		- [Commandline postings](#commandline-postings)
		- [Authentication](#authentication)
		- [User/password file & handling](#userpassword-file--handling)
		- [Page/link previews](#pagelink-previews)
	- [Configuration](#configuration)
	- [URLs](#urls)
	- [Files](#files)
		- [CSS](#css)
		- [Fonts](#fonts)
		- [Images](#images)
		- [Postings](#postings)
		- [Static](#static)
		- [Views](#views)
		- [Contents](#contents)
	- [Libraries](#libraries)
	- [Licence](#licence)

----

## Purpose

The purpose of this package/application was twofold initially. On one hand I needed a project to learn the (then to me new) `Go` language, and on the other hand I wanted a project, that lead me into different domains, like user authentication, configuration, data formats, error handling, filesystem access, data logging, os, network, regex, templating etc. –
And, I wanted no external dependencies (like databases etc.). –
And, I didn't care for Windows<sup>(tm) </sup> compatibility since I left the MS-platform about 25 years ago after using it in the 80s and early 90s of the last century.
(But who, in his right mind, would want to run a web-service on such a platform anyway?)

That's how I ended up with this little blog-system (for lack of a better word; or: "diary", "notes", …).
It's a system that lets you write and add articles from both the command line and a web-interface.
It provides options to add, modify and delete entries using a user/password list for authentication when accessing certain URLs in this system.
Articles can be added, edited (e.g. for correcting typos etc.), or removed altogether.
If you don't like the styles coming with the package you can, of course, change them according to your preferences in your own installation.

The articles you write are then available on the net as _web-pages_.

It is not, however, a discussion platform.
It's supposed to be used as a publication platform, not some kind of _social media_ (where people thrive on insulting each other).
So I intentionally didn't bother with comments or discussion threading.

## Features

* Markdown support
* Multiple user accounts supported
* No database (like SQLite, MariaDB, etc.) required
* No JavaScript dependency
* No cookies needed
* Privacy aware
* Simplicity of use

## Installation

You can use `Go` to install this package for you:

    go get -u github.com/mwat56/nele

## Usage

After downloading this package you go to its directory and compile

    go build app/nele.go

which should produce an executable binary.
On my system it looked (at a certain point in time) like this:

    $ ls -l
	total 11420
	drwxrwxr-x 12 matthias matthias     4096 Mai 23 18:35 .
	drwxrwxr-x 12 matthias matthias     4096 Mai 23 17:58 ..
	-rw-rw-r--  1 matthias matthias      474 Apr 27 00:21 addTest.md
	drwxrwxr-x  3 matthias matthias     4096 Mai 23 18:14 app
	drwxrwxr-x  2 matthias matthias     4096 Mai 23 18:14 certs
	-rw-rw-r--  1 matthias matthias     6583 Mai 23 18:14 cmdline.go
	-rw-rw-r--  1 matthias matthias    10149 Mai 23 18:20 config.go
	-rw-rw-r--  1 matthias matthias     1846 Mai 23 18:14 config_test.go
	drwxrwxr-x  2 matthias matthias     4096 Mai 23 18:14 css
	-rw-rw-r--  1 matthias matthias      823 Mai 23 18:14 doc.go
	drwxrwxr-x  2 matthias matthias     4096 Mai 23 18:14 fonts
	drwxrwxr-x  8 matthias matthias     4096 Mai 23 18:10 .git
	drwxrwxr-x  3 matthias matthias     4096 Mai 23 17:58 .github
	-rw-rw-r--  1 matthias matthias      123 Mai 23 17:58 .gitignore
	-rw-------  1 matthias matthias      507 Mai 23 18:11 go.mod
	-rw-------  1 matthias matthias     4004 Mai 23 18:11 go.sum
	-rw-rw-r--  1 matthias matthias     5010 Mai 23 18:18 hashfile.db
	drwxrwxr-x  2 matthias matthias     4096 Mai 23 18:14 img
	-rw-rw-r--  1 matthias matthias    32474 Mai 23 17:58 LICENSE
	-rwxrwxr-x  1 matthias matthias 11149115 Mai 23 18:19 nele
	-rw-rw-r--  1 matthias matthias    21803 Mai 23 18:14 pagehandler.go
	-rw-rw-r--  1 matthias matthias      619 Mai 23 18:22 pagehandler_test.go
	-rw-rw-r--  1 matthias matthias     9313 Mai 23 18:14 posting.go
	drwxrwxr-x  8 matthias matthias     4096 Mai 23 18:01 postings
	-rw-rw-r--  1 matthias matthias    15319 Mai 23 18:14 posting_test.go
	-rw-rw-r--  1 matthias matthias     8240 Mai 23 18:14 postlist.go
	-rw-rw-r--  1 matthias matthias     7279 Mai 23 18:23 postlist_test.go
	-rw-rw-r--  1 matthias matthias       70 Mai 23 18:19 pwaccess.db
	-rw-rw-r--  1 matthias matthias    22792 Mai 23 18:35 README.md
	-rw-rw-r--  1 matthias matthias    10435 Mai 23 18:24 regex.go
	-rw-rw-r--  1 matthias matthias     8190 Mai 23 18:14 regex_test.go
	drwxrwxr-x  2 matthias matthias     4096 Mai 23 17:58 static
	-rw-rw-r--  1 matthias matthias     3656 Mai 23 18:14 tags.go
	-rw-rw-r--  1 matthias matthias     3811 Mai 23 17:58 template_vars.md
	-rw-rw-r--  2 matthias matthias     3109 Mai 23 18:16 TODO.md
	drwxrwxr-x  3 matthias matthias     4096 Mai 23 18:14 views
	-rw-rw-r--  1 matthias matthias     6787 Mai 23 18:14 views.go
	-rw-rw-r--  1 matthias matthias     6009 Mai 23 18:14 views_test.go
    $ _

You can reduce the binary's size by stripping it:

    $ strip nele
    $ ls -l nele
    -rwxrwxr-x 1 matthias matthias 8146912 Mai 23 18:38 nele
    $ _

As you can see the binary lost about 3MB of its weight.
Or you can further [compress](https://dev.to/akshaybharambe14/guide-to-compress-golang-binaries-2d95) your Go binary by using [UPX](https://upx.github.io/) which would downsize the `Nele` binary to about 2.5MB.

Let's start with the command line:

	$ ./nele -h

	Usage: ./nele [OPTIONS]

	-accessLog string
		<filename> Name of the access logfile to write to
		(default "/home/matthias/nele/access.bla.mwat.de")
	-blogName string
		<string> Name of this Blog (shown on every page)
		(default "Meine Güte, was für'n Blah!")
	-certKey string
		<fileName> Name of the TLS certificate key
	-certPem string
		<fileName> Name of the TLS certificate PEM
	-d	dump
	-dataDir string
		<dirName> Directory with CSS, FONTS, IMG, SESSIONS, and VIEWS sub-directories
		(default "/home/matthias/nele")
	-delWhitespace
		<boolean> Delete superfluous whitespace in generated pages (default true)
	-errorlog string
		<filename> Name of the error logfile to write to
		(default "/home/matthias/error.bla.mwat.de")
	-gzip
		<boolean> use gzip compression for server responses (default true)
	-hashFile string
		<fileName> Name of the file storing #hashtags and @mentions
		(default "/home/matthias/nele/hashfile.db")
	-ini string
		<fileName> the path/filename of the INI file to use
		(default "/home/matthias/.nele.ini")
	-lang string
		<de|en> Default language to use  (default "de")
	-listen string
		<IP number> The host's IP to listen at  (default "127.0.0.1")
	-lst
		<boolean> Log a stack trace for recovered runtime errors  (default true)
	-mfs string
		<filesize> Max. accepted size of uploaded files (default "10mb")
	-pa
		<boolean> (optional) posting add: write a posting from the commandline
	-pf string
		<fileName> (optional) post file: name of a file to add as new posting
	-port int
		<port number> The IP port to listen to  (default 8181)
	-pv
		<boolean> Use page preview images for links (default true)
	-realm string
		<hostName> Name of host/domain to secure by BasicAuth
		(default "Matthias' Bla")
	-theme string
		<name> The display theme to use ('light' or 'dark')
		(default "dark")
	-ua string
		<userName> User add: add a username to the password file
	-uc string
		<userName> User check: check a username in the password file
	-ud string
		<userName> User delete: remove a username from the password file
	-uf string
		<fileName> Passwords file storing user/passwords for BasicAuth
		(default "/home/matthias/nele/pwaccess.db")
	-ul
		<boolean> User list: show all users in the password file
	-uu string
		<userName> User update: update a username in the password file

	Most options can be set in an INI file to keep the command-line short ;-)

	$ _

> Please _note_ that the default values shown above vary depending on the system where `nele` is called and especially on the contents of the INI file (it's read before the help text is produced).

To just run the program, however, you'll usually don't need any of those options to input on the commandline.
There is an INI file called `nele.ini` coming with the package, where you can store the most common settings:

	$ cat nele.ini
	# Nele's default configuration file

	[Default]

	# Name of the optional logfile to write to.
	# NOTE: A relative path/name will be combined with `datadir` (below).
	accessLog = ./access.log

	# Name of this Blog (shown on every page).
	blogName = "Meine Güte, was für'n Blah!"

	# path-/filename of the TLS certificate's private key to enable
	# TLS/HTTPS (if empty standard HTTP is used).
	# NOTE: A relative path/name will be combined with `datadir` (below).
	certKey = ./certs/server.key

	# path-/filename of TLS (server) certificate to enable TLS/HTTPS
	# (if empty standard HTTP is used).
	# NOTE: A relative path/name will be combined with `datadir` (below).
	certPem = ./certs/server.pem

	# The directory root for the "css", "fonts", "img", "postings",
	# "static", and "views" sub-directories.
	# NOTE: This should be an _absolute_ path name.
	dataDir = ./

	# Delete superfluous whitespace in generated pages.
	delWhitespace = yes

	# Name of the optional logfile to write to.
	# NOTE: A relative path/name will be combined with `datadir` (above).
	errorLog =  ./error.log

	# Use gzip compression for server responses.
	gzip = true

	# The file to store #hashtags and @mentions.
	# NOTE: A relative path/name will be combined with `datadir` (above).
	hashFile = ./hashfile.db

	# The default UI language to use ("de" or "en").
	lang = de

	# The host's IP number to listen at.
	# An empty value means: listen on all interfaces.
	listen = 127.0.0.1

	# Whether or not log a stack trace for recovered runtime errors.
	# NOTE: This is merely a debugging aid and should normally be `false`.
	logStack = true

	# The IP port to listen to.
	port = 8181

	# Accepted size of uploaded files.
	maxfilesize = 10MB

	# Password file for HTTP Basic Authentication.
	# NOTE: a relative path/name will be combined with `datadir` (above).
	passFile = ./pwaccess.db

	# Name of host/domain to secure by BasicAuth.
	realm = "This Host"

	# Use screenshot images of linked pages.
	# NOTE: This feature depends on the external `wkhtmltoimage` binary;
	# for more details see: https://godoc.org/github.com/mwat56/screenshot
	Screenshot = true

	# Web/display theme ("dark" or "light").
	theme = dark

	# _EoF_
	$ _

The program, when started, will first look for the INI file in five different places:

1. in your (i.e. the current user's) directory (`./nele.ini`),
2. in the computer's main config directory (`/etc/nele.ini"`),
3. in the current user's home directory (e.g. `$HOME/.nele.ini`),
4. in the current user's configuration directory (e.g. `$HOME/.config/nele.ini`),
5. in the `-ini <filename>` commandline option (if given).

All these files (_if they exist_) are read in the given order at startup before finally parsing the commandline options shown above.
So each step overwrites the previous one, the commandline options having the highest priority. –
But let's look at some of the commandline options more closely.

### Commandline postings

You can post an article directly from the commandline.

`./nele -pa` allows you to write an article/posting directly on the commandline.

    $ ./nele -pa
    This is
    a test
    posting directly
    from the commandline.
    <Ctrl-D>
    2019/05/06 14:57:30 ./nele wrote 54 bytes in a new posting
    $ _

`./nele -pf <fileName>` allows you to include an already existing text file (with possibly some Markdown markup) into the system.

    $ ./nele -pf addTest.md
    2019/05/06 15:09:27 ./nele stored 474 bytes in a new posting
    $ _

These two options (`-pa` and `-pf`) are only usable from the commandline.

### Authentication

Why, you may ask, would you need an username/password file anyway? Well, you remember me mentioning that you can add, edit and delete articles? You wouldn't want _anyone_ on the net being able to do that, now, would you? For that reason, whenever there's no password file given (either in the INI file or the command-line) all functionality requiring authentication will be _disabled_. (Better safe than sorry, right?)

> _Note_ that the password file generated and used by this system resembles the `htpasswd` used by the Apache web-server, but both files are _not_ interchangeable because the actual encryption algorithms used by both respectively are different.

### User/password file & handling

Only usable from the commandline are the `-uXX` options, most of which need a username and the name of the password file to use.
> _Note_ that whenever you're prompted to input a password this will _not_ be echoed to the console.

The `-ua` option allows you to add an user/password pair:

	$ ./nele -ua testuser1 -uf pwaccess.db

	   password:
	  repeat pw:
		added 'testuser1' to list
	$ _

Again: The password input is not echoed to the console, therefore you don't see it.

Since we have the `passfile` setting already in our INI file (see above) we can forget the `-uf` option for the next options.

With `-uc` you can check a user's password:

	$ ./nele -uc testuser1

	password:
		'testuser1' password check successful
	$ _

This `-uc` you'll probably never actually use, it was just easy to implement.

If you want to remove an user account the `-ud` will do the trick (i.e. delete a user):

	$ ./nele -ud testuser1
		removed 'testuser1' from list
	$ _

When you want to know which users are stored in your password file `-ul` is your friend:

	$ ./nele -ul
	matthias

	$ _

Since we deleted the `testuser1` before only one entry remains.

That only leaves `-uu` to update (change) a user's password.

	$ ./nele -ua testuser2

	password:
	repeat pw:
		added 'testuser2' to list

	$ ./nele -uu testuser2

	   password:
	  repeat pw:
		updated user 'testuser2' in list

	$ ./nele -ul
	matthias
	testuser2

	$ _

First we added (`-ua`) a new user, then we updated the password (`-uu`), and finally we asked for the list of users (`-ul`).

### Page/link previews

If you set the `Screenshot` INI- or commandline-option to `true` there will be a preview image generated – by way of calling the [ChromeDP](https://github.com/chromedp/chromedp) library.
Those image files are stored locally (in the `./img/` directory) and may be used as often as you want.

> **Note** that screenshot images are created only for links in a _blockquote_ section:
>
>	`> [link text](http://www.example.org/one.html)`
>
> will be changed to
>
>	`> [![alt text](/httpwwwexampleorgonehtml.png)](http://www.example.org/one.html)`
>
> while
>
>	`bla [link text](http://www.example.org/one.html) bla`
>
> will be left untouched as a normal hyperlink.
>
> This restriction was introduced to avoid messing up the overall layout of a posting: It wouldn't look good if every link in a sentence would be replaced by an image.

The Go library controlling a headless instance of the `Chrome` browser is [ChromeDP](https://github.com/chromedp/chromedp) and is  _required_  for this package to work. Under Linux this browser is usually part of your distribution (as `chromium-browser`).

Generating a screenshot image usually takes between one and five seconds, depending on the actual web-page in question, bandwidth, traffic etc.; however, it can take considerably longer. To avoid hanging the program the `CreateImage()` function uses a timeout of half a minute by default.

And, finally, not all web-pages can be rendered properly and turned into an image. In such a case `ChromeDP` usually aborts with an error and the link in your posting just remains as is (i.e. a normal text link w/o preview/screenshot).

## Configuration

The system's configuration takes two steps:

1. Prepare the required files and directories.
2. Customise the INI file and/or prepare a script with all needed commandline arguments.
3. You most probably want to customise the files [./views/imprint.gohtml](https://github.com/mwat56/nele/blob/main/views/imprint.gohtml#L17), [./views/licence.gohtml](https://github.com/mwat56/nele/blob/main/views/licence.gohtml#L20), and [./views/privacy.gohtml](https://github.com/mwat56/nele/blob/main/views/privacy.gohtml#L17) according to your personal requirements.

## URLs

The system uses a number of slightly different URL groups.

### Static URLs

First, there are the static files served from the `css`, `img`, and `static` directories.
The actual location of which you can configure with the `datadir` INI entry and/or commandline option.

### Common URLs

Second, there are the URLs any normal user might see and use:

* `/` defines the logical root of the presentation; it's effectively the same as `/n/` (see below).
* `/faq`, `/imprint`, `/licence`, and `/privacy` serve static files which have to be filled with content according to your personal and legal needs.
* `/hl/tagname` allows the users to search for `#tagname` (but you'll input it without the number sign `#` because that has a special meaning in an URL).
Provided the given `tagname` was actually used in one or more of your articles a list of the respective postings will be shown.
* `/m/` shows the articles of the current month.
One can, however, specify the month one is interested in by adding a data part defining the month one wants to see (`/m/yyyy-mm`), like `/m/2019-04` to see the articles from April 2019.
* `/ml/mentionedname` allows the users to search for `@mentionedname` (but one will input it without the at sign `@` because that has a special meaning in an URL).
Provided the given `mentionedname` was actually used in one or more of your articles a list of the respective articles will be shown.
* `/n/` gives you the newest 30 articles.
The number of articles to show can be added to the URL like `/n/5` to see only five articles, or `/n/100` to see a hundred.
If one want to see the articles in slices of, say, 10 per page (instead of the default 30/page) one can use the URL `/n/10,10` and to see the second slice use `/n/10,20`, the third with `/n/10,30` and so on.
However, as long as there are more articles available, there will be a `»»` link at the bottom of the page to ease the navigation for the reader.
* `/p/1234567890abcdef` shows a single article/posting (the ID is automatically generated).
This kind of URL your users will see when they choose on another page to see the single article per page by selecting the leading `[*]` link in the overview page(s).
* `/s/searchterm` can be used to search for articles containing a certain word or expression.
All existing articles will be searched for the given `searchterm`.
* `/w/` shows the articles of the current week.
One can, however, specify the week one is interested in by adding a data part defining the week to see (`/w/yyyy-mm-dd`), like `/w/2019-04-13` to see the articles from the week in April 2019 containing the 13th.

### Internal URLs

And, third, there's a group of URLs your users won't see or use, because by design they are reserved for you, the author of your postings.
These URLs are protected by an authentication mechanism called _BasicAuth_ (which is supported by browsers for at least twenty years); this is where the username/password file comes in.
Only users whose credentials (i.e. username and password) are stored in the password file will be given access to the following URLs.
_So don't forget to set up an appropriate password file_.
If you forget that (or the file is not accessible for the program) everybody on the net could read, modify, or delete your articles, or add new ones – which you might not like; therefore the system disables all options that might modify your system.

* `/ap/` add a new posting.
A simple Web form will allow you to input whatever is on your mind.
* `/dp/234567890abcdef1` lets you change an article/posting's _date/time_ if you feel the need for cosmetic or other reasons.
Since you don't usually know/remember the article ID you'll first go to show the article/posting on a single page (`/p/234567890abcdef1`) by selecting the respective `[*]` link on the index page and then just prepend the `p` by a `d` in the URL.
* `/ep/34567890abcdef12` lets you edit the article/posting's _text_ identified by `34567890abcdef12`, e.g. to fix typos or correct the grammar.
* `/il/` (init list): Assuming you configured the `hashfile` INI-/commandline-option this shows you a simple HTML form by which you can start a background process re-initialising the hashlist.
It clears the current list and reads all postings to extract the #hashtags and @mentions.
_Note_: You will barely (if ever) need this option; it's mostly a debugging aid.
* `/pv/` (preview): Assuming you set the `Screenshot` INI-/commandline-option to `true` this shows you a simple HTML form by which you can start a background process checking all postings for page preview/screenshot images.
Again, this was implemented as a debugging aid and you won't usually use this option.
* `/rp/4567890abcdef123` lets you remove (delete) the article/posting identified by `4567890abcdef123` altogether.
_Note_ that there's no `undo` feature: Once you've deleted an article/posting it's gone.
* `/share/https://some.host.domain/somepage` lets you share another page URL.
Whatever you write after the initial `/share/` is assumed to be a remote URL, and a new article will be created and shown for you to edit.
* `/si/` (store image): This shows you a simple HTML form by which you can upload image files into your `/img/` directory.
Once the upload is done you (i.e. the user) will be presented an edit page in which the uploaded image is used.
* `/ss/` (store static): This shows you a simple HTML form by which you can upload static files into your `/static/` directory.
Once the upload is done you (i.e. the user) will be presented an edit page in which the uploaded file is used.
* `/xt/` (eXchange tag): This shows you a simple HTML form by which you can exchange a `#hashtag`/`@mention` with another one, or correct its writing.
_Note_ that the search for the term to replace is done case-insensitive while the replacement string gets inserted as you write it.

## Files

Right at the start I mentioned that I wanted to avoid external dependencies – like databases for example.
Well, that's not exactly true (or even possible), because there is _one_ "database" that's always already there, regardless of the operating system: the _filesystem_.
The trick is to figure out how to best use it for our own purposes.
The solution I came up with here is to use sort of a _timestamp_ as ID and filename for the articles, and use part of that very timestamp as ID and name for the directory names as well.

Both directory- and file-names are automatically handled by the system.
Each directory can hold up to 52 days worth of articles.
After extensive experimentation – with hundreds of thousands of automatically generated (and deleted) test files – that number seemed to be a reasonable compromise between directories not growing too big (_search times_) and keeping the number of directories used low (about seven per year).

All this data (files and directories) will be created under the directory you configure either in the INI file (entry `datadir`) or on the commandline (option `-datadir`).
Under that directory the program expects several sub-directories:

* `css/` for stylesheet files,
* `fonts/` for font files,
* `img/` for image files,
* `postings/` directory root for the articles,
* `static/` for static files (like e.g. PDF files),
* `views/` for page templates

Apart from setting that `datadir` option to your liking you don't have to worry about it any more.

As mentioned before, it's always advisable to use _absolute pathnames_, not relative one.
The latter are converted into absolute ones (based on `datadir`) by the system, but they depend on where you are in the filesystem when you start the program or write the commandline options.
You can use `./nele -h` to see which directories the program will use (see the example above).

### CSS

In the CSS directory (`datadir`/`css`) there are currently four files that are used automatically (i.a. hardcoded) by the system: `stylesheet.css` with some basic styling rules and `dark.css` and `light.css` with different settings for mainly colours, thus implementing two different _themes_ for the web-presentation, and there's the `fonts.css` file setting up the custom fonts to use.
The `theme` INI setting and the `-theme` commandline option determine which of the two `dark` and `light` styles to actually use.

### Fonts

The `datadir`/`fonts/` directory contains some freely available fonts used by the CSS files.

### Images

The `datadir`/`img/` directory can be used to store, well, _images_ to which you then can link in your articles.
You can put there whatever images you like either from the command-line or by using the system's `/si` URL.

Additionally any page preview/screenshot images are stored here (if set the `Screenshot` INI- or commandline-option to `true`).

### Postings

The `datadir/` directory is the base for storing all the articles.
The system creates subdirectories as needed to store new articles.
This directory structure is not accessed via a direct URL but used internally by the system.

### Static

The `datadir`/`static/` directory can be used to store, well, _static_ files to which you then can link in your articles.
You can put there whatever file you like either from the command-line or by using the system's `/ss` URL.

### Views

The `datadir`/`views/` directory holds the templates with which the final HTML pages are generated.
Provided that you feel at home working with _Go_ templates you might change them as you see fit.
I will, however, __not__ provide any support for you changing the default template structure.

A concise overview of the used templates and which variables they use you'll find in the file [template_vars.md](template_vars.md)

### Contents

For all the article you write – either on the commandline or with the web-interface – you can use [Markdown](https://en.wikipedia.org/wiki/Markdown) to enrich the plain text.
In fact, the system _expects_ the postings to be using `MarkDown` syntax if any markup at all.

## Libraries

The following external libraries were used building `Nele`:

* [ApacheLogger](https://github.com/mwat56/apachelogger)
* [BlackFriday](https://gopkg.in/russross/blackfriday.v2)
* [ChromeDP](https://github.com/chromedp/chromedp)
* [Crypto](https://golang.org/x/crypto)
* [CSSfs](https://github.com/mwat56/cssfs)
* [ErrorHandler](https://github.com/mwat56/errorhandler)
* [GzipHandler](https://github.com/NYTimes/gziphandler)
* [Hashtags](https://github.com/mwat56/hashtags)
* [INI](https://github.com/mwat56/ini)
* [Jffs](https://github.com/mwat56/jffs)
* [PassList](https://github.com/mwat56/passlist)
* [ScreenShot](https://github.com/mwat56/screenshot)
* [UploadHandler](https://github.com/mwat56/uploadhandler)
* [WhiteSpace](https://github.com/mwat56/whitespace)
* [wkHtmlToImage](https://wkhtmltopdf.org/)

## Licence

    Copyright © 2019, 2022 M.Watermann, 10247 Berlin, Germany
                    All rights reserved
                EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program.  If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.

----
