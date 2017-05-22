package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liteByte/structval"
	"os"
)

type Config struct {
	Port     string `default:"3002"`
	Nasa_key string `required:"true"`
}

type Apod struct {
	Url string `json:"url"`
}

type NeoFeed struct {
	Id           string             `json:"id"`
	Name         string             `json:"name"`
	Size         map[string]float64 `json:"size"`
	IsDangerous  bool               `json:"isDangerous"`
	Velocity     float64            `json:"velocity"`
	MissDistance float64            `json:"missDistance"`
}

func main() {

	var c Config
	c.Port = os.Getenv("PORT")
	c.Nasa_key = os.Getenv("NASA_KEY")
	structval.Validate(&c)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(Cors())

	router.GET("/neo/feed", getNeoFeed)
	router.GET("/planetary/apod", getPlanetaryApod)

	router.Run(":" + c.Port)
}

type feedParser struct {
	NearEarthObjects map[string][]interface{} `json:"near_earth_objects"`
}
