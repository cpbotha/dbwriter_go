// go get -u github.com/jinzhu/gorm
// go get github.com/jinzhu/gorm/dialects/sqlite
// go get github.com/gin-gonic/gin

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // _ is for import side-effect, no explicit use
)

// GORM model
// this could also go in Models/sample.go
// https://gorm.io/docs/models.html
type Sample struct {
	ID        uint
	Name      string
	TimeStamp time.Time
	v0        float64
	v1        float64
}

// post data is validated into this struct, from where we can populate the GORM model
type CreateSampleInput struct {
	Name      string    `json:"name" binding:"required"`
	TimeStamp time.Time `json:"timestamp" binding:"required"`
	v0        float64   `json:"v0" binding:"-"`
	v1        float64   `json:"v0" binding:"-"`
}

func CreateSample(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Validate input
	var input CreateSampleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create Sample
	// what should we do if v0 or v1 is not supplied?
	sample := Sample{Name: input.Name, TimeStamp: input.TimeStamp, v0: input.v0, v1: input.v1}
	db.Create(&sample)
	c.JSON(http.StatusOK, gin.H{"data": sample})
}

func GetSample(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var sample Sample
	if err := db.Where("id = ?", c.Param("id")).First(&sample).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sample})
}

func ListSamples(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var samples []Sample
	db.Find(&samples)
	c.JSON(http.StatusOK, gin.H{"data": samples})
}

func main() {
	db, err := gorm.Open("sqlite3", "bleh.db")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.AutoMigrate(&Sample{})

	r := gin.Default()

	// inject db into context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/samples", CreateSample)
	r.GET("/samples/:id", GetSample)
	r.GET("/samples", ListSamples)

	fmt.Println("Started up!")

	r.Run()
}
