package menudb

import (
	"log"
	"strings"
	"database/sql"
	"strconv"
	"time"
)

type FoodItem struct {
	Id   		 int
	Name 		 string
	School		 int
	Image     	 string
	Review_count int
	Rating       float64
	Daily		 string
}

type ReviewItem struct {
	Food_id		int	   `json:"food_id"`
	User_id		string `json:"user_id"`
	Rating 	 	int    `json:"rating"`
	Review_text string `json:"review_text"`
	Created_at  string `json:"created_at"`
}

const (
	DB_HOST = "tcp(jj820qt5lpu6krut.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)"
	DB_USER = "vc568j0frxncao3a"
	DB_PASS = "yq4trluh9rq7tvkt"
	DB_NAME = "pzw6eam7d8a9orvg"
)

/* This function:
   1. Queries the foods database to check for any daily food
   2. Parses the rows returned into maps separated by schools and meals and returns the overarching map
 */
func GetDaily(db *sql.DB) map[string][]*FoodItem {

	// Initialize
	foodMap := make(map[string][]*FoodItem)


	for i := 0; i < 7; i++ {
		if (!IsWeekend()) {
			for j := 0; j < 3; j++ {
				foodMap[strconv.Itoa(i)+strconv.Itoa(j)] = make([]*FoodItem, 0)
			}
		} else {
			for j := 2; j <= 3; j++ {
				foodMap[strconv.Itoa(i)+strconv.Itoa(j)] = make([]*FoodItem, 0)
			}
		}
	}

	// Querying daily
	rows, err := db.Query("SELECT id, name, school, image, review_count, rating, daily FROM foods WHERE daily <> '.'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parsing rows by school and meal
	for rows.Next() {
		currFood := new(FoodItem)
		err := rows.Scan(&currFood.Id, &currFood.Name, &currFood.School, &currFood.Image,
			&currFood.Review_count, &currFood.Rating, &currFood.Daily)
		if err != nil {
			log.Fatal(err)
		}
		mealSlice := strings.Split(currFood.Daily[1:], "")
		for _ , val := range mealSlice {
			foodMap[strconv.Itoa(currFood.School) + val] = append(foodMap[strconv.Itoa(currFood.School) + val], currFood)
		}
	}

	return foodMap
}

func GetReviews(db *sql.DB, foodID string) []*ReviewItem {

	rows, err := db.Query("SELECT food_id, user_id, rating, review_text, created_at FROM reviews WHERE food_id = " +
		foodID + " AND NOT (review_text = '')")
	if err != nil {
		log.Fatal(err)
	}

	reviews := make([]*ReviewItem, 0)
	for rows.Next() {
		review := new(ReviewItem)
		err := rows.Scan(&review.Food_id, &review.User_id, &review.Rating, &review.Review_text, &review.Created_at)
		if err != nil {
			log.Fatal(err)
		}
		reviews = append(reviews, review)
	}

	return reviews
}

func GetMeal(db *sql.DB, school string, meal string) []*FoodItem {
	foods := make([]*FoodItem, 0)

	query := "SELECT id, name, school, image, review_count, rating FROM foods WHERE school = " +
		school + " AND daily LIKE '%" + meal + "%'"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		foodItem := new(FoodItem)
		err := rows.Scan(&foodItem.Id, &foodItem.Name, &foodItem.School, &foodItem.Image, &foodItem.Review_count, &foodItem.Rating)
		if err != nil {
			log.Fatal(err)
		}
		foods = append(foods, foodItem)
	}

	return foods

}

func IsWeekend() bool {
	t := time.Now().Weekday()
	return t == 0 || t == 6
}