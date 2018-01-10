package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ianjdarrow/slow-hn/models"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

var DB *sqlx.DB

func GetPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestBegin := time.Now()
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
		  p.time,
		  p.type,
		  p.descendants,
		  sum(s.score) as aggregated
		FROM 
		  posts p INNER JOIN scores s ON s.id = p.id
		WHERE
		  s.time BETWEEN ? AND ?
		GROUP BY p.id
		ORDER BY (sum(s.score) + (p.score / 400) + (p.descendants / 100)) DESC
		LIMIT 30;`, start, end)
	resp, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	w.Write(resp)
	requestTime := time.Now().Sub(requestBegin)
	fmt.Printf("Served /posts, total time: %s\n", requestTime)
}

func GetPostsPreflight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	w.Write([]byte(""))
}
