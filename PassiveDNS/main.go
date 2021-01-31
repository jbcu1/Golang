package main

import (
	"PTwork/DBElastic"
	"PTwork/kafka_client"
	"PTwork/models"
	"PTwork/parser"
	"PTwork/recon"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)


var mainPagesParsResultDir = "main_page_pars_result"
var linksParsResultDir = "links_pars_result"
var startPageParsResultDir = "start_page_pars_result"

//Get current path
func currentPath() string{
	curPath,_:=os.Getwd()
	return curPath
}

//Create file
func CreateFile(fileName string) error{
	f,err:=os.Create(fileName)
	if err!=nil{
		return err
	}
	defer f.Close()
	return nil
}


//Copy data from file to another file
func CopyData(fileToCopy string,fileFromCopy string) error{

	fileOrigin,_:=os.Open(fileFromCopy)
	defer fileOrigin.Close()

	newFile,err:=os.Create(fileToCopy)
	if err!=nil{
		return err
	}
	defer newFile.Close()
	bytesWritten,err:=io.Copy(newFile,fileOrigin)
	if err!=nil{
		return err
	}
	fmt.Printf("Copied %d bytes.\n", bytesWritten)
	err = newFile.Sync()
	if err!=nil{
		return err
	}
	return nil
}

//Delete file
func RemoveFile(pathToFile string)error{
	err:=os.Remove(pathToFile)
	if err!=nil{
		return err
	}
	return nil
}



//Parse json file from start page parser
func OpenAndRead(fileName string, data models.MainPages)[]string{
	file,err:=os.Open(fileName)
	if err!=nil{
		_ = fmt.Errorf("Something went wrong %s\n", err)
	}
	defer file.Close()
	byteData,_:=ioutil.ReadAll(file)
	byteData=byteData[:len(byteData)-2]
	byteData = append(byteData, 93)
	_ = json.Unmarshal(byteData, &data)
	urlsSlice:=make([]string,0)
	for _,i:=range data{
		for _,j:=range i.Links{
			urlsSlice = append(urlsSlice, j)
		}
	}
	return urlsSlice
}

type Links [][]string

//Parse links from scrapping result links
func OpenAndReadLinks(filename string,links Links)[]string{
	file,err:=os.Open(filename)
	if err!=nil{
		_ = fmt.Errorf("sometring went wrong %s", err)
	}
	defer file.Close()
	byteData,_:=ioutil.ReadAll(file)
	byteData=byteData[:len(byteData)-2]
	byteData = append(byteData, 93)
	_ = json.Unmarshal(byteData, &links)
	urlSlice:=make([]string,0)
	for _,i:=range links{
		for _,j:=range i{
			urlSlice = append(urlSlice, j)
		}
	}
	return urlSlice
}

//Get slice directory names
func ReadDir(dirName string)[]string{
	files,err:=ioutil.ReadDir(dirName)
	if err!=nil{
		_ = fmt.Errorf("something went wrong %s", err)
	}
	filesArray:=make([]string,0)
	for _,f:=range files{
		filesArray = append(filesArray, f.Name())
	}
	return filesArray
}


//Parsing main page and links from main page
func ParsingMainPage(){
	curPath:=currentPath()
	_ = os.MkdirAll(curPath+"/pars_result/"+mainPagesParsResultDir, 0777)
	_ = os.MkdirAll(curPath+"/pars_result/"+linksParsResultDir, 0777)
	resultMainPagePars:=curPath+"/pars_result/MainPageInformation.json"
	mainPageFileTemp,_:=os.Create(resultMainPagePars)
	defer mainPageFileTemp.Close()
	parser.MainPageParser(resultMainPagePars)
	mainPageFileSave:=curPath+"/pars_result/"+mainPagesParsResultDir+"/"+time.Now().Format("2006-1-2")+".json"
	_ = CopyData(mainPageFileSave, resultMainPagePars)
	defer RemoveFile(resultMainPagePars)

	var mainPages models.MainPages
	_ = os.MkdirAll(curPath+"/pars_result/"+linksParsResultDir+"/"+time.Now().Format("2006-1-2"), 0777)
	sN:="https://domain-status.com/archives/"+time.Now().AddDate(0,0,-1).Format("2006-1-2")+"/"

	for _,i:=range OpenAndRead(resultMainPagePars,mainPages){
		nameOfFile:=strings.Replace(i,sN,time.Now().Format("2006-1-2")+"_",1)
		nameOfFile=strings.Replace(nameOfFile,"/1",".json",1)
		nameOfFile=strings.ReplaceAll(nameOfFile,"/","_")
		_ = CreateFile(curPath + "/pars_result/" + linksParsResultDir + "/" + time.Now().Format("2006-1-2") + "/" + nameOfFile)
		startUrlArray:=strings.Fields(i)
		parser.PageParser(curPath+"/pars_result/"+linksParsResultDir+"/"+time.Now().Format("2006-1-2")+"/"+nameOfFile,startUrlArray,100)
	}
}

