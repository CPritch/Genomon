package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	apiURL     = "https://api.dotgg.gg/cgfw/getcards?game=pokepocket&mode=indexed"
	outputFile = "cards.json"
)

type CardData struct {
	Data  [][]interface{} `json:"data"`
	Names []string        `json:"names"`
}

type Effect struct {
	Effect string `json:"effect"`
	Info   string `json:"info"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Demand struct {
	Amount string `json:"amount"`
	Image  string `json:"image"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

type Flair struct {
	Count             string   `json:"count"`
	Demands           []Demand `json:"demands"`
	FromDate          int64    `json:"from_date"`
	Image             string   `json:"image"`
	Name              string   `json:"name"`
	Slug              string   `json:"slug"`
	Type              string   `json:"type"`
	PrerequisiteCount *string  `json:"prerequisite_count,omitempty"`
	PrerequisiteName  *string  `json:"prerequisite_name,omitempty"`
}

type Route struct {
	Flairs    []Flair `json:"flairs"`
	RouteName string  `json:"routeName"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables.")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API key not found. Set API_KEY in .env")
	}

	// Fetch raw card data.
	rawData, err := fetchCards(apiKey)
	if err != nil {
		log.Fatalf("Failed to fetch cards: %v", err)
	}

	// Transform raw data into a slice of dictionaries.
	readableCards := transformCards(rawData)

	// Save the transformed data to file.
	if err := saveToFile(readableCards, outputFile); err != nil {
		log.Fatalf("Failed to save cards: %v", err)
	}

	fmt.Printf("Successfully saved %d cards to %s\n", len(readableCards), outputFile)
}

// fetchCards retrieves card data from the API and decodes it into a CardData type.
func fetchCards(apiKey string) (CardData, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return CardData{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return CardData{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CardData{}, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var data CardData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return CardData{}, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return data, nil
}

// transformCards maps each row in CardData.Data to a dictionary using CardData.Names.
func transformCards(raw CardData) []map[string]interface{} {
	var transformed []map[string]interface{}
	for _, row := range raw.Data {
		card := make(map[string]interface{})
		// Map each element of the row to its corresponding field name.
		for i, key := range raw.Names {
			if i < len(row) {
				card[key] = row[i]
			}
		}
		transformed = append(transformed, card)
	}
	return transformed
}

// saveToFile writes JSON data to a file.
func saveToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}
