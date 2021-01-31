package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"strings"
	"time"
)


//Imported function to pars pages with registered domains from domain-status.com
func PageParser(saveFilePath string, urlsArr []string,retryTimes int){

	geziyor.NewGeziyor(&geziyor.Options{

		StartURLs: urlsArr,
		ParseFunc: linksParse,
		Exporters: []export.Exporter{&export.JSON{FileName: saveFilePath}},
		//RequestDelay: time.Second*2,
		RobotsTxtDisabled: true,
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
		RetryTimes: retryTimes,

	}).Start()

}

//Imported function to scrapped archived main pages from domain-status.com
func StartPageParser(fileName string, startPages []string){

	geziyor.NewGeziyor(&geziyor.Options{

		StartURLs: startPages,
		ParseFunc: startPagePars,
		Exporters: []export.Exporter{&export.JSON{FileName: fileName,EscapeHTML: false}},
		RobotsTxtDisabled: true,
		RequestDelay: time.Second,
		RequestDelayRandomize: true,
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
		RetryHTTPCodes: []int{502},
		RetryTimes: 1000,

	}).Start()

}

//Function to scrapping archive start pages from domain-status.com
func startPagePars(g *geziyor.Geziyor, r *client.Response){

	r.HTMLDoc.Find("div.column > ul").Each(func(i int, s *goquery.Selection) {

		domainZone:=s.Find("h3").Text()

		if domainZone!=""{
			domainAmount:=strings.Split(s.Find("li").Text(),"\n ")
			domainRegisterAmount:=make([]string,0)
			links:=make([]string,0)

			for i:=range domainAmount{

				domainAmount[i]=strings.TrimSpace(domainAmount[i])

				if len(domainAmount[i])!=0{
					domainRegisterAmount=append(domainRegisterAmount,domainAmount[i])
				}

			}

			s.Find("a").Each(func(ii int, ss *goquery.Selection) {

				link, ok := ss.Attr("href")

				if ok {
					links = append(links, link)
				}

			})

			g.Exports <- map[string]interface{}{
				"domain_zone": domainZone,
				"domain_amount": domainRegisterAmount,
				"links": links,
			}

		}

	})
}


//Imported function to collect information from main page domain-status.com
func MainPageParser(fileName string){

	geziyor.NewGeziyor(&geziyor.Options{

		StartURLs: []string{"https://domain-status.com/"},
		ConcurrentRequests: 2,
		ParseFunc: mainPagePars,
		Exporters: []export.Exporter{&export.JSON{FileName: fileName,EscapeHTML: false}},
		RobotsTxtDisabled: true,
		RequestDelay: time.Second,
		RequestDelayRandomize: true,
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36",
		RetryHTTPCodes: []int{502},
		RetryTimes: 100,

	}).Start()

}


//Parse main page from domain-status.com
func mainPagePars(g *geziyor.Geziyor, r *client.Response){

	r.HTMLDoc.Find("div.column.small-12.medium-4.large-3.container-hover").Each(func(i int, s *goquery.Selection) {

		domainZone:=s.Find("h3").Text()

		if domainZone!=""{
			domainAmount:=strings.Split(s.Find("li").Text(),"\n ")
			domainRegisterAmount:=make([]string,0)
			links:=make([]string,0)

			for i:=range domainAmount{

				domainAmount[i]=strings.TrimSpace(domainAmount[i])

				if len(domainAmount[i])!=0{
					domainRegisterAmount=append(domainRegisterAmount,domainAmount[i])
				}
			}

			s.Find("a").Each(func(ii int, ss *goquery.Selection) {

				link, ok := ss.Attr("href")

				if ok {
					links = append(links, link)
				}

			})

			g.Exports <- map[string]interface{}{
				"domain_zone": domainZone,
				"domain_amount": domainRegisterAmount,
				"links": links,
			}

		}

	})
}


//Function to collect domains from domain-status.com
func linksParse(g *geziyor.Geziyor, r *client.Response){

	r.HTMLDoc.Find("div.row.expanded.chained-lists").Each(func (i int, s *goquery.Selection){

		links:=make([]string,0)
		linksText:=s.Find("li").Text()
		linksText=strings.ReplaceAll(linksText,""," ")
		s.Find("a").Each(func (ii int, ss *goquery.Selection){

			link,ok:=ss.Attr("href")
			link=strings.ReplaceAll(link,"https://domain-status.com/www/","")
			if ok{
				links=append(links,link)
			}

		})

		if href, ok:=r.HTMLDoc.Find("li.pagination-next > a").Attr("href"); ok{
			g.Get(r.JoinURL(href), linksParse)
		}

		g.Exports <- links
	})
}

