package entity

// Docsets: models

type DocsetSolrResult struct {
	Highlight map[string]Data `json:"highlighting"`
	Response  DocsetSolrDoc   `json:"response"`
	Facet     DocsetSolrFacet `json:"facets"`
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
		Buckets SolrBuckets `json:"buckets"`
	} `json:"document"`
}

type DocsetHits struct {
	Hits   []DocsetData    `json:"hits"`
	Facets DocsetSolrFacet `json:"facets"`
	Total  int             `json:"total"`
}

type DocsetData struct {
	ID       Map `json:"id"`
	Title    Map `json:"title"`
	FileName Map `json:"file_name"`
	Document Map `json:"document"`
	Content  Map `json:"content"`
	Link     Map `json:"link"`
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
