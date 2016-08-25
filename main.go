package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "encoding/json"
    "fmt"
    "io/ioutil"
)

func main() {

    fmt.Print("Hi!")

    port := os.Getenv("PORT")
    // port = "5000"

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.GET("/neo/feed", func(c *gin.Context) {
        c.String(http.StatusOK, "This will return Neo Feed")
    })

    router.GET("/planetary/apod", func(c *gin.Context) {
        ApodJSON := getApodJSON()
        c.JSON(http.StatusOK, ApodJSON)
    })

    // ApodJSON := getApodJSON()
    // fmt.Printf("JSON body: %v \n JSON type: %t \n", ApodJSON, ApodJSON)

    router.Run(":" + port)
}

type Apod struct {
    Url string
}

func getApodJSON() Apod{

    ApodJSON := Apod{}

    resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY")

    if (err != nil){
        return ApodJSON
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    err = json.Unmarshal(body, &ApodJSON)

    if (err != nil){
        return ApodJSON
    }

    return ApodJSON
}