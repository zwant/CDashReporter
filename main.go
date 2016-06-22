package main

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/kataras/iris"
)

const (
	influxURL = "http://localhost:8086"
	influxDb  = "svante"
)

type Build struct {
	Success   bool
	RuntimeMs int64
	Agent     string
}

func HandleBuildNotification(c *iris.Context) {
	buildData := Build{}
	if err := c.ReadJSON(&buildData); err != nil {
		panic(err.Error())
	}
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxURL,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxDb,
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}
	// Create a point and add to batch
	tags := map[string]string{
		"agent": buildData.Agent,
	}

	fields := map[string]interface{}{
		"runtime_ms": buildData.RuntimeMs,
		"success":    buildData.Success,
	}
	pt, err := client.NewPoint("builds", tags, fields, time.Now())

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	cli.Write(bp)

	c.JSON(iris.StatusOK, iris.Map{
		"status": "ok",
	})
}

func main() {
	// Make client

	// render JSON
	iris.Post("/build", HandleBuildNotification)

	iris.Listen(":8080")
}
