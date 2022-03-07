package http

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
	Docs     []SolrField `json:"docs"`
	NumFound int         `json:"numFound"`
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
	Tbm   string   `queryparam:"tbm"`
	ID    string   `queryparam:"id"`
	Path  []string `queryparam:"filter[path]"`
	Lang  []string `queryparam:"filter[lang]"`
	Repo  []string `queryparam:"filter[repo]"`
}

type CodeHits struct {
	Hits   []ContentData `json:"hits"`
	Facets SolrFacet     `json:"facets"`
	Total  int           `json:"total"`
}

type CodeSearchResult struct {
	Response CodeHits `json:"hits"`
}

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
			Content  []string `json:"content"`
		} `json:"docs"`
	} `json:"response"`
}
