package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"flag"
	"search/src/migrations"
	"search/src/types"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/rs/zerolog/log"
)

func NewDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

	return db, nil
}

func fetchForData() ([]types.Item, error) {
	api := "https://furniture-api.fly.dev"
	response, err := http.Get(api + "/v1/products?limit=100")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got this status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	
	var resp types.Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func checkDataset(db *gorm.DB) ([]types.Item, error) {
	var items []types.Item
	res := db.Table("items").Find(&items)
	if res.Error != nil {
		return nil, res.Error
	}
	return items, nil
}

func writeDataset(dataset []types.Item, db *gorm.DB) error {
	for _, item := range dataset {
		if err := db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}

func wipeDataset(db *gorm.DB) error {
	return db.Where("1 = 1").Delete(&types.Item{}).Error
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file")
		return
	}

	db, err := NewDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := migrations.RunMigrations(db); err != nil {
		fmt.Println("Failed to run migrations: " + err.Error())
		return
	}

	if len(os.Args) == 1 {
		fmt.Println("Usage: search <search request>")
		return
	}

	forceRequest := flag.Bool("force-request", false, "force sending request and updating db")
	flag.Parse()

	fmt.Println("Fetching data...")
	dataset, err := checkDataset(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	if *forceRequest || len(dataset) == 0 {
		fmt.Println("Fetching data from API")
		
		dataset, err = fetchForData()
		// dataset2, err := fetchForData()
		// for _, data := range dataset2 {
		// 	dataset = append(dataset, data)
		// }

		if err != nil {
			fmt.Println(err)
			return
		} else if len(dataset) == 0 {
			fmt.Println("API returned empty response")
			return
		}
		
		err = wipeDataset(db)
		if err != nil {
			fmt.Println("Falied to wipe dataset: " + err.Error())
			return
		}
		err = writeDataset(dataset, db)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	e := NewEngine()

	request := strings.Join(os.Args[1:], " ")

	fmt.Println("Searching for: " + request)
	result := e.Search(os.Args[1:], dataset)
	if len(result) == 0 {
		fmt.Println("No results found")
	} else {
		fmt.Println("Found:")
		for i, r := range result {
			fmt.Printf("  %d: %s\n", i, r)
		}
	}
}
