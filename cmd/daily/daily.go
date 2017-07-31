package main

import (
	"log"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"encoding/json"
	"time"
	"net/url"
	"github.com/kevinsangholee/ClaremontMenuGo"
	"database/sql"
	"strconv"
	"strings"
)

type Menu struct {
	DiningHall string   `json:"dining_hall"`
	Meal	   string   `json:"meal"`
	FoodItems  []string `json:"food_items"`
}

type SingleBing struct {
	ContentUrl string `json:"contentUrl"`
}

type BingResult struct {
	Values []SingleBing `json:"value"`
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

	// Set HTTP Client for Bing
	client := &http.Client{Timeout: 10 * time.Second}

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

	// Reset daily
	stmt, err := db.Prepare("UPDATE foods SET daily = '.'")
	res, err := stmt.Exec()
	num, err := res.RowsAffected()
	log.Println("Reset daily for " + strconv.Itoa(int(num)) + " food items!")

	// Iterate through each meal
	for _, menuItem := range menus {
		schoolInt := schoolToInt(menuItem.DiningHall)
		mealInt := mealToInt(menuItem.Meal)
		for _, foodName := range menuItem.FoodItems {
			// See if this item is in the database already
			query := "SELECT daily, id, image FROM foods WHERE school = " + strconv.Itoa(schoolInt) +
				 " AND name = \"" + foodName + "\""
			row := db.QueryRow(query)
			var daily string
			var id int
			var image string
			err := row.Scan(&daily, &id, &image)
			// If it doesn't exist, find the image for it add it to the database
			if err != nil {
				if err == sql.ErrNoRows {
					imageUrl := queryBingImageSearch(client, foodName)
					stmt, err := db.Prepare("INSERT INTO foods (name, school, image, daily) VALUES (?,?,?,?)")
					if err != nil {
						log.Fatal(err)
					}
					_, err = stmt.Exec(foodName, schoolInt, imageUrl, "." + strconv.Itoa(mealInt))
					if err != nil {
						log.Fatal(err)
					}
				} else {
					log.Fatal(err)
				}
				log.Println("Added new food " + foodName + "!")
			// If it does exist, first check to see if there is an image for it.
			} else {
				newDaily := daily + strconv.Itoa(mealInt)
				// If not, update daily with new image url
				if image == "null" {

					stmt, err := db.Prepare("UPDATE foods SET daily = ?, image = ? WHERE id = ?")
					if err != nil {
						log.Fatal(err)
					}
					imageUrl := queryBingImageSearch(client, foodName)
					_, err = stmt.Exec(newDaily, imageUrl, id)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("Updated the daily for " + foodName + "as well as updated the image!")
				// If yes, just update daily
				} else {

					stmt, err := db.Prepare("UPDATE foods SET daily = ? WHERE id = ?")
					if err != nil {
						log.Fatal(err)
					}
					_, err = stmt.Exec(newDaily, id)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("Updated the daily for " + foodName + "!")

				}
			}
		}
	}
}

func queryBingImageSearch(client *http.Client, foodName string) string {

	values := url.Values {
		"q": []string{foodName},
		"size": []string{"Medium"},
	}

	req, err := http.NewRequest("GET", "https://api.cognitive.microsoft.com/bing/v5.0/images/search/?" + values.Encode(), nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Ocp-Apim-Subscription-Key", "6229f6930bc84d4a89a29215b1dc9f47")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close();

	data := BingResult{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	return data.Values[0].ContentUrl

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