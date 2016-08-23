package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.GET("/neo/feed", func(c *gin.Context) {
        c.String(http.StatusOK, "This wiil return Neo feed")
    })

    router.GET("/planetary/apod", func(c *gin.Context) {
        c.String(http.StatusOK, "This wiil return the APOD")
    })

    router.Run(":" + port)
}
