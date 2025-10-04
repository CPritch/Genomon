package core

import "github.com/cpritch/genomon/pkg/tcgdex"

// Card represents our internal, enriched representation of a card. It embeds
// the raw card data from the TCGdex API and adds our own parsed effects.
type Card struct {
	tcgdex.Card
	ParsedAbilities []Effect `json:"parsedAbilities"`
	ParsedAttacks   []Effect `json:"parsedAttacks"`
}
