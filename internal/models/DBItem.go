package models

type DBItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code"`
	Value string `json:"value"`
	Date  string `json:"date"`
}
