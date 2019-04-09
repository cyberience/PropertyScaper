# PropertyScaper

Golang Scraper for Property

Requirements are as follows:

goquery for the web tools and sqlit for the database sqlite driver mayt have a problem if you use an older version of golang
```
github.com/PuerkitoBio/goquery
go get github.com/mattn/go-sqlite3
go get golang.org/x/text/currency
```


```
from rightmove_webscraper import rightmove_data
url = "https://www.rightmove.co.uk/property-for-sale/find.html?searchType=SALE&locationIdentifier=REGION%5E904"
rightmove_object = rightmove_data(url)
```

The following commands work

