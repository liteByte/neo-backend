package main

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"strconv"
	"time"
	"liteByte/neo-backend/validator"
)

type Config struct {
	Port string `default:"3001"`
	Nasa_key string `required:"true"`
}

type Apod struct {
	Url string `json:"url"`
}

type NeoFeed struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Size          map[string]float64 `json:"size"`
	Dangerous     bool `json:"dangerous"`
	Velocity      float64 `json:"velocity"`
	Miss_distance float64 `json:"miss_distance"`
}

func main() {

	var c Config
	c.Port = os.Getenv("PORT")
	c.Nasa_key = os.Getenv("NASA_KEY")
	validator.Validate(&c)

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/neo/feed", getNeoFeed)
	router.GET("/planetary/apod", getPlanetaryApod)

	router.Run(":" + c.Port)
}

func getNeoFeed(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	NeoFeedJSONArray := getNeoFeedJSON(startDate, endDate)
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, NeoFeedJSONArray)
}

func getPlanetaryApod(c *gin.Context) {
	ApodJSON := getApodJSON()
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, ApodJSON)
}

func getApodJSON() Apod {

	ApodJSON := Apod{}

	resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=BvzqGsSJDYhfXLJ94uiaJDF7NLtrKJGdYW42eORT")

	if err != nil {
		return ApodJSON
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &ApodJSON)

	if err != nil {
		return ApodJSON
	}

	return ApodJSON
}

func getNeoFeedJSON(startDate, endDate string) []NeoFeed {

	var neoFeedArray = []NeoFeed{}

	t := time.Now()
	currentDate := t.Format("2006-01-02")

	if startDate == "" {
		startDate = currentDate
		endDate = currentDate
	}

	if endDate == "" {
		endDate = currentDate
	}

	neoFeedURL := "https://api.nasa.gov/neo/rest/v1/feed?start_date=" + startDate +
		"&end_date=" + endDate + "&api_key=BvzqGsSJDYhfXLJ94uiaJDF7NLtrKJGdYW42eORT"

	//neoFeedURL := fmt.Sprintf("https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s",
	//startDate, endDate, )

	println(neoFeedURL)

	resp, err := http.Get(neoFeedURL)

	if err != nil {
		return neoFeedArray
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	f := feedParser{}

	err = json.Unmarshal(body, &f)

	neoFeedArray = parseNeoFeedJSON(f)

	return neoFeedArray
}

func parseNeoFeedJSON(f feedParser) []NeoFeed {

	var neoFeedArray = []NeoFeed{}

	for _, neoMapArray := range f.NearEarthObjects {
		for _, neo := range neoMapArray {
			neoObj := NeoFeed{}
			for k, v := range neo.(map[string]interface{}) {
				switch k {
				case "neo_reference_id":
					neoObj.Id = v.(string)
				case "name":
					neoObj.Name = v.(string)
				case "is_potentially_hazardous_asteroid":
					neoObj.Dangerous = v.(bool)
				case "estimated_diameter":
					size := v.(map[string]interface{})
					km := size["kilometers"].(map[string]interface{})
					neoObj.Size = map[string]float64{
						"min": km["estimated_diameter_min"].(float64),
						"max": km["estimated_diameter_max"].(float64),
						"avg": (km["estimated_diameter_max"].(float64) + km["estimated_diameter_min"].(float64)) / 2,
					}
				case "close_approach_data":
					approach_data := v.([]interface{})
					approach_data_v := approach_data[0].(map[string]interface{})

					velocity := approach_data_v["relative_velocity"].(map[string]interface{})
					miss_distance := approach_data_v["miss_distance"].(map[string]interface{})

					neoObj.Velocity, _ = strconv.ParseFloat((velocity["kilometers_per_hour"].(string)), 64)
					neoObj.Miss_distance, _ = strconv.ParseFloat(miss_distance["kilometers"].(string), 64)
				}
			}
			neoFeedArray = append(neoFeedArray, neoObj)
		}
	}

	return neoFeedArray
}

type feedParser struct {
	NearEarthObjects map[string][]interface{} `json:"near_earth_objects"`
}
