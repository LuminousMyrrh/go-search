package main

import (
	"search/src/types"
	"strings"
)

type Engine struct {
	
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Search(query []string, dataset []types.Item) []string {
	q := strings.Join(query, " ")
	
	found := []string{}
	atl := e.lev(q, dataset[0].Name)

	for _, s := range dataset[1:] {
		score := e.lev(q, s.Name)
		if atl > score && score < 5 {
			atl = score
			found = append(found, s.Name)
		}
	}

	return found
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
