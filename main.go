package main

import (
	"fmt"
	"net/http"
	"time"

	f "groupie-tracker-filter/functions"
)

func main() {
	start := time.Now()
	f.FitchAllData()
	fmt.Println(time.Since(start))
	http.HandleFunc("/styles/", f.ServeStyle)
	http.HandleFunc("/", f.FirstPage)
	http.HandleFunc("/suggest", f.SuggestHandler)
	http.HandleFunc("/filtre", f.FilterHandler)
	http.HandleFunc("/artist", f.OtherPages)
	http.HandleFunc("/search", f.SearchPage)
	fmt.Println("http://localhost:5841")
	http.ListenAndServe(":5841", nil)
}
