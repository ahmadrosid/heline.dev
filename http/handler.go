package http

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/ahmadrosid/heline/solr"
	"github.com/ahmadrosid/heline/utils"
	queryparam "github.com/tomwright/queryparam/v4"
)

type Map map[string]interface{}

type Data struct {
	Content []string `json:"content"`
}

type SolrField struct {
	ID      string `json:"id"`
	FileID  string `json:"file_id"`
	OwnerID string `json:"owner_id"`
	Repo    string `json:"repo"`
	Branch  string `json:"branch"`
}

type SolrDoc struct {
	Docs 			[]SolrField `json:"docs"`
	NumFound	int `json:"numFound"`
}

type SolrBuckets []struct {
	Val   string `json:"val"`
	Count int    `json:"count"`
}

type SolrFacet struct {
	Count int `json:"count"`
	Lang  struct {
		Buckets SolrBuckets `json:"buckets"`
	} `json:"lang"`
	Path struct {
		Buckets SolrBuckets `json:"buckets"`
	} `json:"path"`
	Repo struct {
		Buckets SolrBuckets `json:"buckets"`
	} `json:"repo"`
}

type SolrResult struct {
	Highlight map[string]Data `json:"highlighting"`
	Response  SolrDoc         `json:"response"`
	Facet     SolrFacet       `json:"facets"`
}

type ContentData struct {
	ID      Map `json:"id"`
	OwnerID Map `json:"owner_id"`
	FileID  Map `json:"file_id"`
	Branch  Map `json:"branch"`
	Content Map `json:"content"`
	Repo    Map `json:"repo"`
}

type QueryParamFilter struct {
	Repo     []string `queryparam:"repo"`
	Language []string `queryparam:"lang"`
	Path     []string `queryparam:"path"`
}

type QueryParam struct {
	Query string   `queryparam:"q"`
	Path  []string `queryparam:"filter[path]"`
	Lang  []string `queryparam:"filter[lang]"`
	Repo  []string `queryparam:"filter[repo]"`
}

type DocHits struct {
	Hits		[]ContentData `json:"hits"`
	Facets	SolrFacet     `json:"facets"`
	Total 	int     			`json:"total"`
}

type SearchResult struct {
	Response DocHits `json:"hits"`
}

//go:embed dist
//go:embed dist/_next
//go:embed dist/_next/static/chunks/pages/*.js
//go:embed dist/_next/static/*/*.js
var nextFS embed.FS

func Handler() http.Handler {
	index, err := fs.Sub(nextFS, "dist")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(index)))
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		var query = "/search.html"
		if r.URL.Query().Get("q") != "" {
			query = query + "?q=" + r.URL.Query().Get("q")
		}
		
		if r.URL.Query().Get("filter[repo]") != "" {
			query = query + "&filter[repo]=" + r.URL.Query().Get("filter[repo]")
		}
		
		if r.URL.Query().Get("filter[lang]") != "" {
			query = query + "&filter[lang]=" + r.URL.Query().Get("filter[lang]")
		}
		
		if r.URL.Query().Get("filter[path]") != "" {
			query = query + "&filter[path]=" + r.URL.Query().Get("filter[path]")
		}
		
		http.Redirect(w, r, query, http.StatusSeeOther)
	})
	mux.HandleFunc("/api/search", handleSearch)

	return wrapCORSHandler(mux, &CorsConfig{
		allowedOrigin: "*",
	})
}

func getQueryFilter(param QueryParam) []string {
	var filter []string

	if len(param.Lang) > 0 {
		filter = append(filter, fmt.Sprintf("lang:(%s)", utils.Join(param.Lang, " ", "*")))
	}

	if len(param.Path) > 0 {
		filter = append(filter, fmt.Sprintf("path:(%s)", utils.Join(param.Path, " ", "*")))
	}

	if len(param.Repo) > 0 {
		filter = append(filter, fmt.Sprintf("repo:(%s)", utils.Join(param.Repo, " ", "*")))
	}

	return filter
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	param := QueryParam{}
	err := queryparam.Parse(r.URL.Query(), &param)
	switch err {
	case nil:
		break
	case queryparam.ErrInvalidBoolValue:
		println("Failed parse query param")
		return
	default:
		println("return empty query param", err.Error())
		return
	}

	q := strings.Replace(param.Query, "*", "\\*", -1)
	q = strings.Replace(q, "\"", "\\\"", -1)
	q = strings.Replace(q, "'", "\\'", -1)
	q = "\"" + q + "\""

	result, err := solr.Search(solr.SolrQuery{
		Query:  q,
		Filter: getQueryFilter(param),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(Map{
			"error": err.Error(),
		})
		return
	}

	dec := json.NewDecoder(bytes.NewReader(result))
	var data SolrResult
	err = dec.Decode(&data)
	if err != nil {
		enc.Encode(Map{
			"error": err.Error(),
		})
		return
	}

	println("hints:", data.Response.NumFound, q)

	w.Header().Set("Content-Type", "application/json")
	var content []ContentData
	for _, item := range data.Response.Docs {
		contents := data.Highlight[item.ID].Content
		if len(contents) == 0 {
			continue
		}
		content = append(content, ContentData{
			ID: Map{
				"raw": item.ID,
			},
			Branch: Map{
				"raw": item.Branch,
			},
			OwnerID: Map{
				"raw": item.OwnerID,
			},
			FileID: Map{
				"raw": item.FileID,
			},
			Content: Map{
				"snippet": contents,
				// "snippet": contents[len(contents)-1:],
			},
			Repo: Map{
				"raw": item.Repo,
			},
		})
	}
	enc.Encode(SearchResult{
		Response: DocHits{
			Hits:		content,
			Facets:	data.Facet,
			Total:	data.Response.NumFound,
		},
	})
}
