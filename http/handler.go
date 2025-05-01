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

	"github.com/ahmadrosid/heline/core/entity"
	"github.com/ahmadrosid/heline/core/module/solr"
	"github.com/ahmadrosid/heline/core/utils"
	queryparam "github.com/tomwright/queryparam/v4"
)

//go:embed dist
//go:embed dist/_next
//go:embed dist/_next/static/chunks/pages/*.js
//go:embed dist/_next/static/*/*.js
var nextFS embed.FS

func Handler(analytic http.Handler) http.Handler {
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
	mux.Handle("/stats", analytic)

	return wrapCORSHandler(mux, &CorsConfig{
		allowedOrigin: "*",
	})
}

func getQueryFilter(param entity.QueryParam) []string {
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
	param := entity.QueryParam{}
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
	// q = "*" + q + "*"
	println(q)

	if param.ID != "" && param.Tbm == "docs" {
		handleGetDocsetByID(w, param.ID)
		return
	}

	if param.Tbm == "docs" {
		handleSearchDocset(w, q)
		return
	}

	result, err := solr.Search(solr.SolrQuery{
		Query:  q,
		Filter: getQueryFilter(param),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(entity.Map{
			"error": err.Error(),
		})
		return
	}

	dec := json.NewDecoder(bytes.NewReader(result))
	var data entity.SolrResult
	err = dec.Decode(&data)
	if err != nil {
		enc.Encode(entity.Map{
			"error": err.Error(),
		})
		return
	}

	println("hints:", data.Response.NumFound, q)

	w.Header().Set("Content-Type", "application/json")
	var content []entity.ContentData
	for _, item := range data.Response.Docs {
		contents := data.Highlight[item.ID].Content
		if len(contents) == 0 {
			continue
		}
		content = append(content, entity.ContentData{
			ID: entity.Map{
				"raw": item.ID,
			},
			Branch: entity.Map{
				"raw": item.Branch,
			},
			OwnerID: entity.Map{
				"raw": item.OwnerID,
			},
			FileID: entity.Map{
				"raw": item.FileID,
			},
			Content: entity.Map{
				"snippet": contents,
				// "snippet": contents[len(contents)-1:],
			},
			Repo: entity.Map{
				"raw": item.Repo,
			},
		})
	}
	enc.Encode(entity.CodeSearchResult{
		Response: entity.CodeHits{
			Hits:   content,
			Facets: data.Facet,
			Total:  data.Response.NumFound,
		},
	})
}
