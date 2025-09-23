package main

import (
	"fmt"
	"search/src/types"
	"strings"
	"sync"
)

type Engine struct {
	
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Search(query []string, dataset []types.Item) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	found := []string{}
	wordChan := make(chan string, 100)

	chunks := e.sliceDataset(dataset, 100)

	go func() {
		for word := range wordChan {
			mu.Lock()
			found = append(found, word)
			mu.Unlock()
		}
	}()

	for _, chunk := range chunks {
		wg.Add(1)
		go func(set []types.Item) {
			defer wg.Done()
			for _, qe := range query {
				atl := e.lev(qe, set[0].Name)

				for _, s := range set[1:] {
					wordsInName := strings.SplitSeq(s.Name, " ")
					for w := range wordsInName {
						score := e.lev(qe, w)
						if atl > score {
							atl = score
							wordChan <- s.Name
							fmt.Println("Id: " + s.ID)
						}
					}
				}
			}
		}(chunk)
	}

	wg.Wait()
	close(wordChan)

	return found
}

func (e *Engine) sliceDataset(dataset []types.Item, chSize int) [][]types.Item {
	var chunks [][]types.Item
	for i := 0; i < len(dataset); i += chSize {
		end := i + chSize
		if end > len(dataset) {
			end = len(dataset)
		}
		chunks = append(chunks, dataset[i:end])
	}

	return chunks
}

func (e *Engine) lev(str1 string, str2 string) int {
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
