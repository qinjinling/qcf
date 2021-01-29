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
