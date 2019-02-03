package main

import (
    "fmt"
    "strings"
    "time"

    "github.com/gocolly/colly"
)

const (
    baseURL = "https://www.hwsw.hu/"
    articleLink = ".cikkdiv h3 a"
    articleTitle = "#cikkfejlec h1"
    articleDate = "#cikkfejlec .iro:not(a)"
    articleAuthor = "#cikkfejlec .iro a"
    articleContent = "#content #tartalom"
    firstYear = 2000
)

func main() {
    for year := firstYear; year <= time.Now().Year(); year++ {
        for month := 1; month <= 12; month++ {
            scrape(fmt.Sprintf("%s/archivum/%d/%d", baseURL, year, month))
        }
    }
}

func getContentBySelector(collector *colly.HTMLElement, selector string) string {
    return strings.TrimSpace(collector.DOM.Find(selector).Text())
}

func scrape(listingURL string) {
    listingCollector := colly.NewCollector()

    listingCollector.OnHTML(articleLink, func(articlePage *colly.HTMLElement) {
        articleCollector := colly.NewCollector()

        articleCollector.OnHTML("body", func(article *colly.HTMLElement) {
            title := getContentBySelector(article, articleTitle)
            author := getContentBySelector(article, articleAuthor)
            date := strings.TrimSpace(strings.Split(getContentBySelector(article, articleDate), ",")[1])
            content := getContentBySelector(article, articleContent)

            fmt.Println("Title:", title)
            fmt.Println("Author:", author)
            fmt.Println("Date:", date)
            fmt.Println("Content:", content)
            fmt.Println("=========")
        })

        articleCollector.OnError(func(r *colly.Response, err error) {
            fmt.Println("[articleCollector] Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
        })

        articleCollector.Visit(articlePage.Attr("href"))

        time.Sleep(2 * time.Second)
    })

    listingCollector.OnError(func(r *colly.Response, err error) {
        fmt.Println("[listingCollector] Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
    })

    listingCollector.Visit(listingURL)

    time.Sleep(2 * time.Second)
}