//Create range date
func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

//Create links slice for scrape domain amount
func CreateLinksSlice()[]string{
	urlsSlice:=make([]string,0)
	url:="https://domain-status.com/archives/"
	start:=time.Date(2017,7,22,0,0,0,0, time.UTC)
	end:=time.Now()
	for rd:=rangeDate(start,end); ; {
		date:=rd()
		if date.IsZero(){
			break
		}
		url := url + date.Format("2006-1-2")+"/"
		urlsSlice = append(urlsSlice, url)
	}
	return urlsSlice
}

//Write parser for scrape all start pages
func GetAllStartPages(){
	url:="https://domain-status.com/archives/"
	linkSlice:=CreateLinksSlice()
	curPath:=currentPath()
	workDir:=curPath+"/pars_result/"+startPageParsResultDir+"/"
	for _,i:=range linkSlice{
		dir:=strings.ReplaceAll(i,url,"")
		filename:=strings.ReplaceAll(dir,"/",".json")
		filename=workDir+dir+filename
		_ = os.MkdirAll(workDir+dir, 0777)
		parser.StartPageParser(filename,strings.Fields(i))

	}

}

//Get information for tranco list and indexed
func InformationGatheringForTranco() {
	testPost:=models.TrancoDomainData{
		Host:             "",
		Tld:              "",
		Status:           "",
		Identifier:       "",
		Subdomains:       nil,
		FirstSeen:        "",
		LastSeen:         "",
		Resolvers:        nil,
		WhoISInformation: nil,
		Sources:          nil,
		Tags:             nil,
	}
	delID,err:=DBElastic.SendPost(testPost)
	if err!=nil{
		_ = fmt.Errorf("Cannot insert test doc %s\n", err)
	}
	domains:=recon.OpenAndReadTrancoCSV()
	reader,readerA:=recon.OpenGeoLite()

	var wg sync.WaitGroup

	dataCh := make(chan models.TrancoDomainData)
	hosts:=make(chan string)

	for i:=0;i<=concurrency;i++{
		go recon.GetInformationAboutHostTranco(hosts,dataCh,&wg,reader,readerA)
		wg.Add(1)
	}

	numJobs:=len(domains)
	fmt.Printf("Start gathering information from %v hostst. Please wait.\n",len(domains))

	go func(){
		for i:=0;i<numJobs;i++{
			hosts <- domains[i]
		}
		close(hosts)
	}()

	dataSlice:=make([]models.TrancoDomainData,0,numJobs)
	finishProcess:=float64(cap(dataSlice))

	for i:=0;i<numJobs;i++{
		dataSlice = append(dataSlice, <-dataCh)
		startProcess:=float64(len(dataSlice))
		if int(startProcess)%100==0{
			percentProcess:=(startProcess/finishProcess)*100
			fmt.Printf("The array is %.2f perc full.\n", percentProcess)
		}
	}
	wg.Wait()

	fmt.Println("Array is full")
	fmt.Println("Starting bulk index")
	_ = DBElastic.BulkSendPostTranco(dataSlice)
	_ = DBElastic.DeletePost(delID)
	fmt.Println("Collection information completed. All post are indexed.")
}


//Indexed one day scrapping
func InformationGatheringForOnePage() {

	start:=time.Now()
	testPost:=models.DomainData{
		Host:             "",
		Tld:              "",
		Status:           "",
		Subdomains:       nil,
		FirstSeen:        "",
		LastSeen:         "",
		Resolvers:        nil,
		WhoISInformation: nil,
		Source:           "",
		Tag:              nil,
	}

	delId,err:=DBElastic.SendPost(testPost)
	if err!=nil{
		_ = fmt.Errorf("Cannot insert test doc %s\n", err)
	}


	domains,_:=GetHostsSlice()
	reader,readerA:=recon.OpenGeoLite()

	var wg sync.WaitGroup

	dataCh := make(chan models.DomainData)
	hosts:=make(chan string)

	for i:=0;i<=concurrency;i++{
		go recon.GetInformationAboutHost(hosts,dataCh,&wg,reader,readerA)
		wg.Add(1)
	}

	numJobs:=len(domains)
	fmt.Printf("Start gathering information from %v hostst. Please wait.\n",len(domains))

	go func(){
		for i:=0;i<numJobs;i++{
			hosts <- domains[i]
		}
		close(hosts)
	}()

	dataSlice:=make([]models.DomainData,0,numJobs)
	finishProcess:=float64(cap(dataSlice))

	for i:=0;i<numJobs;i++{
		dataSlice = append(dataSlice, <-dataCh)
		startProcess:=float64(len(dataSlice))
		if int(startProcess)%100==0{
			percentProcess:=(startProcess/finishProcess)*100
			fmt.Printf("The array is %.2f perc full.\n", percentProcess)
		}
	}
	wg.Wait()

	fmt.Println("Array is full")
	fmt.Println("Starting bulk index")
	_ = DBElastic.BulkSendPost(dataSlice)
	_ = DBElastic.DeletePost(delId)
	fmt.Println("Collection information completed. All post are indexed.")
	fmt.Printf("Time taken: %s\n", time.Since(start))
}

