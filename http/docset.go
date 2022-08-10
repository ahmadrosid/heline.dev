package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ahmadrosid/heline/core/entity"
	"github.com/ahmadrosid/heline/core/module/solr"
	"github.com/ahmadrosid/heline/core/utils"
)

func handleGetDocsetByID(w http.ResponseWriter, id string) {
	result, err := solr.GetDocsetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, entity.Map{
			"error": err.Error(),
		})
		return
	}

	dec := json.NewDecoder(bytes.NewReader(result))
	var data entity.DocsetDetail
	err = dec.Decode(&data)
	if err != nil {
		utils.Encode(w, entity.Map{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(data.Response.Docs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		utils.Encode(w, entity.Map{
			"error": id + " not found!",
		})
		return
	}

	utils.Encode(w, data.Response.Docs[0])
}

func handleSearchDocset(w http.ResponseWriter, q string) {
	result, err := solr.DocsetSearch(solr.DocsetQuery{Query: q})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, entity.Map{
			"error": err.Error(),
		})
		return
	}

	dec := json.NewDecoder(bytes.NewReader(result))
	var data entity.DocsetSolrResult
	err = dec.Decode(&data)
	if err != nil {
		utils.Encode(w, entity.Map{
			"error": err.Error(),
		})
		return
	}

	println("hits:", data.Response.NumFound, q)

	w.Header().Set("Content-Type", "application/json")
	var content []entity.DocsetData
	for _, item := range data.Response.Docs {
		contents := data.Highlight[item.ID].Content
		if len(contents) == 0 {
			continue
		}
		content = append(content, entity.DocsetData{
			ID: entity.Map{
				"raw": item.ID,
			},
			Title: entity.Map{
				"raw": item.Title,
			},
			FileName: entity.Map{
				"raw": item.FileName,
			},
			Link: entity.Map{
				"raw": item.Link,
			},
			Content: entity.Map{
				"snippet": contents,
			},
			Document: entity.Map{
				"raw": item.Document,
			},
		})
	}

	utils.Encode(w, entity.DocsetSearchResult{
		Response: entity.DocsetHits{
			Hits:   content,
			Facets: data.Facet,
			Total:  data.Response.NumFound,
		},
	})
}
