package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/darolpz/golang-doc-generator/routes"
	"github.com/darolpz/golang-doc-generator/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	port := utils.GetEnvVariable("PORT")
	t := time.Now()
	f, _ := os.Create(fmt.Sprintf("logs/gin-%s.log", t.Format("2006-01-02")))
	gin.DefaultWriter = io.MultiWriter(f)
	router := routes.InitRoutes()
	router.Static("/docs", "./docs")
	router.Run(fmt.Sprintf(":%v", port))
}
