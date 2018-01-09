package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ianjdarrow/slow-hn/models"
)

var (
	topPostsUrl string = "https://hacker-news.firebaseio.com/v0/topstories.json"
	numPosts    int    = 50
)

func FetchTopPosts() []models.Post {
	body := getHTML(topPostsUrl)
	var topPostNumbers []int
	err := json.Unmarshal(body, &topPostNumbers)
	if err != nil {
		fmt.Println("error:", err)
	}
	var posts []models.Post
	for i := 0; i < numPosts; i++ {
		urlString := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", strconv.Itoa(topPostNumbers[i]))
		postBytes := getHTML(urlString)
		var post models.Post
		err := json.Unmarshal(postBytes, &post)
		if err != nil {
			fmt.Println("error:", err)
		}
		posts = append(posts, post)
	}
	return posts
}

func ScorePosts(posts []models.Post) []models.Score {
	var scores []models.Score
	now := time.Now().Unix()
	for i, post := range posts {
		score := models.Score{
			ID:    post.ID,
			Time:  now,
			Score: 10 / (math.Sqrt(float64(i+1)) + 1),
		}
		scores = append(scores, score)
	}
	return scores
}

func UpdatePosts(db *sqlx.DB) {
	posts := FetchTopPosts()
	scores := ScorePosts(posts)

	postTx := db.MustBegin()
	for _, post := range posts {
		postTx.MustExec(`
      REPLACE INTO posts(by, id, score, time, title, type, url)
      VALUES (?, ?, ?, ?, ?, ?, ?);`,
			post.By, post.ID, post.Score, post.Time, post.Title, post.Type, post.URL)
	}
	err := postTx.Commit()
	if err != nil {
		fmt.Printf("Database write error: %s\n", err)
	}

	scoreTx := db.MustBegin()
	for _, score := range scores {
		scoreTx.MustExec(`
      INSERT INTO scores(id, score, time)
      VALUES(?, ?, ?);`,
			score.ID, score.Score, score.Time)
	}
	err = scoreTx.Commit()
	if err != nil {
		fmt.Printf("Database write error: %s\n", err)
	}

	fmt.Printf("Updated posts at %s\n", time.Now())
	fmt.Printf("Top post: %s\n", posts[0].Title)
	time.Sleep(60 * time.Minute)
	UpdatePosts(db)
}

func getHTML(url string) []byte {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
