package main

import (
  "os"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
)

type (
  MgoCon struct {
    DB *mgo.Database
  }
)


func main() {
  var mc MgoCon
  var err error
  mc.DB = MgoCon_Connect("linkrary_go-production")

  r := mux.NewRouter()
  r.HandleFunc("/", Handle_Home(mc)).Methods("GET")
  r.HandleFunc("/links", Handle_LinkAll(mc)).Methods("GET")
  r.HandleFunc("/links/create", Handle_LinkCreate(mc)).Methods("POST")
  http.Handle("/", r)
  var port = os.Getenv("PORT")
  if (port == "") {
    port = "3000"
  }
  err = http.ListenAndServe(":"+port, nil)
  if err != nil { panic(err) }
}


func MgoCon_Connect(db_name string) (db *mgo.Database) {
  var mongo_url = os.Getenv("MY_MONGOmc.DB.URL")
  if (mongo_url == "") {
    mongo_url = "localhost:27017"
  }
  session, err := mgo.Dial(mongo_url)
  if err != nil { panic(err) }
  session.SetMode(mgo.Monotonic, true)  // Optional. Switch the session to a monotonic behavior.
  db = session.DB(db_name)
  return db
}


func Handle_Home(mc MgoCon) (func(http.ResponseWriter, *http.Request)) {
  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello!")
  }
}


func Handle_LinkAll(mc MgoCon) (func(http.ResponseWriter, *http.Request)) {
  return func(w http.ResponseWriter, r *http.Request) {
    var (
      links Links
      err   error
    )
    if err = mc.Link_All(&links); err != nil {
      panic(err)
    }
    writeJson(w, links)
  }
}


func Handle_LinkCreate(mc MgoCon) (func(http.ResponseWriter, *http.Request)) {
  return func(w http.ResponseWriter, r *http.Request) {
    var link Link
    var folder_suggested Folder
    var err  error
    var tags_filtered []string
    data := struct {
      Success bool `json:"success"`
      FolderName string `json:"folder_name"`
    }{
      Success: false,
    }

    r.ParseForm()
    info_raw := r.FormValue("name") + " " + r.FormValue("extra_info")

    if err = mc.Folder_Suggest(&folder_suggested, &tags_filtered, &info_raw); err != nil {
      panic(err)
    }

    link.Name = r.FormValue("name")
    link.Url = r.FormValue("url")
    link.FolderId = folder_suggested.Id
    link.Tags = tags_filtered

    if err = mc.Link_Create(&link); err != nil {
      panic(err)
    } else {
      data.Success = true
      data.FolderName = folder_suggested.Name
    }
    writeJson(w, data)
  }
}
