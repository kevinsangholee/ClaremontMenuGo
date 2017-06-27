package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	DB_HOST = "jj820qt5lpu6krut.cbetxkdyhwsb.us-east-1.rds.amazonaws.com"
	DB_USER = "vc568j0frxncao3a"
	DB_PASS = "yq4trluh9rq7tvkt"
	DB_NAME = "pzw6eam7d8a9orvg"
)

func main() {
	var review_text string

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html", "templates/foot.html", "templates/head.html", "templates/food_cell.html")
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/img", "./img")

	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to database successfully!")
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM reviews")
	
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&review_text)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(review_text)
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.Run(":8080")
}
