package main

import (
  "fmt"
  "gopkg.in/mgo.v2/bson"
)


func (r FolderRepo) Create(folder *Folder) (err error) {
  if folder.Id.Hex() == "" {
    folder.Id = bson.NewObjectId()
  }
  _, err = r.Collection.UpsertId(folder.Id, folder)
  return
}

func (r FolderRepo) All(folders *Folders) (err error) {
  err = r.Collection.Find(bson.M{}).All(folders)
  return
}

func (r FolderRepo) SuggestFolder(folder_suggested *Folder, linkRepo *LinkRepo, tags_filtered *[]string, extra_info *string) (err error) {
  //TODO: What if there are no Links???
  var links Links
  var best_score int
  var best_folder_id bson.ObjectId

  *tags_filtered = TagsFilter(extra_info) // remove empty strings
  fmt.Println(tags_filtered)
  if err = linkRepo.All(&links); err != nil {
    panic(err)
  }

  scores := make(map[bson.ObjectId]int)
  for _, link := range links {
    fmt.Println("link: ",link.Name)
    scores[link.FolderId] += TagsCommonalityScore(link.Tags, *tags_filtered)
    fmt.Println("score: ",scores[link.FolderId])
  }
  fmt.Println(scores)
  for folder_id, score := range scores {
    if (best_score < score) {
      best_score = score
      best_folder_id = folder_id
    }
  }
  r.Collection.Find(bson.M{"_id": best_folder_id}).One(&folder_suggested)
  fmt.Printf("suggested folder %s\n", folder_suggested)
  return
}

