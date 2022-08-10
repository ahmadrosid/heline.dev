package solr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ahmadrosid/heline/utils"
)

type DocsetQuery struct {
	Query string
}

func GetDocsetByID(id string) ([]byte, error) {
	u, _ := url.Parse("http://localhost:8984/solr/docset/get")
	q := u.Query()
	q.Set("ids", id)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println("ERROR", err.Error())
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}

func DocsetSearch(query DocsetQuery) ([]byte, error) {
	u, _ := url.Parse("http://localhost:8984/solr/docset/select")
	q := u.Query()
	q.Set("hl", "on")
	q.Set("hl.fl", "content")
	q.Set("hl.simple.pre", "<mark>")
	q.Set("hl.simple.post", "</mark>")
	q.Set("hl.snippets", "2")
	// q.Set("hl.usePhraseHighlighter", "true")
	// q.Set("hl.requireFieldMatch", "true")
	// q.Set("hl.highlightMultiTerm", "true")
	// q.Set("hl.mergeContiguous", "true")
	q.Set("hl.fragsize", "500")
	// q.Set("hl.maxAnalyzedChars", "100000")
	// q.Set("hl.method", "unified")
	u.RawQuery = q.Encode()

	data := Map{
		"query":  "content:" + query.Query,
		"fields": "id,file_name,title,link,document",
		"facet": Map{
			"document": Map{
				"type":  "terms",
				"field": "document",
				"limit": 10,
			},
		},
	}

	queryData, _ := json.Marshal(data)

	println("solr_query:", utils.ByteToString(queryData))
	println("raw_query:", u.RawQuery)

	payload := bytes.NewReader(queryData)

	req, _ := http.NewRequest("POST", u.String(), payload)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		println("ERROR", err.Error())
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body, nil
}
