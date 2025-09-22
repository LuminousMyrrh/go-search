package main

import (
	"fmt"
	"os"
	"strings"
	// "github.com/rs/zerolog/log"
)

var DATASET []string = []string{
	"lether sofa",
	"comfy chair",
	"huge table",
	"small table",
	"cute table",
	"cute sofa",
	"cute chair",
	"puffy couch",
	"coffee table",
	"family table",
	"spacious wardrobe",
}

func lev(str1 string, str2 string) int {
	m := len(str1)
	n := len(str2)

	d := make([][]int, m + 1)
	for i := 0; i <= m; i++ {
		d[i] = make([]int, n + 1)
	}

	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	 
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	subCost := 0
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {

			if str1[i-1] == str2[j-1] {
				subCost = d[i-1][j-1]
			} else {
				subCost = min(
					d[i-1][j],
					d[i][j-1],
					d[i-1][j-1],
					) + 1
			}

			d[i][j] = subCost
		}
	}

	return d[m][n]
}

func search(query []string) []string {
	q := strings.Join(query, " ")
	
	found := []string{}
	atl := lev(q, DATASET[0])

	for _, s := range DATASET[1:] {
		score := lev(q, s)
		if atl > score {
			atl = score
			found = append(found, s)
		}
	}

	return found
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: search <search request>")
		return
	}

	request := strings.Join(os.Args[1:], " ")

	fmt.Println("Your request is: " + request)

	result := search(os.Args[1:])
	if len(result) == 0 {
		fmt.Println("No results found")
	} else {
		fmt.Println("Found:")
		for i, r := range result {
			fmt.Printf("  %d: %s\n", i, r)
		}
	}
}
