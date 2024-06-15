package models

import "encoding/xml"

type Rates struct {
	XMLName xml.Name       `xml:"rates"`
	Items   []CurrencyItem `xml:"item"`
	Date    string         `xml:"date"`
}

type CurrencyItem struct {
	Title string `xml:"fullname"`
	Code  string `xml:"title"`
	Value string `xml:"description"`
}