//Gathering information from all scrapping dir
func InformationGatheringForAllExistDir() {
	domains:=ReadDirWithLinks()
	reader,readerA:=recon.OpenGeoLite()

	var wg sync.WaitGroup

	dataCh := make(chan models.DomainData)
	hosts:=make(chan string)
	for i:=0;i<=concurrency;i++{
		go recon.GetInformationAboutHost(hosts,dataCh,&wg, reader, readerA)
		wg.Add(1)
	}

	fmt.Println(len(domains))
	go func(){
		for _,i:=range domains{
			hosts <- i
		}
		close(hosts)
	}()


	for i:=0;i<=len(domains);i++{
		_, _ = DBElastic.SendPost(<-dataCh)
	}
	wg.Wait()
	fmt.Println()

}

//Get all dir with links for indexed next step
func ReadDirWithLinks()[]string{
	var link Links
	dirSlice:=make([]string,0)
	curPath:=currentPath()
	dirPath:=curPath+"/pars_result"
	for _,i:=range ReadDir(dirPath){
		fmt.Println(i)
	}
	var dir string
	fmt.Println("Enter directory witch contain scrapping hosts:")
	_, _ = fmt.Fscan(os.Stdin, &dir)

	for _,i:=range ReadDir(dirPath+"/"+dir+"/"){
		dirSlice = append(dirSlice, dirPath+"/"+dir+"/"+i+"/")
	}

	fileSlice:=make([]string,0)
	for _,i:=range dirSlice{
		fmt.Println(i)
		for _,j:=range ReadDir(i){
			fmt.Println(j)
			fileSlice = append(fileSlice, i+"/"+j)
		}
	}

	hosts:=make([]string,0)
	for _,i:=range fileSlice{

		for _,j:=range OpenAndReadLinks(i,link){

			hosts = append(hosts, j)

		}
	}

	return hosts
}


//One directory with user choose
func GetHostsSlice()([]string,string){

	source:=""
	domains:=make([]string,0)
	curPath:=currentPath()
	resultPath:=curPath+"/pars_result"

	for _,i:=range ReadDir(resultPath){

		fmt.Println(i)

	}

	var path string
	fmt.Println("Enter directory:")
	_, _ = fmt.Fscan(os.Stdin, &path)
	if path=="links_pars_result"{
		source="domain-status.com"
	}

	var dir string
	for _,i:=range ReadDir(resultPath+"/"+path){
		fmt.Println(i)
	}

	fmt.Println("Enter directory from which will be created url slice:")
	_, _ = fmt.Fscan(os.Stdin, &dir)

	var links Links
	for _, i := range ReadDir(curPath + "/pars_result/" + linksParsResultDir + "/" + dir) {

		for _, j := range OpenAndReadLinks(curPath+"/pars_result/"+linksParsResultDir+"/"+dir+"/"+i, links) {

			domains = append(domains, j)

		}

	}

	return domains,source
}


//Get hosts slice from today's scrapping
func GetHostsSliceTodayHosts()([]string,string){

	source:="domain-status.com"
	domains:=make([]string,0)
	curPath:=currentPath()
	resultPath:=curPath+"/pars_result/links_pars_result/"+time.Now().Format("2006-1-2")

	var links Links
	for _, i := range ReadDir(resultPath) {

		for _, j := range OpenAndReadLinks(resultPath+"/"+i, links) {

			domains = append(domains, j)

		}

	}

	return domains,source
}


//kafka producer with goroutine
func KafkaProducer(){

	c:=context.Background()
	var wg sync.WaitGroup
	var kM models.KafkaMessage

	a,source:=GetHostsSliceTodayHosts()

	hosts:=make(chan string)

	kM.Tag=[]string{"test","test"}
	kM.Source=source

	for i:=0;i<=kafkaConcurrency;i++{

		go kafka_client.Produce(hosts,c,&wg,kM)
		wg.Add(1)

	}

	numJobs:=len(a)
	go func(){

		for i:=0;i<numJobs;i++{

			hosts <- a[i]

		}

		close(hosts)

	}()

	wg.Wait()
}

