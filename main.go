package main

import (
    "database/sql"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "log"
    "net/http"
    "os"
    "regexp"
    "strconv"
    "strings"
    "sync"
    "time"
    _ "github.com/mattn/go-sqlite3"
)

type property struct {
    id      int64
    title   string
    price   int64
    address string
    postcode string
    bedrooms int64
    bathrooms int64
    receptions int64
    description string
    agent string
    agentadd string
    agentNo string
    history string
    url string
    averageVal int64
    estimatedRent int64
}

var Property property
var waitGroup sync.WaitGroup //to ensue the download is complete before next file

func main() {
    makeDB()
//    getPages("All")
    extractZooPage(40630072)
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
        if r != nil { log.Println("Request Disconnection ignored:", r) }
    }()
    pageNo , totalPages := 1 , 0
    var fullIdList []int64
    for {
        var url = fmt.Sprintf("https://www.zoopla.co.uk/for-sale/houses/manchester/?q=Manchester&radius=0&results_sort=newest_listings&search_source=refine&pn=%d",pageNo)
        document, err := pullUrl(url)
        if err != nil{log.Println("Error loading HTTP response body. ")}
        if pageNo == 1 { // Get a page count on the first load
            totalPages = getTotalZooPages(document)
        }
        fullIdList = append(fullIdList, getIdZooList(document )... )
        if pageNo == totalPages {
            break
        } else {
            pageNo++
        }
    }
    log.Print("====================================================================\n","Pages",totalPages,"\n", fullIdList)
    //Should I save the ID list?
    for _, propPage := range fullIdList {
        //Check the ID has not already been found
        extractZooPage(propPage)
        // save the record here to DB
    }
}


func extractZooPage(pageId int64) {
    defer func() {
        r := recover()
        if r != nil { log.Println("Request Disconnection ignored:", r) }
    }()
    Property.id = pageId
    var url = fmt.Sprintf("https://www.zoopla.co.uk/for-sale/details/%d",pageId)
    Property.url = url
    document, err := pullUrl(url)
    checkErr(err)
    reg := regexp.MustCompile("[^0-9]+") // Removes Currency and other text

    document.Find("div.dp-grid-wrapper").Each(func(index int, section *goquery.Selection) {
        section.Find("h1.ui-property-summary__title").Each(func(index int, element *goquery.Selection) {
            Property.title = strings.TrimSpace(element.Text())
        })
        section.Find("p.ui-pricing__main-price").Each(func(index int, element *goquery.Selection) {
            Property.price, _ = strconv.ParseInt(reg.ReplaceAllString(element.Text(), ""), 0, 64)
        })
        section.Find("h2.ui-property-summary__address").Each(func(index int, element *goquery.Selection) {
            Property.address = element.Text()
            addArr := strings.Split(element.Text(), " ")
            Property.postcode = addArr[len(addArr)-1] // Last element has postcode
        })
        section.Find("h4.ui-agent__name").Each(func(index int, element *goquery.Selection) {
            Property.agent = element.Text()
        })
        section.Find("address.ui-agent__address").Each(func(index int, element *goquery.Selection) {
            Property.agentadd = element.Text()
        })
        section.Find("p.ui-agent__tel").Each(func(index int, element *goquery.Selection) {
            element.Find("a.ui-link").Each(func(index int, item *goquery.Selection) {
                Property.agentNo = reg.ReplaceAllString(item.Text(), "")
            })
        })
        section.Find("section#property-details-tab").Each(func(index int, element *goquery.Selection) {
            element.Find("li.dp-features-list__item").Each(func(index int, details *goquery.Selection) {
                details.Find("span").Each(func(index int, item *goquery.Selection) {
                    if strings.Contains(item.Text(), "bedrooms") {
                        Property.bedrooms, _ = strconv.ParseInt(reg.ReplaceAllString(item.Text(), ""), 0, 64)
                    }
                    if strings.Contains(item.Text(), "bathrooms") {
                        Property.bathrooms, _ = strconv.ParseInt(reg.ReplaceAllString(item.Text(), ""), 0, 64)
                    }
                    if strings.Contains(item.Text(), "reception") {
                        Property.receptions, _ = strconv.ParseInt(reg.ReplaceAllString(item.Text(), ""), 0, 64)
                    }
                })
            })
            element.Find("div.dp-description__text").Each(func(index int, details *goquery.Selection) {
                Property.description = strings.TrimSpace(details.Text())
            })
        })
        section.Find("section#market-stats-tab").Each(func(index int, element *goquery.Selection) {
            element.Find("ul.dp-market-stats__price-list").Each(func(index int, list *goquery.Selection) {
                element.Find("li.dp-market-stats__price-list-item").Each(func(index int, line *goquery.Selection) {
                    line.Find("span.dp-market-stats__price").Each(func(index int, item *goquery.Selection) {
                        Property.averageVal, _ = strconv.ParseInt(reg.ReplaceAllString(item.Text(), ""), 0, 64)
                    })
                })
            })
            element.Find("div.dp-market-stats--border-top").Each(func(index int, list *goquery.Selection) {
                list.Find("span.dp-market-stats__price").Each(func(index int, item *goquery.Selection) {
                    Property.estimatedRent, _ = strconv.ParseInt(reg.ReplaceAllString(item.Text(), ""), 0, 64)
                })
            })
        })

    })
    document.Find("div.ui-layout__halves").Each(func(index int, section *goquery.Selection) {
        section.Find("section.dp-price-history-block").Each(func(index int, element *goquery.Selection) {
            element.Find("div.dp-price-history__item").Each(func(index int, rows *goquery.Selection) {
                rows.Find("span.dp-price-history__item-date").Each(func(index int, row *goquery.Selection) {
                    Property.history = Property.history + strings.TrimSpace(row.Text()) + " - "
                })
                rows.Find("span.dp-price-history__item-price").Each(func(index int, row *goquery.Selection) {
                    Property.history = Property.history + strings.TrimSpace(row.Text()) + " - "
                })
                rows.Find("span.dp-price-history__item-detail").Each(func(index int, row *goquery.Selection) {
                    Property.history = Property.history + strings.TrimSpace(row.Text()) + " \n"
                })
            })
        })
    })
    saveData(Property)
    outputCsv(Property)

}


