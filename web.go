package main


import (
  "os"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  //"gopkg.in/mgo.v2/bson"
)


type Link struct {
  Name string
  Url string
  SubmittedBy string
}


func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", HomeHandler).Methods("GET")
  r.HandleFunc("/links/create", CreateLinkHandler).Methods("POST")
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


func HomeHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}


func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  fmt.Printf("name: !%s!\n", r.FormValue("name"))
  //fmt.Printf("Vars: !%s!\n", vars)
  session, err := mgo.Dial("localhost:27017")
  if err != nil {
    panic(err)
  }
  defer session.Close()
  // Optional. Switch the session to a monotonic behavior.
  session.SetMode(mgo.Monotonic, true)
  con := session.DB("test").C("link")
  err = con.Insert( &Link{ r.FormValue("name"), r.FormValue("url"), r.FormValue("submitted_by") } )
  //fmt.Printf("Name: !%s!\n", vars["name"])
  if err != nil {
    panic(err)
  }
  fmt.Fprintf(w, "SemiWorking...")
}
