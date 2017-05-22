package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func getPlanetaryApod(c *gin.Context) {
	ApodJSON := getApodJSON()
	c.JSON(http.StatusOK, ApodJSON)
}

func getApodJSON() Apod {

	ApodJSON := Apod{}

	//TODO Erase API KEY from repository history and move to env
	resp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=BvzqGsSJDYhfXLJ94uiaJDF7NLtrKJGdYW42eORT")
	if err != nil {
		println(err)
		return ApodJSON
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &ApodJSON)
	if err != nil {
		println(err)
		return ApodJSON
	}

	return ApodJSON
}
