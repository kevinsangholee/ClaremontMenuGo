package main

import (
	"log"
	"strings"
	"database/sql"
	"strconv"
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
func GetDaily() map[string][]*FoodItem {

	// Open Connection
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize
	foodMap := make(map[string][]*FoodItem)

	for i := 0; i < 7; i++ {
		for j := 0; j < 3; j++ {
			foodMap[strconv.Itoa(i) + strconv.Itoa(j)] = make([]*FoodItem, 0)
		}
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
		err := rows.Scan(&currFood.Id, &currFood.Name, &currFood.School, &currFood.Image,
			&currFood.Review_count, &currFood.Rating, &currFood.Daily)
		if err != nil {
			log.Fatal(err)
		}
		mealSlice := strings.Split(currFood.Daily, "")
		for _ , val := range mealSlice {
			foodMap[strconv.Itoa(currFood.School) + val] = append(foodMap[strconv.Itoa(currFood.School) + val], currFood)
		}
	}

	return foodMap
}