package searcher

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"strings"
)

func Search(indexBase string, terms []string) {
	index, _ := bleve.Open(indexBase)
	query := bleve.NewQueryStringQuery(strings.Join(terms, " "))
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := index.Search(searchRequest)
	for res, hit := range searchResult.Hits {
		fmt.Println(res, "https://kb.local.meyn.fr/"+hit.ID, hit.Score)
	}
}
