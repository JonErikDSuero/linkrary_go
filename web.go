package main


import (
  "os"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "strings"
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
  for _, e := range os.Environ() {
    pair := strings.Split(e, "=")
    fmt.Println( pair[0] + "=" + pair[1] )
  }
  var mongo_url = os.Getenv("MONGO_URL")
  if (mongo_url == "") {
    mongo_url = "localhost:27017"
  }
  fmt.Printf("mongo_url is [%s]\n", mongo_url)
  session, err := mgo.Dial(mongo_url)
  if err != nil {
    panic(err)
  }
  defer session.Close()
}


func CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  var mongo_url = r.FormValue("mongo_url")
  if (mongo_url == "") {
    mongo_url = "localhost:27017"
  }
  fmt.Printf("mongo_url !%s!\n", mongo_url)
  session, err := mgo.Dial(mongo_url)
  if err != nil {
    panic(err)
  }
  defer session.Close()
  // Optional. Switch the session to a monotonic behavior.
  session.SetMode(mgo.Monotonic, true)
  con := session.DB("linkrary_go-production").C("link")
  err = con.Insert( &Link{ r.FormValue("name"), r.FormValue("url"), r.FormValue("submitted_by") } )
  if err != nil {
    panic(err)
  }
  fmt.Fprintf(w, "SemiWorking...")
}