package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

func main() {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{"https://www.exploit-db.com/google-hacking-database"},
		ParseFunc: quotesParse,
		Exporters: []export.Exporter{&export.JSON{}},
	}).Start()

}

func quotesParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.card-body div#exploits-table_wrapper div.row div.col-sm-12 table#exploits-table tbody tr").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(r)
		g.Exports <- map[string]interface{}{
			"dorks-google": s.Find("tr.even").Text(),
		}
	})
	if href, ok := r.HTMLDoc.Find("li#exploits-table_next.paginate_button.page-item.next > a.page-link").Attr("href"); ok {
		g.Get(r.JoinURL(href), quotesParse)
	}
}
