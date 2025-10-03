package tcgdex

// This file defines the Go structs that match the JSON
// structure returned by the TCGdex API. This allows for easy and type-safe
// decoding of the data we fetch.

// Card represents a single, detailed card object from the TCGdex API.
type Card struct {
	ID             string     `json:"id"`
	LocalID        string     `json:"localId"`
	Name           string     `json:"name"`
	Image          string     `json:"image,omitempty"`
	Category       string     `json:"category"`
	Illustrator    string     `json:"illustrator,omitempty"`
	Rarity         string     `json:"rarity"`
	Set            Set        `json:"set"`
	Variants       *Variants  `json:"variants,omitempty"`
	HP             int        `json:"hp,omitempty"`
	Types          []string   `json:"types,omitempty"`
	EvolveFrom     string     `json:"evolveFrom,omitempty"`
	Stage          string     `json:"stage,omitempty"`
	Attacks        []Attack   `json:"attacks,omitempty"`
	Abilities      []Ability  `json:"abilities,omitempty"`
	Weaknesses     []Weakness `json:"weaknesses,omitempty"`
	Retreat        int        `json:"retreat,omitempty"`
	RegulationMark string     `json:"regulationMark,omitempty"`
	Legal          Legal      `json:"legal"`
	Text           string     `json:"effect,omitempty"` // For Trainer/Item cards
}

// Set contains basic information about the set a card belongs to.
type Set struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo,omitempty"`
}

// Variants indicates the different print variants of a card (holo, reverse, etc.).
type Variants struct {
	FirstEdition bool `json:"firstEdition"`
	Holo         bool `json:"holo"`
	Normal       bool `json:"normal"`
	Reverse      bool `json:"reverse"`
	WPromo       bool `json:"wPromo"`
}

// Attack defines an attack, including its cost, name, effect text, and damage.
type Attack struct {
	Cost   []string    `json:"cost"`
	Name   string      `json:"name"`
	Effect string      `json:"effect,omitempty"`
	Damage interface{} `json:"damage,omitempty"` // Can be int or string like "30+"
}

// Ability defines a Pokémon's ability.
type Ability struct {
	Name   string `json:"name"`
	Effect string `json:"effect"`
	Type   string `json:"type"`
}

// Weakness defines a weakness, including its type and value (e.g., "×2").
type Weakness struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Legal defines the legality of a card in different formats.
type Legal struct {
	Standard bool `json:"standard"`
	Expanded bool `json:"expanded"`
}

// SeriesDetails is used to decode the response from the /series/{id} endpoint.
type SeriesDetails struct {
	Sets []SetSummary `json:"sets"`
}

// SetSummary contains the basic info for a set as returned by the series endpoint.
type SetSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SetDetails represents the full details of a set, including a list of all its cards.
type SetDetails struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Cards []CardStump `json:"cards"`
}

// CardStump is a smaller Card struct used for decoding the list of cards in a set response.
type CardStump struct {
	ID      string `json:"id"`
	LocalID string `json:"localId"`
	Name    string `json:"name"`
}
