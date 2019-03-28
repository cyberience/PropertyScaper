#!/usr/bin/env python3

from rightmove_webscraper import rightmove_data
from pandas import DataFrame

def test_valid_urls(urls):
    for u in urls:
        print("Getting URL: {}".format(u))
        try:
            data = rightmove_data(u)
            if len(data.get_results) > 0:
                print("Rows:",len(data.get_results))
                print("\n".data.get_results)
                print("\n".data)
                if isinstance(data.summary(), DataFrame):
                    print("summary()\n".format(DataFrame))
            else:
                print("no results collected.")

        except Exception as e:
            print("URL: {}\n> FAILED with Exception:\n\t{}".format(u, e))
        print()

urls = ["https://www.rightmove.co.uk/property-for-sale/find.html?searchType=SALE&locationIdentifier=REGION%5E904"]
test_valid_urls(urls)


# from rightmove_webscraper import rightmove_data
# url = "https://www.rightmove.co.uk/property-for-sale/find.html?searchType=SALE&locationIdentifier=REGION%5E904&insId=1&radius=0.0&minPrice=&maxPrice=&minBedrooms=&maxBedrooms=&displayPropertyType=&maxDaysSinceAdded=&_includeSSTC=on&sortByPriceDescending=&primaryDisplayPropertyType=&secondaryDisplayPropertyType=&oldDisplayPropertyType=&oldPrimaryDisplayPropertyType=&newHome=&auction=false"
# rightmove_object = rightmove_data(url)