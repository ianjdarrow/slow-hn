package router

import (
  "github.com/julienschmidt/httprouter"

  "github.com/ianjdarrow/slow-hn/controllers"
)

func InitRouter() *httprouter.Router {
  router := httprouter.New()
  router.GET("/", controllers.GetIndex)

  return router
}
