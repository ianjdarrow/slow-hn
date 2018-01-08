package controllers

import (
  "net/http"

  "github.com/julienschmidt/httprouter"

  "github.com/ianjdarrow/slow-hn/models"
)

var AllPosts = make(map[int]models.Post)

func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  w.Write([]byte("Hello!"))
}
