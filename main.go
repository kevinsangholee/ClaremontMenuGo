package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//const (
//	DB_HOST = "tcp(claremontmenu.com:3306)"
//	DB_USER = "claremo7_klee"
//	DB_PASS = "Miyukguk369"
//	DB_NAME = "claremo7_claremontmenu"
//)

func main() {
	var review_text string

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html", "templates/foot.html", "templates/head.html", "templates/food_cell.html")
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/img", "./img")

	//dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME
	db, err := sql.Open("mysql", "vc568j0frxncao3a:yq4trluh9rq7tvkt@tcp(jj820qt5lpu6krut.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/pzw6eam7d8a9orvg")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to database successfully!")
	}
	defer db.Close()
	row := db.QueryRow("SELECT review_text FROM reviews WHERE id=319")

	if err := row.Scan(&review_text); err != nil {
		log.Fatal(err)
	}

	log.Println(review_text)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.Run(":8080")
}
