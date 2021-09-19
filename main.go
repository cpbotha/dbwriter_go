package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // _ is for import side-effect, no explicit use
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/cpbotha/dbwriter_go/docs"
)

// GORM model
// this would usually go in Models/sample.go
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
	TimeStamp time.Time `json:"timestamp" binding:"required" example:"2021-09-19T10:41:33.333Z"`
	V0        float64   `json:"v0" binding:"-"`
	V1        float64   `json:"v1" binding:"-"`
}

func APIRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "Hello, World!"})
}

// from https://github.com/swaggo/swag/blob/17c1766b6349df2ab31d52e87be0ae3abca0d239/example/celler/httputil/error.go
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// CreateSample writes a time sample to DB
// @Summary Write single sample
// @Param sample body CreateSampleInput true "create sample"
// @Success 200 {object} Sample
// @Failure 400 {object} HTTPError
// @Router /samples/ [post]
func CreateSample(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// Validate input
	var input CreateSampleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// instead of c.JSON(400, gin.H{"error": err.Error()}) we use func + struct from swaggo
		NewError(c, http.StatusBadRequest, err)
		return
	}
	// Create Sample
	// what should we do if v0 or v1 is not supplied?
	sample := Sample{Name: input.Name, TimeStamp: input.TimeStamp, v0: input.V0, v1: input.V1}
	db.Create(&sample)
	c.JSON(http.StatusOK, gin.H{"data": sample})
}

// GetSample retrieves single sample from database by ID
// @Summary Retrieve sample by ID
// @Param id path int true "Sample ID"
// @Success 200 {object} Sample
// @Failure 404 {object} HTTPError
// @Router /samples/{id} [get]
func GetSample(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var sample Sample
	if err := db.Where("id = ?", c.Param("id")).First(&sample).Error; err != nil {
		NewError(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sample})
}

// ListSamples godoc
// @Summary List all samples
// @Produce json
// @Success 200 {object} []Sample
// @Router /samples [get]
func ListSamples(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var samples []Sample
	db.Find(&samples)
	c.JSON(http.StatusOK, gin.H{"data": samples})
}

// @title DBWriter API in Go with Gin and GORM
// @version 1.0
// @description Bleh bleh bleh
// @BasePath /
// @schemes http
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

	r.GET("/", APIRoot)
	r.POST("/samples", CreateSample)
	r.GET("/samples/:id", GetSample)

	r.GET("/samples", ListSamples)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Started up!")

	r.Run("0.0.0.0:8080")
}
