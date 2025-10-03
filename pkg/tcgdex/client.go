package tcgdex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	apiBaseURL   = "https://api.tcgdex.net/v2/en"
	pocketSeries = "tcgp"
)

// Client is a client for interacting with the TCGdex API.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new TCGdex API client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchTCGPSetIDs fetches all set IDs belonging to the TCG Pocket series.
func (c *Client) FetchTCGPSetIDs() ([]string, error) {
	seriesURL := fmt.Sprintf("%s/series/%s", apiBaseURL, pocketSeries)
	req, err := http.NewRequest("GET", seriesURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for TCGP series: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TCGP series: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code when fetching TCGP series: %d", resp.StatusCode)
	}

	var seriesDetails SeriesDetails
	if err := json.NewDecoder(resp.Body).Decode(&seriesDetails); err != nil {
		return nil, fmt.Errorf("failed to decode TCGP series details: %w", err)
	}

	var setIDs []string
	for _, set := range seriesDetails.Sets {
		setIDs = append(setIDs, set.ID)
	}

	return setIDs, nil
}

// FetchSetCards fetches the full details for every card in a given set.
func (c *Client) FetchSetCards(setID string) ([]Card, error) {
	// First, fetch the set to get the list of card IDs
	setURL := fmt.Sprintf("%s/sets/%s", apiBaseURL, setID)
	req, err := http.NewRequest("GET", setURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for set %s: %w", setID, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch set %s: %w", setID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code when fetching set %s: %d", setID, resp.StatusCode)
	}

	var setDetails SetDetails
	if err := json.NewDecoder(resp.Body).Decode(&setDetails); err != nil {
		return nil, fmt.Errorf("failed to decode set details for %s: %w", setID, err)
	}

	var fullCards []Card
	// Now, fetch each card individually to get full details
	for _, cardSummary := range setDetails.Cards {
		fmt.Printf("Fetching card: %s (%s)\n", cardSummary.Name, cardSummary.ID)
		card, err := c.FetchCard(cardSummary.ID)
		if err != nil {
			// Log the error but continue trying to fetch other cards
			fmt.Printf("Warning: failed to fetch card %s: %v\n", cardSummary.ID, err)
			continue
		}
		// Small delay to be respectful to the API
		time.Sleep(50 * time.Millisecond)
		fullCards = append(fullCards, *card)
	}

	return fullCards, nil
}

// FetchCard fetches a single card by its full ID (e.g., "pock-1").
func (c *Client) FetchCard(cardID string) (*Card, error) {
	cardURL := fmt.Sprintf("%s/cards/%s", apiBaseURL, cardID)
	req, err := http.NewRequest("GET", cardURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for card %s: %w", cardID, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch card %s: %w", cardID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code when fetching card %s: %d", cardID, resp.StatusCode)
	}

	var card Card
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, fmt.Errorf("failed to decode card %s: %w", cardID, err)
	}

	return &card, nil
}
