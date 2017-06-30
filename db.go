package main

import (
	"log"
	"strings"
	"database/sql"
	"strconv"
)

type FoodItem struct {
	id   		 int
	name 		 string
	school		 int
	image     	 string
	review_count int
	rating       float64
	daily		 string
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
func GetDaily() map[int]map[int][]*FoodItem {

	// Open Connection
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize
	foodMap := make(map[int](map[int][]*FoodItem))
	for i := 0; i < 7; i++ {
		foodMap[i] = make(map[int][]*FoodItem)
	}

	// Querying daily
	rows, err := db.Query("SELECT id, name, school, image, review_count, rating, daily FROM foods WHERE daily <> -1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Parsing rows by school and meal
	for rows.Next() {
		currFood := new(FoodItem)
		err := rows.Scan(&currFood.id, &currFood.name, &currFood.school, &currFood.image,
			&currFood.review_count, &currFood.rating, &currFood.daily)
		if err != nil {
			log.Fatal(err)
		}
		mealSlice := strings.Split(currFood.daily, "")
		for _ , val := range mealSlice {
			mealInt, err := strconv.Atoi(val)
			if err != nil {
				log.Fatal(err)
			}
			foodMap[currFood.school][mealInt] = append(foodMap[currFood.school][mealInt], currFood)
		}
	}

	return foodMap
}