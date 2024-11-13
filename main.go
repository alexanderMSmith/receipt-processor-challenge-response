package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// The structure for keeping track of Receipt information.
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Id           string `json:"id"`
	Total        string `json:"total"`
	Score        int    `json:"points"`
	Items        []Item `json:"items"`
}

// The structure for tracking item information.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// THe test method.
func main() {
	jsonData := []byte(`{
		"retailer": "Target",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "13:01",
		"items": [
			{
			"shortDescription": "Mountain Dew 12PK",
			"price": "6.49"
			},{
			"shortDescription": "Emils Cheese Pizza",
			"price": "12.25"
			},{
			"shortDescription": "Knorr Creamy Chicken",
			"price": "1.26"
			},{
			"shortDescription": "Doritos Nacho Cheese",
			"price": "3.35"
			},{
			"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
			"price": "12.00"
			}
		],
		"total": "35.35"
	}`)

	receipt := calculate(jsonData)
	fmt.Printf("%+v\n", receipt.Score)
}

func calculate(in []byte) Receipt {
	receipt := Receipt{}
	err := json.Unmarshal(in, &receipt)
	if err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}

	/* This segment converts the receipts retailer into alphanumerical characters
	   and then adds score based on its length. */
	bytRetailer := []byte(receipt.Retailer)
	alphaRetailer := clean(bytRetailer)
	fmt.Printf("Cleaned retailer: '%s'\n", alphaRetailer)
	receipt.Score += len(alphaRetailer)
	fmt.Printf("Retailer points (alpha length): %d\n", len(alphaRetailer))

	// This segment converts the last 2 digits of the price into an integer.
	lastTwo := receipt.Total[len(receipt.Total)-2:]
	cents, err := strconv.Atoi(lastTwo)
	if err != nil {
		log.Fatalf("Unable to calculate cents due to %s", err)
	}
	// This segment parses the cents and adds score if appropriate.
	if cents == 0 {
		receipt.Score += 50 // 50 points for perfect 0
	}
	if cents%25 == 0 {
		receipt.Score += 25 // 25 points for multiple of .25
	}
	fmt.Printf("Total points (round dollar and multiple of 0.25): %d\n", receipt.Score)

	// This line adds 5 points for every 2 items on the receipt.
	receipt.Score += 5 * (len(receipt.Items) / 2)
	fmt.Printf("Item count points (5 points per pair): %d\n", 5*(len(receipt.Items)/2))

	// This segment goes through every item to add points.
	for i := 0; i < len(receipt.Items); i++ {
		// Trims the item description.
		trimName := strings.TrimSpace(receipt.Items[i].ShortDescription)
		fmt.Printf("Trimmed item description: %s (Length: %d)\n", trimName, len(trimName))

		// Only apply points calculation if the length of the trimmed name is a multiple of 3.
		if len(trimName)%3 == 0 {
			cost, err := strconv.ParseFloat(receipt.Items[i].Price, 64)
			if err != nil {
				log.Fatalf("Unable to calculate item price due to %s", err)
			}
			preTotal := cost * .2
			receipt.Score += int(math.Ceil(preTotal))
			fmt.Printf("Item %d points (price * 0.2): %d\n", i+1, int(math.Ceil(preTotal)))
		}
	}

	// This segment calculates the day of the purchase date and converts it to an int.
	lastOfDate := receipt.PurchaseDate[len(receipt.PurchaseDate)-2:]
	day, err := strconv.Atoi(lastOfDate)
	if err != nil {
		log.Fatalf("Unable to calculate day due to %s", err)
	}

	// This segment adds 6 points if the purchase date is odd.
	if day%2 == 1 {
		receipt.Score += 6
	}
	fmt.Printf("Day of the month points (odd day): %d\n", receipt.Score)

	// This segment converts the first two digits of purchaseTime to integers.
	firstTwo := receipt.PurchaseTime[:2]
	hours, err := strconv.Atoi(firstTwo)
	if err != nil {
		log.Fatalf("Unable to calculate hours due to %s", err)
	}

	// This segments adds 10 points if the purchase time is between 14:00 and 16:00
	if hours >= 14 && hours <= 16 {
		receipt.Score += 10
	}
	fmt.Printf("Purchase time points (2pm-4pm): %d\n", receipt.Score)

	// Generates the ID.
	receipt.Id = generateID()

	return receipt
}

/*
Cleans the input to only use alphanumeric characters. Obtained and modified to exclude
spaces from peterSO on Stack Overflow. https://stackoverflow.com/a/54463943
*/
func clean(s []byte) string {
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}

func generateID() string {
	return uuid.NewString()
}
