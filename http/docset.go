package http

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ahmadrosid/heline/core/entity"
	"github.com/ahmadrosid/heline/core/module/solr"
	"github.com/ahmadrosid/heline/core/utils"
)

// Docsets: models

type DocsetSolrResult struct {
	Highlight map[string]entity.Data `json:"highlighting"`
	Response  DocsetSolrDoc          `json:"response"`
	Facet     DocsetSolrFacet        `json:"facets"`
}

type DocsetSolrDoc struct {
	Docs     []DocsetSolrField `json:"docs"`
	NumFound int               `json:"numFound"`
}

type DocsetSolrField struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	FileName string `json:"file_name"`
	Document string `json:"document"`
	Link     string `json:"link"`
}

type DocsetSolrFacet struct {
	Count    int `json:"count"`
	Document struct {
		Buckets entity.SolrBuckets `json:"buckets"`
	} `json:"document"`
}

type DocsetHits struct {
	Hits   []DocsetData    `json:"hits"`
	Facets DocsetSolrFacet `json:"facets"`
	Total  int             `json:"total"`
}

type DocsetData struct {
	ID       entity.Map `json:"id"`
	Title    entity.Map `json:"title"`
	FileName entity.Map `json:"file_name"`
	Document entity.Map `json:"document"`
	Content  entity.Map `json:"content"`
	Link     entity.Map `json:"link"`
}

type DocsetSearchResult struct {
	Response DocsetHits `json:"docs"`
}

// Docset by ID
type DocsetDetail struct {
	Response struct {
		Docs []struct {
			ID       string   `json:"id"`
			FileName string   `json:"file_name"`
			Document string   `json:"document"`
			Title    string   `json:"title"`
			Link     string   `json:"link"`
			Content  []string `json:"content"`
		} `json:"docs"`
	} `json:"response"`
}

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
	var data DocsetDetail
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
	var data DocsetSolrResult
	err = dec.Decode(&data)
	if err != nil {
		utils.Encode(w, entity.Map{
			"error": err.Error(),
		})
		return
	}

	println("hits:", data.Response.NumFound, q)

	w.Header().Set("Content-Type", "application/json")
	var content []DocsetData
	for _, item := range data.Response.Docs {
		contents := data.Highlight[item.ID].Content
		if len(contents) == 0 {
			continue
		}
		content = append(content, DocsetData{
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

	utils.Encode(w, DocsetSearchResult{
		Response: DocsetHits{
			Hits:   content,
			Facets: data.Facet,
			Total:  data.Response.NumFound,
		},
	})
}
