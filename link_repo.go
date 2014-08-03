package main

import (
  "gopkg.in/mgo.v2/bson"
)


func (r LinkRepo) Create (link *Link) (err error) {
  if link.Id.Hex() == "" {
    link.Id = bson.NewObjectId()
  }
  _, err = r.Collection.UpsertId(link.Id, link)
  return
}

func (r LinkRepo) All(links *Links) (err error) {
  err = r.Collection.Find(bson.M{}).All(links)
  return
}
