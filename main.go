package main

import (
    "flag"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"
)

var tpl = template.Must(template.ParseFiles("index.html"))
var apiKey *string

type Source struct {
    ID   interface{} `json:"id"`
    Name string      `json:"name"`
}

type Article struct {
    Source      Source    `json:"source"`
    Author      string    `json:"author"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    URL         string    `json:"url"`
    URLToImage  string    `json:"urlToImage"`
    PublishedAt time.Time `json:"publishedAt"`
    Content     string    `json:"content"`
}

type Results struct {
    Status       string    `json:"status"`
    TotalResults int       `json:"totalResults"`
    Articles     []Article `json:"articles"`
}

type Search struct {
    SearchKey  string
    NextPage   int
    TotalPages int
    Results    Results
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    u, err := url.Parse(r.URL.String())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal server error"))
        return
    }

    params := u.Query()
    searchKey := params.Get("q")
    page := params.Get("page")
    if page == "" {
        page = "1"
    }

    fmt.Println("Search Query is: ", searchKey)
    fmt.Println("Results page is: ", page)
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    apiKey = flag.String("apikey", "", "Newsapi.org access key")
    flag.Parse()

    if *apiKey == "" {
        log.Fatal("apiKey must be set")
    }

    mux := http.NewServeMux()

    fs := http.FileServer(http.Dir("assets"))
    mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

    mux.HandleFunc("/search", searchHandler)
    mux.HandleFunc("/", indexHandler)
    http.ListenAndServe(":"+port, mux)
}