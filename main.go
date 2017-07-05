package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"database/sql"
	//"os"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html", "templates/foot.html", "templates/head.html", "templates/food_cell.html")
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/img", "./img")

	// Open Connection
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router.GET("/", func(c *gin.Context) {
		foodMap := GetDaily(db)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"foodData": foodMap,
		})
	})

	router.GET("/getReviews/:id", func(c *gin.Context) {
		id := c.Param("id")
		reviews := GetReviews(db, id)
		c.JSON(http.StatusOK, reviews)
	})

	router.Run(":" + port)
}
