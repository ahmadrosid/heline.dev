package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ahmadrosid/heline/solr"
)

func handleGetDocsetByID(w http.ResponseWriter, id string) {
	enc := json.NewEncoder(w)

	result, err := solr.GetDocsetByID(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(Map{
			"error": err.Error(),
		})
		return
	}

	dec := json.NewDecoder(bytes.NewReader(result))
	var data DocsetDetail
	err = dec.Decode(&data)
	if err != nil {
		enc.Encode(Map{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(data.Response.Docs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(Map{
			"error": id + " not found!",
		})
		return
	}

	enc.Encode(data.Response.Docs[0])
}

func handleSearchDocset(w http.ResponseWriter, q string) {
	enc := json.NewEncoder(w)
	result, err := solr.DocsetSearch(solr.DocsetQuery{Query: q})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(Map{
			"error": err.Error(),
		})
		return
	}

	dec := json.NewDecoder(bytes.NewReader(result))
	var data DocsetSolrResult
	err = dec.Decode(&data)
	if err != nil {
		enc.Encode(Map{
			"error": err.Error(),
		})
		return
	}

	println("hints:", data.Response.NumFound, q)

	w.Header().Set("Content-Type", "application/json")
	var content []DocsetData
	for _, item := range data.Response.Docs {
		contents := data.Highlight[item.ID].Content
		if len(contents) == 0 {
			continue
		}
		content = append(content, DocsetData{
			ID: Map{
				"raw": item.ID,
			},
			Title: Map{
				"raw": item.Title,
			},
			FileName: Map{
				"raw": item.FileName,
			},
			Content: Map{
				"snippet": contents,
			},
			Document: Map{
				"raw": item.Document,
			},
		})
	}
	enc.Encode(DocsetSearchResult{
		Response: DocsetHits{
			Hits:   content,
			Facets: data.Facet,
			Total:  data.Response.NumFound,
		},
	})
}
