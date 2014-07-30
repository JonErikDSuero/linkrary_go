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

  var mongo_url= os.Getenv("MONGO_URL")
  if (mongo_url == "") {
    mongo_url = "localhost:27017"
  }
  session, err := mgo.Dial(mongo_url)
  if err != nil {
    panic(err)
  }
  defer session.Close()
  // Optional. Switch the session to a monotonic behavior.
  session.SetMode(mgo.Monotonic, true)
  con := session.DB("test").C("link")
  err = con.Insert( &Link{ r.FormValue("name"), r.FormValue("url"), r.FormValue("submitted_by") } )
  if err != nil {
    panic(err)
  }
  fmt.Fprintf(w, "SemiWorking...")
}

// MONGODB_DATABASE: "linkrary_go-production"
// MONGODB_HOST:     "172.17.0.129"
// MONGODB_PORT:     "27017"
// MONGODB_USERNAME: "linkrary_go"
// MONGODB_PASSWORD: "QkhxVE5JSzd3SFRMWW9ZTmV1TElocWNjbHlLVGd0SXlhQ0dEQ3ZJMm8zaz0K"
// MONGO_URL:        "mongodb://linkrary_go:QkhxVE5JSzd3SFRMWW9ZTmV1TElocWNjbHlLVGd0SXlhQ0dEQ3ZJMm8zaz0K@172.17.0.129:27017/linkrary_go-production"

