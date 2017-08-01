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

/*
	This function gets all review items of a food given its id.
 */
func GetReviews(db *sql.DB, foodID string) []*ReviewItem {

	rows, err := db.Query("SELECT food_id, user_id, rating, review_text, created_at FROM reviews WHERE food_id = " +
		foodID + " AND NOT (review_text = '')")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

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

/*
	This function is more for the mobile side, gets the array of food items for a meal given the school and meal number
 */
func GetMeal(db *sql.DB, school string, meal string) []*FoodItem {
	foods := make([]*FoodItem, 0)

	query := "SELECT id, name, school, image, review_count, rating FROM foods WHERE school = " +
		school + " AND daily LIKE '%" + meal + "%'"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

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

/*
	This is for the mobile app to update the review page once a review has been added
 */
func GetSingleFood(db *sql.DB, id string) *FoodItem {
	food := new(FoodItem)
	row := db.QueryRow("SELECT id, name, school, image, review_count, rating FROM foods WHERE id = " + id)
	err := row.Scan(&food.Id, &food.Name, &food.School, &food.Image, &food.Review_count, &food.Rating)
	if err != nil {
		log.Fatal(err)
	}
	return food
}

/*
	This is for the mobile app to actually add a review
 */
func AddReview(db *sql.DB, food_id string, user_id string, rating string, review_text string, created_at string) int64 {

	// Adding the review itself
	stmt, err := db.Prepare("INSERT INTO reviews (food_id, user_id, rating, review_text, created_at) " +
		"VALUES (?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec(food_id, user_id, rating, review_text, created_at)
	if err != nil {
		log.Fatal(err)
	}
	// Keeping the id of the review just created to send back to the phone
	review_id, _ := result.LastInsertId()

	// First get current review count and total score to increment both inside this function
	row := db.QueryRow("SELECT review_count, total_score FROM foods WHERE id = " + food_id)
	var review_count int
	var total_score int
	row.Scan(&review_count, &total_score)

	// Increment review count and total score
	review_count++
	parsedRating, _ := strconv.ParseInt(rating, 10, 64)
	total_score += int(parsedRating)

	// Calculate new average
	new_average := strconv.FormatFloat(float64(total_score) / float64(review_count), 'E', 2, 64)

	// Updating foods database to get the correct review count and total score
	stmt, err = db.Prepare("UPDATE foods SET review_count = ?, total_score = ?, rating = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(review_count, total_score, new_average, food_id)
	if err != nil {
		log.Fatal(err)
	}

	return review_id
}

/*
	This function is for the phone to delete a review and update all averages accordingly
 */
func DeleteReview(db *sql.DB, review_id string, food_id string) {
	// First get the actual rating from the review given the review_id
	row := db.QueryRow("SELECT rating FROM reviews WHERE id = " + review_id)
	var rating int
	row.Scan(&rating)

	// Get the current review count and total score from the foods table
	row = db.QueryRow("SELECT review_count, total_score FROM foods WHERE id = " + food_id)
	var review_count int
	var total_score int
	row.Scan(&review_count, &total_score)

	// Decrement review count and subtract total_score from rating, and calculate the average again
	review_count--
	total_score -= rating
	new_average := strconv.FormatFloat(float64(total_score) / float64(review_count), 'E', 2, 64)

	// Insert the new values back into foods table
	stmt, err := db.Prepare("UPDATE foods SET review_count = ?, total_score = ?, rating = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(review_count, total_score, new_average, food_id)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, delete the review from reviews table
	stmt, err = db.Prepare("DELETE FROM reviews WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(review_id)
	if err != nil {
		log.Fatal(err)
	}
}

/*
	This function takes care of updating a review and updating the average count and what not
 */
func UpdateReview(db *sql.DB, review_id string, rating string, review_text string, created_at string) {
	// Get food_id and current rating from the review to update
	row := db.QueryRow("SELECT food_id, rating FROM reviews WHERE review_id = " + review_id)
	var food_id int
	var old_rating int
	row.Scan(&food_id, &old_rating)

	// Get current review coutn and score from foods table
	row = db.QueryRow("SELECT review_count, total_score FROM foods WHERE id = " + strconv.Itoa(food_id))
	var review_count int
	var total_score int
	row.Scan(&review_count, &total_score)

	// Calculate new average
	parsedRating, _ := strconv.ParseInt(rating, 10, 64)
	total_score = total_score - old_rating + int(parsedRating)
	new_average := strconv.FormatFloat(float64(total_score) / float64(review_count), 'E', 2, 64)

	// Update foods table
	stmt, err := db.Prepare("UPDATE foods SET total_score = ?, rating = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(total_score, new_average, food_id)
	if err != nil {
		log.Fatal(err)
	}

	// Update reviews table
	stmt, err = db.Prepare("UPDATE reviews SET rating = ?, review_text = ?, created_at = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(rating, review_text, created_at, review_id)
	if err != nil {
		log.Fatal(err)
	}

}

/*
	Checks to see if it is a weekend so that the website can render the proper menu
 */
func IsWeekend() bool {
	t := time.Now().Weekday()
	return t == 0 || t == 6
}