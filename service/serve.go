package service

import (
	"fmt"
	"log"
	"net/http"
	"smesearch/indexer"
	"smesearch/searcher"
	"strings"
)

var (
	s_indexDir  string
	s_indexName string
	s_indexer   string
)

func writeResponse(w http.ResponseWriter, r *http.Request) http.ResponseWriter {
	q := r.FormValue("q")
	fmt.Println(q)
	res := searcher.Search(indexer.GetIndex(s_indexDir, s_indexName), strings.Fields(q))
	w.Header().Set("Access-Control-Allow-Origin", "https://kb.local.meyn.fr")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	for _, hit := range res.Hits {
		fmt.Fprintln(w, "<br/>")
		fmt.Fprintln(w, hit.Fragments[""], "<a class='bluebtn' href='https://kb.local.meyn.fr/"+strings.TrimSuffix(hit.ID, ".md")+"/'>", hit.ID, "</a>")
	}
	return w
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-- init --")
	switch r.Method {
	case "GET":
		w = writeResponse(w, r)
		return
	case "POST":
		w = writeResponse(w, r)
		return
	}
	fmt.Println(r)
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func Serve(indexDir string, indexName string) {
	s_indexDir = indexDir
	s_indexName = indexName
	s_indexer = indexer.GetIndex(s_indexDir, s_indexName)
	http.HandleFunc("/", searchHandler)
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServeTLS(":8030", "search.crt", "search.key", nil))
}
