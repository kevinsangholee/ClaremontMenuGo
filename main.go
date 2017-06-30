package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	//"os"
)

func main() {

	//port := os.Getenv("PORT")
	//
	//if port == "" {
	//	log.Fatal("$PORT must be set")
	//}

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html", "templates/foot.html", "templates/head.html", "templates/food_cell.html")
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/img", "./img")

	foodMap := GetDaily()

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"foodData": foodMap,
		})
	})

	router.Run(":8080")
}
