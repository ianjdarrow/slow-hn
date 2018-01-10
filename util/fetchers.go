package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
			Score: float64(10.0 / ((i + 1) + 1)),
		}
		scores = append(scores, score)
	}
	return scores
}

func UpdatePosts(db *sqlx.DB) {
	posts := FetchTopPosts()

	postTx := db.MustBegin()
	for _, post := range posts {
		postTx.MustExec(`
      REPLACE INTO posts(by, id, score, time, title, type, url, descendants)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
			post.By, post.ID, post.Score, post.Time, post.Title, post.Type, post.URL, post.Descendants)
	}
	err := postTx.Commit()
	if err != nil {
		fmt.Printf("Database write error: %s\n", err)
	}

	var lastUpdate int64
	updateCheck := `SELECT max(time) FROM scores;`
	err = db.Get(&lastUpdate, updateCheck)
	if err != nil {
		fmt.Printf("Error checking last update time: %s\n", err)
	}
	now := time.Now().Unix()
	if now-lastUpdate > 60*60 {
		scores := ScorePosts(posts)
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
		fmt.Println("Score snapshot complete")
	} else {
		fmt.Printf("Next score snapshot in %v seconds\n", 60*60-(now-lastUpdate))
	}

	fmt.Printf("Updated posts at %s\n", time.Now())
	fmt.Printf("Top post: %s\n", posts[0].Title)
	time.Sleep(5 * time.Minute)
	UpdatePosts(db)
}

func getHTML(url string) []byte {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
