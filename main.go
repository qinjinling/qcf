package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
)

// newDataSet 查询结果外层包装
type newDataSet struct {
	Tables []*table `xml:"Table"`
}

// table 查询结果表格
type table struct {
	MasterPartID string `xml:"MasterPartID"`
	PartNumber   string `xml:"PartNumber"`
	IsStock      string `xml:"IsStock"`
}

// PartDetail 零件明细
type PartDetail struct {
	PartNumber string
	PartType   string
	PartTags   string
	PartStatus string
	Reach      string
	EURoHS     string
	ChinaRoHS  string
	Details    []*PartAvailability
}

// PartAvailability 零件库存信息
type PartAvailability struct {
	Manufacturer string
	AvailableQty string
	ShipDate     string
	Quantity     string
	Pricing      string
	// 底边框是否高亮显示
	BorderBottomHighligt bool
}

// search 根据零件编号查询零件标识信息
func search(pn string) (*newDataSet, error) {
	baseURL := "https://www.questcomp.com/baseform.aspx/searchpartsforautocomplete"
	req := bytes.NewBuffer([]byte(`{"pn":"` + pn + `"}`))
	res, err := http.Post(baseURL, "application/json; charset=utf-8", req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	js, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	d := js.Get("d").MustString()
	ds := &newDataSet{}
	err = xml.Unmarshal([]byte(d), ds)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// getDetail 根据网址抓取零件明细信息
func getDetail(url string) (*PartDetail, error) {
	timeout := time.Duration(30 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	var body io.Reader
	request, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	detail := &PartDetail{}
	divPartNumber := document.Find("#divPartNumber")
	detail.PartNumber = divPartNumber.Find("h1").Text()
	detail.PartType = divPartNumber.Find("span").Text()
	detail.PartTags = divPartNumber.Next().Find("h2").Text()
	// detail.PartStatus = document.Find("#divPartStatus").Text()
	// detail.Reach = document.Find("#reach").Text()
	// detail.EURoHS = document.Find("#eurohs").Text()
	// detail.ChinaRoHS = document.Find("#china").Text()
	var pas []*PartAvailability
	m := make(map[string]int)
	document.Find("#MasterPageContent_ucProductHeader_ucPartResults_divPartGrid .rpt-items").Each(func(i int, selection *goquery.Selection) {
		pa := &PartAvailability{}
		if selection.HasClass("rpt-any-mfg") {
			pa.BorderBottomHighligt = true
		}
		if selection.HasClass("instock") {
			m["instock"] = i
		}
		pa.Manufacturer = selection.Find("#MasterPageContent_ucProductHeader_ucPartResults_rptPartResults_divRptPartResultsMfg_" + strconv.Itoa(i)).Text()
		pa.AvailableQty = selection.Find("#rpt-available-qty").Text()
		pa.ShipDate = selection.Find(".rpt-leadtime").Text()
		pa.Quantity, _ = selection.Find("#rpt-range-qty").Html()
		pa.Pricing, _ = selection.Find("#rpt-range-price").Html()
		pas = append(pas, pa)
	})
	// set last instock class div highligt
	if i, ok := m["instock"]; ok {
		pas[i].BorderBottomHighligt = true
	}
	detail.Details = pas
	return detail, nil
}

// 首页处理器
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, nil)
}

// 查询处理器
func searchHandler(w http.ResponseWriter, r *http.Request) {
	pn := r.FormValue("pn")
	nd, err := search(pn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t := template.Must(template.ParseFiles("result.html"))
	err = t.Execute(w, nd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 明细页面处理器
func detailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	pn := r.FormValue("pn")
	pn = strings.ToLower(pn)
	pn = strings.Replace(pn, ":", "", -1)
	t := template.Must(template.ParseFiles("detail.html"))
	url := fmt.Sprintf("https://www.questcomp.com/part/4/%s/%s", pn, id)
	fmt.Println(url)
	detail, err := getDetail(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, detail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/detail", detailHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
