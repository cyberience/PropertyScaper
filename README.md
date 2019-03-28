# PropertyScaper

Run python

execute the following

```
from rightmove_webscraper import rightmove_data
url = "https://www.rightmove.co.uk/property-for-sale/find.html?searchType=SALE&locationIdentifier=REGION%5E904"
rightmove_object = rightmove_data(url)
```

The following commands work


rightmove_object.average_price
rightmove_object.results_count
rightmove_object.get_results
rightmove_object.summary()
rightmove_object.summary(by = "postcode")
