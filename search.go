package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"

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

// search 根据零件编号查询零件标识信息
func search(pn string) (*newDataSet, error) {
	url := "https://www.questcomp.com/baseform.aspx/searchpartsforautocomplete"
	json := []byte(`{"pn":"` + pn + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	js, err := simplejson.NewFromReader(resp.Body)
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

// 查询处理器
func searchHandler(w http.ResponseWriter, r *http.Request) {
	pn := r.FormValue("pn")
	ds, err := search(pn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ds)
}
