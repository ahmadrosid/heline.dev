package entity

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
