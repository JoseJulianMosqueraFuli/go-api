package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Delivery struct {
	ID           string    `json:"id" gorm:"primary_key"`
	CreationDate time.Time `json:"creation_date"`
	State        string    `json:"state"`
	Pickup       struct {
		PickupLat float64 `json:"pickup_lat" gorm:"-"`
		PickupLon float64 `json:"pickup_lon" gorm:"-"`
	} `json:"pickup" gorm:"-"`
	Dropoff struct {
		DropoffLat float64 `json:"dropoff_lat" gorm:"-"`
		DropoffLon float64 `json:"dropoff_lon" gorm:"-"`
	} `json:"dropoff" gorm:"-"`
	ZoneID string `json:"zone_id"`
}

var db *gorm.DB

func main() {
	// Replace the database connection string with the correct values for your PostgreSQL setup
	dsn := "host=db user=user password=password dbname=db port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Automatically create the Delivery table if it doesn't exist
	err = db.AutoMigrate(&Delivery{})
	if err != nil {
		panic("failed to create table")
	}

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	router.POST("/deliveries", func(c *gin.Context) {
		var delivery Delivery
		if err := c.ShouldBindJSON(&delivery); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save the delivery to the database
		if err := db.Create(&delivery).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create delivery"})
			return
		}

		c.JSON(http.StatusCreated, delivery)
	})

	router.Run(":8080")
}
