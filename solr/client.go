package solr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ahmadrosid/heline/utils"
)

type Map map[string]interface{}

type SolrQuery struct {
	Query  string
	Filter []string
}

func Search(query SolrQuery) ([]byte, error) {
	// u, _ := url.Parse("https://heline.dev/solr/heline/select")
	u, _ := url.Parse("http://localhost:8984/solr/heline/select")
	q := u.Query()
	q.Set("hl", "on")
	q.Set("hl.fl", "content")
	q.Set("hl.simple.pre", "<mark>")
	q.Set("hl.simple.post", "</mark>")
	q.Set("hl.snippets", "3")
	// q.Set("hl.usePhraseHighlighter", "true")
	// q.Set("hl.requireFieldMatch", "true")
	// q.Set("hl.highlightMultiTerm", "true")
	// q.Set("hl.mergeContiguous", "true")
	q.Set("hl.fragsize", "2500")
	q.Set("hl.maxAnalyzedChars", "100000")
	q.Set("hl.method", "unified")
	u.RawQuery = q.Encode()

	data := Map{
		"query": "content:" + query.Query,
		"fields": "id,file_id,repo,lang,branch,owner_id",
		"facet": Map{
			"lang": Map{
				"type":  "terms",
				"field": "lang",
				"limit": 10,
			},
			"path": Map{
				"type":  "terms",
				"field": "path",
				"limit": 8,
			},
			"repo": Map{
				"type":  "terms",
				"field": "repo",
				"limit": 7,
			},
		},
	}

	if query.Filter != nil {
		data["filter"] = query.Filter
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
