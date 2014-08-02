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
  "time"
)

type (
  LinkRepo struct {
    Collection *mgo.Collection
  }

  FolderRepo struct {
    Collection *mgo.Collection
  }

  Links []Link
  Link struct {
    Id bson.ObjectId `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
    Url string `json:"url" bson:"url"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bsom:"created_at"`
    FolderId bson.ObjectId `json:"folder_id" bson:"folder_id"`
  }

  Folders []Folder
  Folder struct {
    Id bson.ObjectId `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
    Tags []string
  }
)

var linkRepo LinkRepo
var folderRepo FolderRepo

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
  folderRepo.Collection = session.DB("linkrary_go-production").C("folder")

  r := mux.NewRouter()
  r.HandleFunc("/", handleHome).Methods("GET")
  r.HandleFunc("/links", handleLinkAll).Methods("GET")
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
  var link Link
  var folder_suggested Folder
  var err  error
  data := struct {
    Success bool `json:"success"`
    FolderName string `json:"folder_name"`
  }{
    Success: false,
  }
  r.ParseForm()
  link.Name = r.FormValue("name")
  link.Url = r.FormValue("url")
  extra_info := link.Name + " " + r.FormValue("extra_info")
  folderRepo.SuggestFolder(&folder_suggested, &extra_info)
  link.FolderId = folder_suggested.Id
  if err = linkRepo.Create(&link); err != nil {
    panic(err)
  } else {
    data.Success = true
    data.FolderName = folder_suggested.Name
  }
  writeJson(w, data)
}


func handleLinkAll(w http.ResponseWriter, r *http.Request) {
  var (
    links Links
    err   error
  )
  if err = linkRepo.All(&links); err != nil {
    log.Printf("%v", err)
    http.Error(w, "500 Internal Server Error", 500)
    return
  }
  writeJson(w, links)
}


