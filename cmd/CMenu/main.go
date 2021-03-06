package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"database/sql"
	"github.com/kevinsangholee/ClaremontMenuGo"
	"strconv"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html", "templates/foot.html", "templates/head.html", "templates/food_cell.html",
						       "templates/index_weekend.html")
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/img", "./img")

	// Open Connection
	dsn := menudb.DB_USER + ":" + menudb.DB_PASS + "@" + menudb.DB_HOST + "/" + menudb.DB_NAME
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if(!menudb.IsWeekend()) {
		router.GET("/", func(c *gin.Context) {
			foodMap := menudb.GetDaily(db)
			c.HTML(http.StatusOK, "index.html", gin.H{
				"foodData": foodMap,
			})
		})
	} else {
		router.GET("/", func(c *gin.Context) {
			foodMap := menudb.GetDaily(db)
			c.HTML(http.StatusOK, "index_weekend.html", gin.H{
				"foodData": foodMap,
			})
		})
	}

	router.GET("/getReviews/:id", func(c *gin.Context) {
		id := c.Param("id")
		reviews := menudb.GetReviews(db, id)
		c.JSON(http.StatusOK, reviews)
	})

	router.GET("/getMeal", func(c *gin.Context) {
		school := c.Query("school")
		meal := c.Query("meal")

		foods := menudb.GetMeal(db, school, meal)
		c.JSON(http.StatusOK, foods)
	})

	router.GET("/getSingleFood/:id", func(c *gin.Context) {
		id := c.Param("id")
		food := menudb.GetSingleFood(db, id)
		c.JSON(http.StatusOK, food)
	})

	router.POST("/addReview", func(c *gin.Context) {
		food_id := c.PostForm("food_id")
		user_id := c.PostForm("user_id")
		rating := c.PostForm("rating")
		review_text := c.PostForm("review_text")
		created_at := c.PostForm("created_at")

		review_id := menudb.AddReview(db, food_id, user_id, rating, review_text, created_at)

		c.String(http.StatusOK, strconv.Itoa(int(review_id)))
	})

	router.POST("/deleteReview", func(c *gin.Context) {
		review_id := c.PostForm("review_id")
		food_id := c.PostForm("food_id")

		menudb.DeleteReview(db, review_id, food_id)

		c.String(http.StatusOK, "Review deleted!")
	})

	router.POST("/updateReview", func(c *gin.Context) {
		review_id := c.PostForm("review_id")
		rating := c.PostForm("rating")
		review_text := c.PostForm("review_text")
		created_at := c.PostForm("created_at")

		menudb.UpdateReview(db, review_id, rating, review_text, created_at)

		c.String(http.StatusOK, "Review updated!")
	})

	//router.Run(":" + port)
	router.Run(":" + port)
}
