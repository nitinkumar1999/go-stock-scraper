package main

import (
	"encoding/csv"

	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Stock struct {
	company, price, change string
}

func main() {
	ticker := []string{"GOOG", "CLOV", "GOOGL", "K", "CSCO", "SERV", "SIGA", "LUMN", "GNLN", "BVNRY"}
	stocks := []Stock{}

	c := colly.NewCollector() // it initializes new instance of colly

	c.OnRequest(func(r *colly.Request) { //to make requests on website
		fmt.Println("Visiting:", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{}

		stock.company = e.ChildText("h1")                                                   // company data
		fmt.Println("Company:", stock.company)                                              // printing the company data
		stock.price = e.ChildText("fin-streamer[data-field='regularMarketPrice']")          // access to price
		fmt.Println("Price:", stock.price)                                                  // printing the stock price
		stock.change = e.ChildText("fin-streamer[data-field='regularMarketChangePercent']") // price change
		fmt.Println("Change:", stock.change)                                                // printing price change

		stocks = append(stocks, stock)

	})

	for _, t := range ticker {
		fmt.Println("fetching data for:", t)
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")

	}

	c.Wait()

	fmt.Println(stocks)

	file, err := os.Create("stocks.csv")

	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file) // write data to the file

	headers := []string{
		"company",
		"price",
		"change",
	}
	writer.Write(headers)

	for _, stock := range stocks {
		record := []string{
			stock.company,
			stock.price,
			stock.change,
		}
		writer.Write(record)
	}
	writer.Flush()

}
