# Variables used in templates

`./views/layout/`

* `01htmlpage.gohtml`: defines the overall structure of an HTML page.
  + `Lang` == the page's language
  + `Headline` == the page's `H1` headline

* `02htmlhead.gohtml`: includes the following HTML/HEAD entries:
  + `CSS` == markup for `<style...>` head entries
  + `Robots` == directive for web-crawlers ("(no)index,(no)follow")
  + `Title` == the page's HTML/HEAD `<title>` entry

* `03header.gohtml`: includes the BODY/HEADER element.
  + `Blogname` == the "name" of the blog
  + `Lang` == the page's language

* `05rightbar.gohtml`: fills the right side of the page.
  + `Taglist` == List of #hashtags/@mentions

* `06footer.gohtml`: includes the BODY/FOOTER element.
  + `Lang` == the page's language

`./views/`

* `ap.gohtml`: called for the URL `"/a"` to add a new posting.
  + `Lang` == the page's language

* `article.gohtml`: called for the URL `"/p/…"` to show a single posting.
  + `Posting` == a single posting with the elements
    - `Date` == the date of the single posting
    - `ID` == the identifier of the single posting
    - `Post` == the actual text of the single posting

* `dc.gohtml`: called for the URL `"/d/…"` to change an article's date/time.
  + `HMS` == the posting's hour-minute-second
  + `ID` == the ID of the current article
  + `Lang` == the page's language
  + `Manuscript` == the posting's text
  + `NOW` == the current date's year-month-date
  + `YMD` == the posting's year-month-day

* `ed.gohtml`: called for the URL `"/e/…"` to edit an article's text.
  + `ID` == the ID of the current article
  + `Lang` == the page's language
  + `Manuscript` == the posting's text

* `error.gohtml`: internally called to send error messages to the remote user.
  + `Error` == the respective error message

* `faq.gohtml`: called for the URL `"/faq"` to to show some FAQs
  + `Lang` == the page's language

* `imprint.gohtml`: called fo the URLs `"/imprint"` and `"/impressum"` to show the site's imprint.
  + `Lang` == the page's language
  + Please be aware that the actual contents of this file is subject to your own country's laws and legislation.

* `index.gohtml`: called for the URLs `"/"`, `"/index"` and `"/n/…"` to show a list of postings
  + `Postings` == a list of postings, each consisting of entries with the elements.
    - `Date` == the date of the respective posting
    - `ID` == the identifier of the respective posting
    - `Post` == the actual text of the respective posting

* `licence.gohtml`: called for URLs `"/licence"`, `"/license"`, and `"/lizenz"` to show the site content's licence.
  + `Lang` == the page's language
  + Please be aware that the actual contents of this file is subject to your own personal consideration.

* `privacy.gohtml`: called for the URL `"/privacy"` and `"/datenschutz"` to show the site's privacy statement.
  + `Lang` == the page's language
  + Please be aware that the actual contents of this file is subject to your own country's laws and legislation.

* `rp.gohtml`: called for the URL `"/r/…"` to remove an article's altogether.
  + `ID` == the ID of the current article
  + `Lang` == the page's language
  + `Manuscript` == the posting's text

* `searchressult.gohtml`: called for the URL `"/s/…"` to show the results of a search.
  + `Lang` == the page's language
  + `Matches` == number of search results
  + `Postings` == a list of postings, each consisting of entries with the elements
    - `Date` == the date of the respective posting
    - `ID` == the identifier of the respective posting
    - `Post` == the actual text of the respective posting

* `si.gohtml`: called for the URL `"/si"` to upload an image.
  + `Lang` == the page's language

* `ss.gohtml`: called for the URL `"/ss"` to upload a static file.
  + `Lang` == the page's language

