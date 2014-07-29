package main

import (
  "os"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler).Methods("GET")
  http.Handle("/", r)

  var port = os.Getenv("PORT")
  if (port == "") {
    port = "3000"
  }

  err := http.ListenAndServe(":"+port, nil)
  if err != nil {
    panic(err)
  }
}