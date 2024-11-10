package main

import (
    "leetcode_backend/config"
    "leetcode_backend/controllers"
    "github.com/gin-gonic/gin"
)

func main() {
    config.ConnectDatabase()
    r := gin.Default()
    controllers.RegisterRoutes(r)
    r.Run()
}
