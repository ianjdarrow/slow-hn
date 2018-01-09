package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ianjdarrow/slow-hn/models"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

var DB *sqlx.DB

func GetPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	start := queryValues.Get("start")
	end := queryValues.Get("end")
	var posts []models.Post
	err := DB.Select(&posts,
		`SELECT
		  p.id,
		  p.title,
		  p.score,
		  p.by,
		  p.url,
		  sum(s.score) as aggregated
		FROM 
		  posts p INNER JOIN scores s ON s.id = p.id
		WHERE
		  s.time BETWEEN ? AND ?
		GROUP BY p.id
		ORDER BY sum(s.score) DESC
		LIMIT 20;`, start, end)
	resp, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write(resp)
}

func GetPostsPreflight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Write([]byte(""))
}
