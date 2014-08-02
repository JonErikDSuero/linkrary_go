package main

import (
  "fmt"
  "strings"
  "sort"
  "io/ioutil"
  "gopkg.in/mgo.v2/bson"
  "github.com/reiver/go-porterstemmer"
)


func (r FolderRepo) Create(folder *Folder) (err error) {
  if folder.Id.Hex() == "" {
    folder.Id = bson.NewObjectId()
  }
  _, err = r.Collection.UpsertId(folder.Id, folder)
  return
}

func (r FolderRepo) All(folders *Folders) (err error) {
  err = r.Collection.Find(bson.M{}).All(&folders)
  return
}

func (r FolderRepo) SuggestFolder(folder_suggested *Folder, extra_info *string) (err error) {
  csv_content, err := ioutil.ReadFile("stopwords.csv")
  if err != nil { panic(err) }

  stopwords := strings.Split(string(csv_content), ",")
  tags_raw := strings.Fields(*extra_info)
  sort.Sort(sort.StringSlice(tags_raw))

  i_stopwords := 0;
  i_tags_raw := 0;
  i_tags_filtered := 0;
  tags_filtered := make([]string, len(tags_raw), len(tags_raw))
  for (i_stopwords < len(stopwords)) && (i_tags_raw < len(tags_raw)) {
    if (tags_raw[i_tags_raw] == stopwords[i_stopwords]) {
      i_tags_raw++
    } else if (tags_raw[i_tags_raw] > stopwords[i_stopwords]) {
      i_stopwords++
    } else {
      tags_filtered[i_tags_filtered] = porterstemmer.StemString(tags_raw[i_tags_raw])
      i_tags_filtered++
      i_tags_raw++
    }
  }

  fmt.Printf("tags_filtered: %s", tags_filtered)
  fmt.Printf("tags_raw: %s", tags_raw)
  err = r.Collection.Find(bson.M{}).One(&folder_suggested)
  return
}
