package main

import (
  "fmt"
  "net/http"

  "github.com/ianjdarrow/slow-hn/router"
  "github.com/ianjdarrow/slow-hn/util"
)

func main() {
  server := router.InitRouter()
  util.UpdatePosts(50)
  fmt.Println("Listening on 4000!")
  http.ListenAndServe(":4000", server)
}
