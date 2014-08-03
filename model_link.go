package main

import (
  "gopkg.in/mgo.v2/bson"
  "time"
)

type (
  Links []Link
  Link struct {
    Id bson.ObjectId `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
    Url string `json:"url" bson:"url"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bsom:"created_at"`
    FolderId bson.ObjectId `json:"folder_id" bson:"folder_id"`
    Tags []string `json:"tags" bson:"tags"`
  }
)


func (mc MgoCon) Link_Create (link *Link) (err error) {
  if link.Id.Hex() == "" {
    link.Id = bson.NewObjectId()
  }
  _, err = mc.DB.C("link").UpsertId(link.Id, link)
  return
}


func (mc MgoCon) Link_All(links *Links) (err error) {
  err = mc.DB.C("link").Find(bson.M{}).All(links)
  return
}
