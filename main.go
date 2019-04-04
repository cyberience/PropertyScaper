package main

import (
    "github.com/PuerkitoBio/goquery"
    "io"
    "log"
    "net/http"
    "os"
    "sync"
    "time"
)

var waitGroup sync.WaitGroup //to ensue the download is complete before next file
var urls = getUrls()

func main() {
    for _ , url := range urls {
        getPage(url)
        log.Print(url)
    }
}

func getPage(_url string) {
    defer func() {
        r := recover()
        if r != nil {
            log.Println("Request Disconnection ignored:", r)
        }
    }()

    client := &http.Client{Timeout:30 * time.Second}
    request, err := http.NewRequest("GET", _url,nil )
    if err != nil {
        log.Println("Make New request")
        log.Fatal(err) }
    request.Header.Set("User-Agent", "Stealing your Data")

    // Make HTTP GET request
    response, err := client.Do(request )
    if err != nil {
        log.Println("Make HTTP request")
    }
    defer response.Body.Close()
    document, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil { log.Println("Error loading HTTP response body. ") }
    getRow(document)



    return
}

func getRow(documment *goquery.Document) (getRow string){
    // This gets the UL Block
    documment.Find("ul.listing-results").Each(func(index int, element *goquery.Selection) {
        // Now get a Row, what about row count?
        log.Print(element.Html())
    })

    return
}

func parseRow(row string) {
    log.Print(row)
}


func getUrls() []string {
    return []string{"https://www.zoopla.co.uk/for-sale/houses/manchester/?q=Manchester&radius=0&results_sort=newest_listings&search_source=refine",
           "https://www.zoopla.co.uk/for-sale/houses/manchester/?identifier=manchester&property_type=houses&q=Manchester&search_source=refine&radius=0&pn=2"}
}


// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) (ret int) {
    defer func() {
        r := recover()
        if r != nil { log.Println("File error Ignored:", r) }
    }()
    defer waitGroup.Done()
    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        log.Print(err)
        ret = 0
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        log.Print(err)
        ret = 0
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        log.Print(err)
        ret = 0
    } else {
        ret = 1
    }

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