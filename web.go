package main


import (
  "os"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  //"strings"
  "gopkg.in/mgo.v2/bson"
  "log"
)

type (
  mgoRepo struct {
    Collection *mgo.Collection
  }

  Links []Link
  Link struct {
    Id           bson.ObjectId  `json:"id"            bson:"_id"`
    Name         string         `json:"name"          bson:"name"`
    Url          string         `json:"url"           bson:"url"`
    SubmittedBy  string         `json:"submitted_by"  bson:"submitted_by"`
  }
)


var linkRepo mgoRepo


func main() {
  var mongo_url = os.Getenv("MY_MONGODB_URL")
  if (mongo_url == "") {
    mongo_url = "localhost:27017"
  }
  session, err := mgo.Dial(mongo_url)
  if err != nil { panic(err) }
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)  // Optional. Switch the session to a monotonic behavior.
  linkRepo.Collection = session.DB("linkrary_go-production").C("link")

  r := mux.NewRouter()
  r.HandleFunc("/", handleHome).Methods("GET")
  r.HandleFunc("/links/create", handleLinkCreate).Methods("POST")
  http.Handle("/", r)
  var port = os.Getenv("PORT")
  if (port == "") {
    port = "3000"
  }
  err = http.ListenAndServe(":"+port, nil)
  if err != nil { panic(err) }
}


func handleHome(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello!")
}


func handleLinkCreate(w http.ResponseWriter, r *http.Request) {
  var (
    link Link
    err  error
  )
  data := struct {
    Success bool `json:"success"`
    Link    Link `json:"link"`
  }{
    Success: false,
  }
  r.ParseForm()
  link.Name = r.FormValue("name")
  link.Url = r.FormValue("url")
  link.SubmittedBy = r.FormValue("submitted_by")
  if err = linkRepo.Create(&link); err != nil {
    panic(err)
  } else {
    data.Success = true
    data.Link = link
  }

  writeJson(w, data)
}


func (r mgoRepo) Create(link *Link) (err error) {
  if link.Id.Hex() == "" {
    link.Id = bson.NewObjectId()
  }
  //if link.Created.IsZero() {
  //  link.Created = time.Now()
  //}
  //link.Updated = time.Now()
  _, err = r.Collection.UpsertId(link.Id, link)
  return
}


func handleLinks(w http.ResponseWriter, r *http.Request) {
  var (
    links Links
    err   error
  )
  if links, err = linkRepo.All(); err != nil {
    log.Printf("%v", err)
    http.Error(w, "500 Internal Server Error", 500)
    return
  }
  writeJson(w, links)
}

func (r mgoRepo) All() (links Links, err error) {
  err = r.Collection.Find(bson.M{}).All(&links)
  return
}



