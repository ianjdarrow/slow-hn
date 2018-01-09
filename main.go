package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/ianjdarrow/slow-hn/controllers"
	"github.com/ianjdarrow/slow-hn/db"
	"github.com/ianjdarrow/slow-hn/router"
	"github.com/ianjdarrow/slow-hn/util"
)

var database *sqlx.DB

func main() {
	database = db.InitDB()
	controllers.DB = database
	server := router.InitRouter()
	go util.UpdatePosts(database)
	fmt.Println("Server listening on port 4000!")
	log.Fatal(http.ListenAndServe(":4000", server))
}
