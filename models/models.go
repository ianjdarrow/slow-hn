package models

import (
	"fmt"
)

type Score struct {
	ID    int     `json:"id"`
	Score float64 `json:"score"`
	Time  int64   `json:"time"`
}

type Post struct {
	By         string  `json:"by"`
	ID         int     `json:"id"`
	Score      int     `json:"score"`
	Time       int     `json:"time"`
	Title      string  `json:"title"`
	Type       string  `json:"type"`
	URL        string  `json:"url"`
	Aggregated float64 `json:"aggregated_score"`
}

func (p Post) Summarize() {
	fmt.Printf("Post %v: %s\n", p.ID, p.Title)
}
