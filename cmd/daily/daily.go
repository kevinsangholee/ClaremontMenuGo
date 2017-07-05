package main

import (
	"log"
	"github.com/kevinsangholee/ClaremontMenuGo"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
	"net/url"
	"encoding/json"
	"time"
	"strings"
	"strconv"
)

type Menu struct {
	DiningHall string   `json:"dining_hall"`
	Meal	   string   `json:"meal"`
	FoodItems  []string `json:"food_items"`
}

const apiURL = "https://aspc.pomona.edu/api/menu/day/"

func main() {
	log.Println("Adding daily for today!")

	// Converts current day of week into 3 letter string
	day := strings.ToLower(time.Now().Weekday().String()[0:3])

	// Authentication token set
	values := url.Values {
		"auth_token": []string{"2a9c963d749d6de4933579c611723625b74521c0"},
	}

	// Get Request for current day
	resp, err := http.Get(apiURL + day + "/?" + values.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Response check
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Request was not OK: " + resp.Status)
	}

	// Decoding the JSON into menus
	menus := make([]Menu, 0)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&menus)
	if err != nil {
		log.Fatal(err)
	}

	// Open Database Connection
	dsn := menudb.DB_USER + ":" + menudb.DB_PASS + "@" + menudb.DB_HOST + "/" + menudb.DB_NAME
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database opened!")
	defer db.Close()

	// Iterate through each meal
	for _, menuItem := range menus {
		schoolInt := schoolToInt(menuItem.DiningHall)
		mealInt := mealToInt(menuItem.Meal)
		for _, foodName := range menuItem.FoodItems {
			// See if this item is in the database already
			query := "SELECT daily, id FROM foods WHERE school = " + strconv.Itoa(schoolInt) +
				 " AND name = \"" + foodName + "\""
			row := db.QueryRow(query)
			var daily string
			var id int
			err := row.Scan(&daily, &id)
			if err != nil {
				if err == sql.ErrNoRows {
					daily = "DOESNT EXIST"
				} else {
					log.Fatal(err)
				}
			} else {
				newDaily := daily + strconv.Itoa(mealInt)
				stmt, err := db.Prepare("UPDATE foods SET daily = ? WHERE id = ?")
				if err != nil {
					log.Fatal(err)
				}
				_, err = stmt.Exec(newDaily, id)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Added " + foodName + " to daily!")
			}
		}
	}
}

func schoolToInt(school string) int {
	switch school {
	case "frank":
		return 0
	case "frary":
		return 1
	case "oldenborg":
		return 2
	case "cmc":
		return 3
	case "scripps":
		return 4
	case "pitzer":
		return 5
	case "mudd":
		return 6
	default:
		return -1
	}
}

func mealToInt(meal string) int {
	switch meal {
	case "breakfast":
		return 0
	case "lunch":
		return 1
	case "dinner":
		return 2
	case "brunch":
		return 3
	default:
		return -1
	}
}