package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func getNeoFeed(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	NeoFeedJSONArray := getNeoFeedJSON(startDate, endDate)
	c.JSON(http.StatusOK, NeoFeedJSONArray)
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

	//TODO Erase API KEY from repository history and move to env
	neoFeedURL := "https://api.nasa.gov/neo/rest/v1/feed?start_date=" + startDate +
		"&end_date=" + endDate + "&api_key=BvzqGsSJDYhfXLJ94uiaJDF7NLtrKJGdYW42eORT"

	resp, err := http.Get(neoFeedURL)
	if err != nil {
		println(err)
		return neoFeedArray
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err)
		return neoFeedArray
	}

	parser := feedParser{}

	err = json.Unmarshal(body, &parser)
	if err != nil {
		println(err)
		return neoFeedArray
	}

	neoFeedArray = parseNeoFeedJSON(parser)

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
					neoObj.IsDangerous = v.(bool)
				case "estimated_diameter":
					size := v.(map[string]interface{})
					m := size["meters"].(map[string]interface{})
					neoObj.Size = map[string]float64{
						"min": m["estimated_diameter_min"].(float64),
						"max": m["estimated_diameter_max"].(float64),
						"avg": (m["estimated_diameter_max"].(float64) + m["estimated_diameter_min"].(float64)) / 2,
					}
				case "close_approach_data":
					approach_data := v.([]interface{})
					approach_data_v := approach_data[0].(map[string]interface{})

					velocity := approach_data_v["relative_velocity"].(map[string]interface{})
					missDistance := approach_data_v["miss_distance"].(map[string]interface{})

					neoObj.Velocity, _ = strconv.ParseFloat(velocity["kilometers_per_hour"].(string), 64)
					neoObj.MissDistance, _ = strconv.ParseFloat(missDistance["kilometers"].(string), 64)
				}
			}
			neoFeedArray = append(neoFeedArray, neoObj)
		}
	}

	return neoFeedArray
}
