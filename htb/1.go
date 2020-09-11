package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"strings"
)

var key string

func main(){

	c:=colly.NewCollector()
	//ch:=make(chan string,1)
	c.OnHTML("h3", func(e *colly.HTMLElement)  {
		key = e.Text
	})
	c.Visit("http://docker.hackthebox.eu:30520")
	enc:=md5.New()
	enc.Write([]byte(key))
	key2:=hex.EncodeToString(enc.Sum(nil))
	data:=strings.NewReader("hash="+key2)
	resp,_:=http.Post("http://docker.hackthebox.eu:30520", "application/x-www-form-urlencoded", data )
	//resp,_:=http.NewRequest("POST","http://docker.hackthebox.eu:30520", data)
	//resp.Header.Set("Content-Type:", "application/x-www-form-urlencoded")
	fmt.Println(resp.Request)
	defer resp.Body.Close()
	html,_:=ioutil.ReadAll(resp.Body)
	fmt.Println(string(html))
}
