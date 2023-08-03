package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	ID           string    `json:"id"`
	CreationDate time.Time `json:"creation_date"`
	State        string    `json:"state"`
	Pickup       struct {
		PickupLat float64 `json:"pickup_lat"`
		PickupLon float64 `json:"pickup_lon"`
	} `json:"pickup"`
	Dropoff struct {
		DropoffLat float64 `json:"dropoff_lat"`
		DropoffLon float64 `json:"dropoff_lon"`
	} `json:"dropoff"`
	ZoneID string `json:"zone_id"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	r.Run(":8080")
}
