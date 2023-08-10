package main

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Delivery struct {
	ID           string    `json:"id" gorm:"primary_key"`
	CreationDate time.Time `json:"creation_date"`
	State        string    `json:"state"`
	Pickup       struct {
		PickupLat float64 `json:"pickup_lat"`
		PickupLon float64 `json:"pickup_lon"`
	} `json:"pickup" gorm:"-"`
	Dropoff struct {
		DropoffLat float64 `json:"dropoff_lat"`
		DropoffLon float64 `json:"dropoff_lon"`
	} `json:"dropoff" gorm:"-"`
	ZoneID        string `json:"zone_id"`
	CreatorID     string `json:"creator_id"`
	AssignedBotID string `json:"assigned_bot_id"`
}

type Bot struct {
	ID       string `json:"id" gorm:"primary_key"`
	Status   string `json:"status"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location" json:"location" gorm:"-"`
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

	// Automatically create the Delivery and Bot table if it doesn't exist
	err = db.AutoMigrate(&Delivery{}, &Bot{})
	if err != nil {
		panic("failed to create tables")
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

		// Check latitude and longitude restrictions
		if !isValidLatitude(delivery.Pickup.PickupLat) || !isValidLongitude(delivery.Pickup.PickupLon) ||
			!isValidLatitude(delivery.Dropoff.DropoffLat) || !isValidLongitude(delivery.Dropoff.DropoffLon) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude or longitude"})
			return
		}

		// Generate ID and set CreationDate
		delivery.ID = generateID()
		delivery.CreationDate = time.Now()

		// Save the delivery to the database
		if err := db.Create(&delivery).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create delivery"})
			return
		}

		c.JSON(http.StatusCreated, delivery)
	})

	router.GET("/deliveries/:id", func(c *gin.Context) {
		id := c.Param("id")

		var delivery Delivery
		if err := db.Where("id = ?", id).First(&delivery).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
			return
		}

		c.JSON(http.StatusOK, delivery)
	})

	router.GET("/deliveries/creator/:creatorID", func(c *gin.Context) {
		creatorID := c.Param("creatorID")

		var deliveries []Delivery
		if err := db.Where("creator_id = ?", creatorID).Find(&deliveries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deliveries"})
			return
		}

		c.JSON(http.StatusOK, deliveries)
	})

	router.GET("/deliveries", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))

		if page < 1 {
			page = 1
		}
		if perPage < 1 {
			perPage = 10
		}

		var deliveries []Delivery
		offset := (page - 1) * perPage

		if err := db.Offset(offset).Limit(perPage).Find(&deliveries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deliveries"})
			return
		}

		total := getTotalDeliveries()

		c.JSON(http.StatusOK, gin.H{
			"page":       page,
			"perPage":    perPage,
			"total":      total,
			"deliveries": deliveries,
		})
	})

	router.GET("/deliveries/by-date", func(c *gin.Context) {
		dateFilter := c.DefaultQuery("date", "")
		layout := "2006-01-02"

		if dateFilter == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
			return
		}

		filterDate, err := time.Parse(layout, dateFilter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use yyyy-mm-dd"})
			return
		}

		var deliveries []Delivery

		// Filter deliveries by date
		if err := db.Where("DATE(creation_date) = ?", filterDate).Find(&deliveries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch deliveries"})
			return
		}

		c.JSON(http.StatusOK, deliveries)
	})

	router.POST("/bots", func(c *gin.Context) {
		var bot Bot
		if err := c.ShouldBindJSON(&bot); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check latitude and longitude restrictions
		if !isValidLatitude(bot.Location.Lat) || !isValidLongitude(bot.Location.Lon) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude or longitude"})
			return
		}

		// Generate ID for the bot
		bot.ID = generateID()

		// Save the bot to the database
		if err := db.Create(&bot).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bot"})
			return
		}

		c.JSON(http.StatusCreated, bot)
	})

	router.GET("/bots/by-zone/:zoneID", func(c *gin.Context) {
		zoneID := c.Param("zoneID")

		var bots []Bot
		if err := db.Where("zone_id = ?", zoneID).Find(&bots).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bots by zone"})
			return
		}

		c.JSON(http.StatusOK, bots)
	})

	router.PUT("/deliveries/assign-bot/:id", func(c *gin.Context) {
		deliveryID := c.Param("id")
		var delivery Delivery

		if err := db.Where("id = ?", deliveryID).First(&delivery).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Delivery not found"})
			return
		}

		// Check if the delivery is already assigned
		if delivery.AssignedBotID != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Delivery is already assigned to a bot"})
			return
		}

		// Find the nearest available bot in the same zone as the delivery
		nearestBot, err := findNearestAvailableBot(delivery.Pickup.PickupLat, delivery.Pickup.PickupLon, delivery.ZoneID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find available bot"})
			return
		}

		// Assign the bot to the delivery
		delivery.AssignedBotID = nearestBot.ID
		delivery.State = "assigned"

		if err := db.Save(&delivery).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign bot to delivery"})
			return
		}

		// Update the bot's status to busy
		nearestBot.Status = "busy"
		if err := db.Save(&nearestBot).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bot status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Bot assigned successfully"})
	})

	router.Run(":8080")
}

func generateID() string {
	return uuid.New().String()
}

func isValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func isValidLongitude(lon float64) bool {
	return lon >= -180 && lon <= 180
}

func getTotalDeliveries() int64 {
	var total int64
	if err := db.Model(&Delivery{}).Count(&total).Error; err != nil {
		return 0
	}
	return total
}

// Function to Calculate the distance between two positions, to create distance.
const earthRadius = 6371.0 // Earth's radius in kilometers

// CalculateHaversineDistance calculates the distance between two points on Earth using the Haversine formula.
func CalculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Differences in coordinates
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

func findNearestAvailableBot(pickupLat, pickupLon float64, zoneID string) (Bot, error) {
	var availableBots []Bot
	if err := db.Where("status = 'available' AND zone_id = ?", zoneID).Find(&availableBots).Error; err != nil {
		return Bot{}, err
	}

	var nearestBot Bot
	minDistance := math.MaxFloat64

	for _, bot := range availableBots {
		distance := CalculateHaversineDistance(pickupLat, pickupLon, bot.Location.Lat, bot.Location.Lon)
		if distance < minDistance {
			nearestBot = bot
			minDistance = distance
		}
	}

	return nearestBot, nil
}