func outputCsv(row property) {
    //    log.Println(Property)
    outFile := strconv.FormatInt( row.id,10) + `,`
    outFile += `"`+ row.title + `",`
    outFile += strconv.FormatInt( row.price ,10) + `,`
    outFile += `"`+ row.address + `",`
    outFile += `"`+ row.postcode + `",`
    outFile += strconv.FormatInt(row.bedrooms ,10)+ `,`
    outFile += strconv.FormatInt(row.bathrooms ,10)+ `,`
    outFile += strconv.FormatInt(row.receptions ,10)+ `,`
    outFile += `"`+ row.description + `",`
    outFile += `"`+ row.agent + `",`
    outFile += `"`+ row.agentadd + `",`
    outFile += `"`+ row.agentNo + `",`
    outFile += `"`+ row.history + `",`
    outFile += `"`+ row.url + `",`
    outFile += `"`+ strconv.FormatInt( row.averageVal ,10) + `",`
    outFile += `"`+ strconv.FormatInt( row.estimatedRent ,10) + `"`
    log.Print(outFile)
    writeFile(outFile)
}

func getIdZooList(document *goquery.Document) (getIdList []int64){
    document.Find("ul.listing-results").Each(func(index int, elementList *goquery.Selection) {
        elementList.Find("li").Each(func(index int, element *goquery.Selection) {
            itemId, valid := element.Attr("data-listing-id")
            if valid {
                idValue ,_ := strconv.ParseInt(itemId,0,64)
                getIdList = append(getIdList,idValue)
            }
        })
    })
    return // Returns a page of ID's
}

func getTotalZooPages(documment *goquery.Document) (getTotalPages int) {
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
    return // Returns the Total nu,mber of pages
}

func makeDB() {
    if !exists("./properties.db") {
        db, err := sql.Open("sqlite3", "./properties.db")
        checkErr(err)
        stmt, err := db.Prepare("CREATE TABLE `properties` (" +
            " `uid` INTEGER PRIMARY KEY AUTOINCREMENT," +
            "`id`           INTEGER, " +
            "`title`        VARCHAR(64) NULL, " +
            "`price`        INTEGER, " +
            "`address`      VARCHAR(64) NULL, " +
            "`postcode`     VARCHAR(64) NULL, " +
            "`bedrooms`     INTEGER, " +
            "`bathrooms`    INTEGER, " +
            "`receptions`   INTEGER, " +
            "`description`  VARCHAR(64) NULL, " +
            "`agent`        VARCHAR(64) NULL, " +
            "`agentadd`     VARCHAR(64) NULL, " +
            "`agentNo`      VARCHAR(64) NULL, " +
            "`history`      VARCHAR(64) NULL, " +
            "`url`          VARCHAR(64) NULL, " +
            "`averageVal`   INTEGER, " +
            "`estimatedRent` INTEGER," +
            "`created`      DATE NULL ) " )
        checkErr(err)
        res, err := stmt.Exec()
        checkErr(err)
        id, err := res.LastInsertId()
        checkErr(err)
        log.Println(id)
        db.Close()
    }
}

func saveData(p property) (saveData bool){
    defer func() {
        r := recover()
        if r != nil { log.Println("Request Disconnection ignored:", r) }
    }()
    db, err := sql.Open("sqlite3", "./properties.db")
    checkErr(err)
    stmt , err := db.Prepare("INSERT INTO properties(id, title, price, address, postcode, bedrooms, bathrooms, receptions, description, agent, agentadd, agentNo, history, url, averageVal, estimatedRent ) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,DateTime('now'))")
    checkErr(err)
    res, err := stmt.Exec(p.id,p.title,p.price,p.address,p.postcode,p.bedrooms,p.bathrooms,p.receptions,p.description,p.agent,p.agentadd,p.agentNo,p.history,p.url,p.averageVal)
    checkErr(err)
    affect, err := res.RowsAffected()
    checkErr(err)
    log.Print(affect)
    db.Close()
    return true
}

func pullUrl(_url string) (*goquery.Document, error ) {
    client := &http.Client{Timeout:30 * time.Second}
    request, err := http.NewRequest("GET", _url, nil )
    if err != nil{log.Println("Make New request", err)}
    request.Header.Set("User-Agent", "Stealing your Data")

    // Make HTTP GET request
    response, err := client.Do(request )
    if err != nil{ log.Println("Make HTTP request") }
    defer response.Body.Close()

    return goquery.NewDocumentFromReader(response.Body)
}


func checkErr(err error) {
    if err != nil {
        log.Print(err)
    }
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
    file, err := os.OpenFile("./properties.csv", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
    if err != nil { log.Println(err) }
    file.WriteString(msg+"\n")
    file.Sync()
    file.Close()
    return
}