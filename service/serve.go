package service

import (
	"log"
	"net/http"
	"smesearch/searcher"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searcher.Search()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return
}

func Serve() {
	http.HandleFunc("/search", searchHandler)
	log.Fatal(http.ListenAndServeTLS(":8030", "search.crt", "search.key", nil))
}
