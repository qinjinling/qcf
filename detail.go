package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PartDetail 零件明细
type PartDetail struct {
	PartNumber  string
	Category    string
	Description string
	Details     []*PartAvailability
}

// PartAvailability 零件库存信息
type PartAvailability struct {
	Manufacturer             string
	AvailableQty             string
	ShipDate                 string
	Quantity                 string
	Pricing                  string
	IsAnyManufacturer        bool
	IsAnyManufacturerInStock bool
	BorderBottomHighligt     bool
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
	detail.Category = divPartNumber.Find("span").Text()
	detail.Description = divPartNumber.Next().Find("h2").Text()
	var pas []*PartAvailability
	m := make(map[string]int)
	document.Find("#MasterPageContent_ucProductHeader_ucPartResults_divPartGrid .rpt-items").Each(func(i int, selection *goquery.Selection) {
		pa := &PartAvailability{}
		if selection.HasClass("rpt-any-mfg") {
			pa.IsAnyManufacturer = true
			pa.BorderBottomHighligt = true
		}
		if selection.HasClass("instock") {
			pa.IsAnyManufacturerInStock = true
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

// 明细页面处理器
func detailHandler(w http.ResponseWriter, r *http.Request) {
	mpid := r.FormValue("mpid")
	pn := r.FormValue("pn")
	pn = strings.ToLower(pn)
	url := fmt.Sprintf("https://www.questcomp.com/part/4/%s/%s", pn, mpid)
	detail, err := getDetail(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}
