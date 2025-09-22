package types

type Item struct {
	ID string      `json:"id"`
	Name string `json:"name"`
	Price float32   `json:"price"`
}

type Response struct {
	Data []Item `json:"data"`
}
