# PropertyScaper

Golang Scraper for Property

Requirements are as follows:

goquery for the web tools and sqlite for the database sqlite driver mayt have a problem if you use an older version of golang.
To update Go use:
```
curl -o go.pkg https://dl.google.com/go/go1.11.1.darwin-amd64.pkg
shasum -a 256 go.pkg | grep 5cbd5505288bc2741091561dbfca05f4451824d557e275f373b8449112b84dff
sudo open go.pkg
```

```
github.com/PuerkitoBio/goquery
go get github.com/mattn/go-sqlite3
```


```
from rightmove_webscraper import rightmove_data
url = "https://www.rightmove.co.uk/property-for-sale/find.html?searchType=SALE&locationIdentifier=REGION%5E904"
rightmove_object = rightmove_data(url)
```

The following commands work

