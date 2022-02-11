package github

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ahmadrosid/heline/utils"
	"github.com/gocolly/colly/v2"
)

type GithubFile struct {
	ID      string   `json:"id"`
	FileID  string   `json:"file_id"`
	OwnerID string   `json:"owner_id"`
	Path    string   `json:"path"`
	Repo    string   `json:"repo"`
	Branch  string   `json:"branch"`
	Lang    string   `json:"lang"`
	Content []string `json:"content"`
}

func getId(path *url.URL) string {
	paths := strings.Split(path.Path, "/")
	if len(paths) > 5 {
		paths = append(paths[:3], paths[5:]...)
	}

	return "g" + strings.Join(paths, "/")
}

func getPath(path *url.URL) string {
	paths := strings.Split(path.Path, "/")
	if len(paths) > 5 {
		paths = paths[5:]
	}

	return strings.Join(paths, "/")
}

func getBranch(path *url.URL) string {
	fullPath := path.Path
	dir := filepath.Dir(fullPath)
	for i := len(dir) - 1; i >= 0; i-- {
		if dir[i] == '\\' {
			dir = dir[i+1:]
			if len(dir) >= 40 {
				dir = dir[:10]
			}
			break
		}
	}

	return dir
}

func getRepo(path *url.URL) string {
	paths := strings.Split(path.Path, "/")

	return strings.Join(paths[1:3], "/")
}

func getLangFromUrl(path *url.URL) string {
	ext := filepath.Ext(path.Path)

	switch ext {
	case ".clj", ".cljs":
		ext = "Clojure"
	case ".cc", ".c++", "hpp", "h":
		ext = "c++"
	case ".ex":
		ext = "Elixir"
	case ".erl":
		ext = "Erlang"
	case ".hs":
		ext = "Haskell"
	case ".js":
		ext = "JavaScript"
	case ".ts", ".tsx":
		ext = "TypeScript"
	case ".gitignore", ".npmignore", ".dockerignore", ".eslintignore", ".prettierignore":
		ext = "Ignore"
	case ".md":
		ext = "Markdown"
	case ".rs":
		ext = "Rust"
	case ".rb":
		ext = "Ruby"
	case ".scss":
		ext = "SCSS"
	case ".sh":
		ext = "Shell"
	case ".txt":
		ext = "Text"
	case ".json":
		ext = "JSON"
	case ".yaml", ".yml":
		ext = "YAML"
	case ".html":
		ext = "HTML"
	case ".css":
		ext = "CSS"
	case ".php":
		ext = "PHP"
	case ".mod":
		file := filepath.Base(path.Path)
		if file == "go.mod" {
			ext = "Go"
		}
	case ".go":
		ext = "Go"
	default:
		if len(ext) > 1 {
			ext = ext[1:]
		}
	}

	return ext
}

func minifyHtml(source string) string {
	var re = regexp.MustCompile(`(?m)<!--(.*?)-->|\s\B`)
	return re.ReplaceAllString(strings.TrimSpace(source), "")
}

func ExtractFile(h *colly.HTMLElement) GithubFile {
	var file = GithubFile{}
	file.ID = h.Request.URL.Path
	file.FileID = getId(h.Request.URL)
	file.Path = getPath(h.Request.URL)
	file.Lang = getLangFromUrl(h.Request.URL)
	file.Branch = getBranch(h.Request.URL)
	file.Repo = getRepo(h.Request.URL)

	h.ForEachWithBreak("meta[name]", func(i int, e *colly.HTMLElement) bool {
		if utils.ItHas(e.Attr("name"), "octolytics-dimension-user_id") {
			file.OwnerID = e.Attr("content")
		}

		return file.OwnerID == ""
	})

	h.ForEachWithBreak("table.highlight", func(i int, e *colly.HTMLElement) bool {
		var codeSources []string
		e.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			html, _ := tr.DOM.Html()
			htmlContent := "<tr>" + html + "</tr>"
			codeSources = append(codeSources, htmlContent)
		})

		var contentSources = []string{}
		var chunk = 0
		var maxLen = len(codeSources)
		for i := 0; i < maxLen; i++ {
			chunk += 5
			if chunk > maxLen-1 {
				chunk = maxLen - 1
			}
			ctn := `<table class="highlight tab-size js-file-line-container" data-tab-size="8" data-paste-markdown-skip="">
			<tbody>`
			ctn += strings.Join(codeSources[i:chunk], "")
			ctn += `</tbody></table>`
			contentSources = append(contentSources, ctn)
			i = chunk
		}

		file.Content = contentSources
		return i > 0
	})

	return file
}
