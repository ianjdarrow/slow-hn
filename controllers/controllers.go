package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/ianjdarrow/slow-hn/models"
)

func GetIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts := make([]models.Post, 0)
	for _, post := range AllPosts {
		posts = append(posts, post)
	}
	result, _ := json.Marshal(posts)
	w.Write(result)
}
