package util

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "math"
  "net/http"
  "strconv"
  "time"

  "github.com/ianjdarrow/slow-hn/controllers"
  "github.com/ianjdarrow/slow-hn/models"
)

var (
  topPostsUrl string = "https://hacker-news.firebaseio.com/v0/topstories.json"
)

func FetchTopPosts(numPosts int) []models.Post {
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
  now := time.Now()
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

func UpdatePosts(numPosts int) {
  posts := FetchTopPosts(numPosts)
  for _, post := range posts {
    if _, ok := controllers.AllPosts[post.ID]; !ok {
      controllers.AllPosts[post.ID] = post
    }
  }
  fmt.Printf("Updated posts at %s, total count: %v\n", time.Now(), len(controllers.AllPosts))
  time.Sleep(time.Hour)
  UpdatePosts(numPosts)
}

func getHTML(url string) []byte {
  resp, _ := http.Get(url)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  return body
}