func KafkaProducerWithChoice(){

	c:=context.Background()
	var wg sync.WaitGroup
	var kM models.KafkaMessage

	a,source:=GetHostsSlice()

	hosts:=make(chan string)

	kM.Tag=[]string{"test","test"}
	kM.Source=source

	for i:=0;i<=kafkaConcurrency;i++{

		go kafka_client.Produce(hosts,c,&wg,kM)
		wg.Add(1)

	}

	numJobs:=len(a)
	go func(){

		for i:=0;i<numJobs;i++{

			hosts <- a[i]

		}

		close(hosts)

	}()

	wg.Wait()
}


func KafkaConsumer() {


	start:=time.Now()
	testPost:=models.DomainData{
		Host:             "",
		Tld:              "",
		Status:           "",
		Subdomains:       nil,
		FirstSeen:        "",
		LastSeen:         "",
		Resolvers:        nil,
		WhoISInformation: nil,
		Source:           "",
		Tag:              nil,
	}

	delId,err:=DBElastic.SendPost(testPost)
	if err!=nil{
		_ = fmt.Errorf("Cannot insert test doc %s\n", err)
	}

	reader,readerA:=recon.OpenGeoLite()

	var wg sync.WaitGroup

	msgKey := make(chan string)
	msgVal := make(chan string)
	dataCh := make(chan models.DomainData)
	hosts:=make(chan string)

	ctx := context.Background()

	go kafka_client.Consumer(ctx,msgKey,msgVal)


	for i:=0;i<=concurrency;i++{

		go recon.GetInformationAboutHost(hosts,dataCh,&wg,reader,readerA)
		wg.Add(1)

	}

	go func(){
		for i:=range msgVal{
			fmt.Println(i)
		}
	}()

	go func(){

		for mK:=range msgKey{
			fmt.Println(mK)
			hosts <- mK

		}
		close(hosts)
	}()


	chunkSize:=15000
	dataSlice:=make([]models.DomainData,0,chunkSize)

	for i:=range dataCh{
		dataSlice = append(dataSlice, i)


		if len(dataSlice)>=chunkSize{
			_ = DBElastic.BulkSendPost(dataSlice)
			dataSlice = dataSlice[:0]

		}


	}
	_ = DBElastic.BulkSendPost(dataSlice)

	wg.Wait()

	_ = DBElastic.DeletePost(delId)

	fmt.Println("Collection information completed. All post are indexed.")
	fmt.Printf("Time taken: %s\n", time.Since(start))

}


var concurrency int
var kafkaConcurrency int
var parse bool
var informationGathering bool
var parsingAndGatheringMainPage bool
var indexAllScrapeHosts bool
var parsAllStartPages bool
var trancoInformGathering bool
var kafkaProduce bool
var kafkaProduceWithChoice bool
var kafkaConsumer bool

func init(){
	flag.IntVar(&concurrency,"c", 100, "Number of concurrency for gathering information about hosts\n")
	flag.IntVar(&kafkaConcurrency, "kc", 100, "Number of concurrency fo working with kafka\n")
	flag.BoolVar(&parse, "p", false, "Start parsing main page\n")
	flag.BoolVar(&informationGathering,"g",false,"Start information gathering about hosts scraped at those day\n")
	flag.BoolVar(&parsingAndGatheringMainPage,"pG",false,"Start parsing main page and index all\nscrapping url\n")
	flag.BoolVar(&indexAllScrapeHosts,"gA",false,"Start index scrapping hosts of all time\n")
	flag.BoolVar(&parsAllStartPages,"pS",false,"Scrapping all start pages from domain-status archive\n")
	flag.BoolVar(&trancoInformGathering,"gT", false, "Start information gathering from tranco list\n")
	flag.BoolVar(&kafkaProduce, "k",false,"Starting upload hosts slice to kafka queue\n")
	flag.BoolVar(&kafkaProduceWithChoice,"kPC", false,"Custom choice of folder to upload to kafka queue\n")
	flag.BoolVar(&kafkaConsumer,"kC", false, "Getting hosts from the kafka queue and further enriching host information\n")
}

func main() {

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, strings.Join([]string{
			"\"Get information about hosts\" this is pDNS system to get useful information about host.",
			"",
			"Usage: resolve [option ...]",
			"",
		}, "\n"))
		flag.PrintDefaults()
	}


	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}
	if parse{
		ParsingMainPage()
	}
	if informationGathering{
		InformationGatheringForOnePage()
	}
	if parsingAndGatheringMainPage{
		ParsingMainPage()
		InformationGatheringForOnePage()
	}
	if indexAllScrapeHosts{
		InformationGatheringForAllExistDir()
	}
	if parsAllStartPages{
		GetAllStartPages()
	}
	if trancoInformGathering{
		InformationGatheringForTranco()
	}
	if kafkaProduce{
		go KafkaProducer()
		KafkaConsumer()
	}
	if kafkaProduceWithChoice{
		go KafkaProducerWithChoice()
	}
	if kafkaConsumer{
		KafkaConsumer()
	}

}



