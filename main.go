package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	_ "github.com/cpbotha/dbwriter_go/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GORM model
// this would usually go in Models/sample.go
// https://gorm.io/docs/models.html
// getting the optional v0, v1 types right took up more time than anything else
// (also lost time because lowercase v* were ignored)
// sql.NullFloat64 would be more explicit than *float64, but complicates (de)serialization
type Sample struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	TimeStamp time.Time `json:"timestamp"`
	V0        *float64  `json:"v0,omitempty"`
	V1        *float64  `json:"v1,omitempty"`
}

// post data is validated into this struct, from where we can populate the GORM model
// I would like omission OR null to result in NULL for that sensor value in the database
// https://www.calhoun.io/how-to-determine-if-a-json-key-has-been-set-to-null-or-not-provided/
type CreateSampleInput struct {
	Name      string    `json:"name" binding:"required"`
	TimeStamp time.Time `json:"timestamp" binding:"required" example:"2021-09-19T10:41:33.333Z"`
	V0        *float64  `json:"v0" binding:"-"`
	V1        *float64  `json:"v1" binding:"-"`
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
// @Summary Add a single sensor sample to the database
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

	sample := Sample{Name: input.Name, TimeStamp: input.TimeStamp, V0: input.V0, V1: input.V1}
	db.Create(&sample)
	c.JSON(http.StatusOK, sample)
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
	c.JSON(http.StatusOK, sample)
}

// List all samples in the database.
//
// @Summary List all samples
// @Produce json
// @Success 200 {object} []Sample
// @Router /samples [get]
func ListSamples(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var samples []Sample
	db.Find(&samples)
	c.JSON(http.StatusOK, samples)
}

// @title DBWriter API in Go with Gin and GORM
// @version 1.0
// @description Bleh bleh bleh
// @BasePath /
// @schemes http
func main() {
	//db, err := gorm.Open(sqlite.Open("bleh.db"), &gorm.Config{})
	db, err := gorm.Open(postgres.Open("host=localhost user=dbwriter password=blehbleh dbname=dbwriter_go port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Sample{})

	// you can use something like this to add a row to the database:
	// v := 0.0
	// db.Create(&Sample{0, "bleh blah bleh", time.Now(), nil, &v})

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

	fmt.Printf("NumCPU: %d and GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(-1))
	fmt.Println("Started up!")

	r.Run("0.0.0.0:8080")
}
