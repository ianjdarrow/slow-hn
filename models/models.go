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
	By              string  `json:"by"`
	ID              int     `json:"id"`
	Score           int     `json:"score"`
	Time            int     `json:"time"`
	Title           string  `json:"title"`
	Type            string  `json:"type"`
	URL             string  `json:"url"`
	AggregateScores []Score `json:"aggregate_scores"`
}

func (p Post) AddScore(score Score) Post {
	p.AggregateScores = append(p.AggregateScores, score)
	return p
}

func (p Post) Summarize() {
	var totalScore float64
	for _, score := range p.AggregateScores {
		totalScore += score.Score
	}
	fmt.Printf("Post %v: %s\nScores: %.2f\n", p.ID, p.Title, totalScore)
}
