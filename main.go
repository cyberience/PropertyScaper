package main

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "log"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"
)

var waitGroup sync.WaitGroup //to ensue the download is complete before next file

func main() {
    getPages("All")
}

// This may be modified with search criteria
func getUrls(key string) []string {
    return []string{"https://www.zoopla.co.uk/for-sale/houses/manchester/?q=Manchester&radius=0&results_sort=newest_listings&search_source=refine"}
}

func getPages( filter string) {
    getZoopla(filter)

    return
}

func getZoopla(filter string){
    defer func() {
        r := recover()
        if r != nil {
            log.Println("Request Disconnection ignored:", r)
        }
    }()
    pageNo , totalPages := 1 , 2
    var fullIdList []int64
    for {
        var url = fmt.Sprintf("https://www.zoopla.co.uk/for-sale/houses/manchester/?q=Manchester&radius=0&results_sort=newest_listings&search_source=refine&pn=%d",pageNo)
        client := &http.Client{Timeout:30 * time.Second}
        request, err := http.NewRequest("GET", url, nil )
        if err != nil{log.Println("Make New request", err)}
        request.Header.Set("User-Agent", "Stealing your Data")

        // Make HTTP GET request
        response, err := client.Do(request )
        if err != nil{ log.Println("Make HTTP request") }
        defer response.Body.Close()
        document, err := goquery.NewDocumentFromReader(response.Body)
        if err != nil{log.Println("Error loading HTTP response body. ")}
        if pageNo == 1 { // Get a page count on the first load
            totalPages = getTotalPages(document)
        }
        fullIdList = append(fullIdList, getIdList(document )... )
        log.Print("=============================================\n", fullIdList)
        if pageNo == totalPages {
            break
        } else {
            pageNo++
        }
    }
    //Should I save the ID list?
    for _, page := range fullIdList {
        extractPage(page)
    }
}

func extractPage(pageId int64) {
    var url = fmt.Sprintf("https://www.zoopla.co.uk/for-sale/details/%d",pageId)

}


func getIdList(documment *goquery.Document) (getIdList []int64){
    documment.Find("ul.listing-results").Each(func(index int, elementList *goquery.Selection) {
        elementList.Find("li").Each(func(index int, element *goquery.Selection) {
            itemId, valid := element.Attr("data-listing-id")
            if valid {
                idValue ,_ := strconv.ParseInt(itemId,0,64)
                getIdList = append(getIdList,idValue)
            }
        })
    })
    log.Println(getIdList)

    return
}

func getTotalPages(documment *goquery.Document) (getTotalPages int) {
    getTotalPages = 0
    documment.Find("div.paginate").Each(func(index int, elementList *goquery.Selection) {
        elementList.Find("a").Each(func(index int, element *goquery.Selection) {
            page := element.Text()
            currentPage, _ := strconv.Atoi(page)
            if currentPage > getTotalPages {
                getTotalPages = currentPage
            }
        })
    })

    return
}

// --------------------------------------------------
//    ___  _  _            ___       _      _
//   | __|(_)| | ___  ___ | __|__ __(_) ___| |_  ___
//   | _| | || |/ -_)|___|| _| \ \ /| |(_-<|  _|(_-<
//   |_|  |_||_|\___|     |___|/_\_\|_|/__/ \__|/__/
func exists(filePath string) (exists bool) {
    _,err := os.Stat(filePath)
    if err != nil {
        exists = false
    } else {
        exists = true
    }
    return
}

// --------------------------------------------------
// write the csv file for the content, change this
// to direct DB update later
func writeFile(msg string) (writeFile bool) {
    file, err := os.OpenFile("zinc.csv", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
    if err != nil { log.Println(err) }
    file.WriteString(msg+"\n")
    file.Sync()
    file.Close()
    return
}