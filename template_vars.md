# Variables used in templates

`./views/layout/`

* `01htmlpage.gohtml`:
  + `Lang` == the page's language
  + `Headline` == the page's `H1` headline
* `02htmlhead.gohtml`:
  + `CSS` == markup for `<style...>` head entries
  + `Robots` == directive for web-crawlers ("(no)index,(no)follow")
  + `Script` ==markup for `<script ...>` head entries
  + `Title` == the page's HTML `<title>` entry
* `05rightbar.gohtml`
  + `Taglist` == List of #hashtags/@mentions
* `06footer.gohtml`:
  + `Lang` == the page's language

`./views/`

* `article.gohtml`:
  + `Postings` == a list of postings, here consisting of a single entry with the elements
    - `ID` == the identifier of a single posting
    - `Posting` == the actual text of a single posting
* `index.gohtml`:
  + `Postings` == a list of postings, each consisting of entries with the elements
    - `ID` == the identifier of the respective posting
    - `Posting` == the actual text of the respective posting
* `imprint.gohtml`:
  + `Lang` == the page's language
* `privacy.gohtml`:
  + `Lang` == the page's language
* `searchressult.gohtml`:
  + `Lang` == the page's language
  + `Matches` == number of search results
  + `Postings` == a list of postings, each consisting of entries with the elements
    - `ID` == the identifier of the respective posting
    - `Posting` == the actual text of the respective posting
