package models

import (
  "fmt"
  "time"
)

type Score struct {
  ID    int       `json:"id"`
  Score float64   `json:"score"`
  Time  time.Time `json:"time"`
}

type Post struct {
  By             string  `json:"by"`
  ID             int     `json:"id"`
  Score          int     `json:"score"`
  Time           int     `json:"time"`
  Title          string  `json:"title"`
  Type           string  `json:"type"`
  URL            string  `json:"url"`
  AggregateScore []Score `json:"aggregate_score"`
}

func (p Post) Summarize() {
  fmt.Printf("Post %v: %s\n", p.ID, p.Title)
}
