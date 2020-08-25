package routes

import (
	"log"
	"net/http"

	"github.com/darolpz/golang-doc-generator/models"
	"github.com/darolpz/golang-doc-generator/services"
	"github.com/darolpz/golang-doc-generator/utils"
	"github.com/gin-gonic/gin"
)

//InitRoutes init routes
func InitRoutes() *gin.Engine {
	router := gin.New()
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(utils.LogFormat))
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	v2 := router.Group("/v2")

	v1.POST("/json", func(c *gin.Context) {
		var params models.Parameter
		if err := c.BindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if params.To == "" {
			params.To = "develop"
		}

		commits := services.GetCommits(&params)
		c.JSON(200, gin.H{
			"commits": commits,
		})
	})

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1.POST("/pdf", func(c *gin.Context) {
		var params models.Parameter

		if err := c.BindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if params.To == "" {
			params.To = "develop"
		}

		commits := services.GetCommits(&params)
		err := services.GeneratePdf(&params, &commits)
		if err != nil {
			log.Printf("Error %s\n", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			panic(err)
		}
		c.JSON(200, gin.H{
			"message": "Release notes has been created succesfully",
		})
	})

	v2.POST("/pdf", func(c *gin.Context) {

		var params models.Parameter
		if err := c.BindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if params.To == "" {
			params.To = "develop"
		}

		commits := services.GetCommits(&params)
		filename, err := services.GeneratePdf2(&params, &commits)
		if err != nil {
			log.Printf("Error %s\n", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			panic(err)
		}
		services.NotifyThroughtSlack(filename, params.Channel)
		c.JSON(200, gin.H{
			"message": "Release notes has been created succesfully",
		})
	})
	return router
}
