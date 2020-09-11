package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		fmt.Println(os.Args[1:])
		links, err := findLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		}
		arr := make([]string, 0, len(links))
		fmt.Printf("%#+v\n", links)
		for _, link := range links {
			fmt.Println(link)
			if strings.Contains(link, "http") {
				arr = append(arr, link)
			}

		}
		fmt.Println(arr)
	}
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	fmt.Println(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}
