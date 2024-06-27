package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type extJob struct {
	title, jobDate, corpName string
}

var baseURL string = "https://www.saramin.co.kr"
var searchWords string = ""

func init() {
	fmt.Println("init call...")
}

func Scrape(words string) {
	start := time.Now()

	// global var
	searchWords = words

	var pageJobs []extJob
	ch := make(chan []extJob)
	totalPage := getPages()
	fmt.Println(`totalPage : `, totalPage)

	
	for i:= 1; i <= totalPage; i++ {
		go getPage(i, ch)
	}
	for i:= 1; i <= totalPage; i++ {
		jobs := <-ch
		pageJobs = append(pageJobs, jobs...)
	}
	// fmt.Println(pageJobs)
	writeAllJobs(pageJobs)
	end := time.Now()
	fmt.Println("elapsed time : ", end.Sub(start))
}

func writeAllJobs(jobs []extJob) {
	fmt.Println("total arr : ", len(jobs))
	file, err := os.Create("jobs.csv")
	defer func() {if err == nil {file.Close()}}()
	checkError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	header := []string{"Title", "Date", "Corp"}
	wErr := w.Write(header)
	checkError(wErr)

	for _, job := range jobs {
		record := []string{job.title, job.jobDate, job.corpName}
		jErr := w.Write(record)
		checkError((jErr))
	}
}

func getPage(page int, mainCh chan []extJob) {
	var pageJobs []extJob
	pageURL := fmt.Sprintf("%s/zf_user/search/recruit?=&searchword=%s&recruitPageCount=40&recruitPage=%d", baseURL, searchWords, page)
	res, err := http.Get(pageURL)
	checkError(err)
	checkRespCode(res)
	
	fmt.Println("Requesting ", pageURL)
	
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
  checkError(err)
	
	ch := make(chan extJob)
	// Find the review items
	cards := doc.Find(".item_recruit")
	cards.Each(func(i int, card *goquery.Selection) {
		// For each item found, get the title
		go extractCard(card, ch)
	})

	for i:=0; i<cards.Length(); i++ {
		pageJobs = append(pageJobs, <-ch)
	}
	mainCh<-pageJobs
}

func extractCard(card *goquery.Selection, c chan extJob) {
	title, _ := card.Find(".area_job .job_tit a").Attr("title")
	jobDate := card.Find(".area_job .job_date span").Text()
	corpName := card.Find(".area_corp .corp_name a").Text()

	// fmt.Println(title, ":", jobDate, ":", corpName)
	c<-extJob{title: title, jobDate: jobDate, corpName: corpName}
}

func getPages() int {
	pageURL := fmt.Sprintf("%s/zf_user/search/recruit?=&searchword=%s&recruitPageCount=40", baseURL, searchWords)
	res, err := http.Get(pageURL)

	checkError(err)
	checkRespCode(res)

	defer res.Body.Close()
	// Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  checkError(err)

	page := 0
	// Find the review items
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		page = s.Find("a").Length()
	})

	return page
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkRespCode(res *http.Response) {
	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
}