package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gin-gonic/gin"
    "encoding/json"
    // "fmt"
    "io/ioutil"
    "strconv"
    "time"
)

func main() {

    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.GET("/neo/feed", func(c *gin.Context) {
        startDate := c.Query("start_date")
        endDate := c.Query("end_date")
        NeoFeedJSONArray := getNeoFeedJSON(startDate, endDate)
        c.JSON(http.StatusOK, NeoFeedJSONArray)
    })

    router.GET("/planetary/apod", func(c *gin.Context) {
        ApodJSON := getApodJSON()
        c.JSON(http.StatusOK, ApodJSON)
    })

    router.Run(":" + port)
}

type Apod struct {
    Url string `json:"url"`
}

type NeoFeed struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Size map[string]float64 `json:"size"`
    Dangerous bool `json:"dangerous"`
    Velocity float64 `json:"velocity"`
    Miss_distance float64 `json:"miss_distance"`
}

func getApodJSON() Apod {

    ApodJSON := Apod{}

    resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=BvzqGsSJDYhfXLJ94uiaJDF7NLtrKJGdYW42eORT")

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

func getNeoFeedJSON(startDate, endDate string) []NeoFeed{

    var neoFeedArray = []NeoFeed{}

    // Get current date
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

    println(neoFeedURL)

    resp, err := http.Get(neoFeedURL)

    if (err != nil){
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
    NearEarthObjects map[string] []interface{} `json:"near_earth_objects"`
}