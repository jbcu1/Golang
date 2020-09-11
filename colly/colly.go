package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
	//colly.AllowedDomains("https://www.exploit-db.com/google-hacking-database"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visited", r.URL)
	})



	c.OnHTML("tr", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		e.ForEach("tbody tr td", func(_ int, el *colly.HTMLElement){
			
		})
	})
	/*
		c.OnResponse(func(r *colly.Response) {
			fmt.Println(string(r.Body))
		})
	*/
	c.Visit("https://www.exploit-db.com/google-hacking-database")
}
