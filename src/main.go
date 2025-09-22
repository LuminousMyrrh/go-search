package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"search/src/types"
	"strings"
	// "github.com/rs/zerolog/log"
)

func fetchForData() ([]types.Item, error) {
	api := "https://furniture-api.fly.dev"
	response, err := http.Get(api + "/v1/products")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to fetch for data")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(body))
	
	var resp types.Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: search <search request>")
		return
	}

	fmt.Println("Fetching data...")
	dataset, err := fetchForData()
	if err != nil {
		fmt.Println(err)
		return
	} else if len(dataset) == 0 {
		fmt.Println("API returned empty response")
		return
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
