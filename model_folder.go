package main

import (
  "gopkg.in/mgo.v2/bson"
)

type (
  Folders []Folder
  Folder struct {
    Id bson.ObjectId `json:"id" bson:"_id"`
    Name string `json:"name" bson:"name"`
  }
)


func (mc MgoCon) Folder_Create(folder *Folder) (err error) {
  if folder.Id.Hex() == "" {
    folder.Id = bson.NewObjectId()
  }
  _, err = mc.DB.C("folder").UpsertId(folder.Id, folder)
  return
}


func (mc MgoCon) Folder_All(folders *Folders) (err error) {
  err = mc.DB.C("folder").Find(bson.M{}).All(folders)
  return
}


func (mc MgoCon) Folder_Suggest(folder_suggested *Folder, tags_filtered *[]string, extra_info *string) (err error) {
  var links Links
  var best_score int
  var best_folder_id bson.ObjectId

  *tags_filtered = Tag_Filter(extra_info) // remove empty strings

  if err = mc.Link_All(&links); err != nil {
    panic(err)
  }

  scores := make(map[bson.ObjectId]int)
  for _, link := range links {
    scores[link.FolderId] += Tag_CommonalityScore(link.Tags, *tags_filtered)
  }

  for folder_id, score := range scores {
    if (best_score < score) {
      best_score = score
      best_folder_id = folder_id
    }
  }

  mc.DB.C("folder").Find(bson.M{"_id": best_folder_id}).One(&folder_suggested)
  return err
}
