package searcher

import (
	"fmt"
	bleve "github.com/blevesearch/bleve/v2"
	"os"
	"strings"
)

func Search(indexBase string, terms []string) *bleve.SearchResult {
	index, _ := bleve.Open(indexBase)
	defer index.Close()
	query := bleve.NewQueryStringQuery(strings.Join(terms, " "))
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchRequest.Highlight.AddField("")
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return searchResult
}
