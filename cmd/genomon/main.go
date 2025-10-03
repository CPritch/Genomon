package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/cpritch/genomon/internal/core"
	"github.com/cpritch/genomon/internal/effects"
	"github.com/cpritch/genomon/pkg/tcgdex"
)

const (
	rawOutputFile      = "ptcgp-cards.json"
	enrichedOutputFile = "genomon-cards.json"
)

func main() {
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	syncOutputFile := syncCmd.String("o", rawOutputFile, "Output file for the synced card data")

	processCmd := flag.NewFlagSet("process", flag.ExitOnError)
	processInputFile := processCmd.String("i", rawOutputFile, "Input file for processing")
	processOutputFile := processCmd.String("o", enrichedOutputFile, "Output file for processed data")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sync":
		syncCmd.Parse(os.Args[2:])
		handleSyncCommand(syncOutputFile)
	case "process":
		processCmd.Parse(os.Args[2:])
		handleProcessCommand(processInputFile, processOutputFile)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: genomon <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  sync       Fetches the latest card data from the TCGdex API.")
	fmt.Println("    -o <file>    Output file for card data (default: ptcgp-cards.json)")
	fmt.Println("\n  process    Parses effects from raw card data into a structured format.")
	fmt.Println("    -i <file>    Input file for processing (default: ptcgp-cards.json)")
	fmt.Println("    -o <file>    Output file for processed data (default: genomon-cards.json)")
}

// ... existing handleSyncCommand code ...
func handleSyncCommand(outputFile *string) {
	fmt.Println("Starting card data sync from TCGdex...")
	client := tcgdex.NewClient()

	fmt.Println("Fetching TCG Pocket set list...")
	pocketSetIDs, err := client.FetchTCGPSetIDs()
	if err != nil {
		fmt.Printf("Error fetching TCG Pocket set list: %v\n", err)
		os.Exit(1)
	}

	if len(pocketSetIDs) == 0 {
		fmt.Println("No TCG Pocket sets found. Exiting.")
		os.Exit(0)
	}

	fmt.Printf("Found %d sets: %v\n", len(pocketSetIDs), pocketSetIDs)

	allCards := make(map[string]tcgdex.Card) // Use a map to prevent duplicates

	for _, setID := range pocketSetIDs {
		fmt.Printf("\n--- Fetching set: %s ---\n", setID)
		setCards, err := client.FetchSetCards(setID)
		if err != nil {
			fmt.Printf("Error fetching set %s: %v\n", setID, err)
			continue
		}
		for _, card := range setCards {
			allCards[card.ID] = card
		}
		fmt.Printf("--- Finished set: %s ---\n", setID)
	}

	// Convert map to slice for sorting
	cardSlice := make([]tcgdex.Card, 0, len(allCards))
	for _, card := range allCards {
		cardSlice = append(cardSlice, card)
	}

	// Sort cards by ID for a consistent output file
	sort.Slice(cardSlice, func(i, j int) bool {
		return cardSlice[i].ID < cardSlice[j].ID
	})

	fmt.Printf("\nTotal unique cards fetched: %d\n", len(cardSlice))

	// Write the data to the output file
	fileData, err := json.MarshalIndent(cardSlice, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling card data to JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, fileData, 0644)
	if err != nil {
		fmt.Printf("Error writing to output file %s: %v\n", *outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully synced all card data to %s\n", *outputFile)
}

func handleProcessCommand(inputFile, outputFile *string) {
	fmt.Printf("Loading raw card data from %s...\n", *inputFile)
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	var rawCards []tcgdex.Card
	if err := json.Unmarshal(data, &rawCards); err != nil {
		fmt.Printf("Error unmarshalling raw card data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing %d cards to parse effects...\n", len(rawCards))
	enrichedCards := make([]core.Card, 0, len(rawCards))
	unknownEffectCount := 0

	for _, rawCard := range rawCards {
		enrichedCard := core.Card{
			Card:            rawCard,
			ParsedAbilities: []core.Effect{},
			ParsedAttacks:   []core.Effect{},
		}

		cardHasUnknownEffect := false

		// Parse abilities
		for _, ability := range rawCard.Abilities {
			if ability.Effect != "" {
				parsedEffect := effects.Parse(ability.Effect)
				parsedEffect.Name = ability.Name // Add name for mapping
				if parsedEffect.Type == core.EffectUnknown {
					cardHasUnknownEffect = true
				}
				enrichedCard.ParsedAbilities = append(enrichedCard.ParsedAbilities, parsedEffect)
			}
		}

		// Parse attacks
		for _, attack := range rawCard.Attacks {
			if attack.Effect != "" {
				parsedEffect := effects.Parse(attack.Effect)
				parsedEffect.Name = attack.Name // Add name for mapping
				if parsedEffect.Type == core.EffectUnknown {
					cardHasUnknownEffect = true
				}
				enrichedCard.ParsedAttacks = append(enrichedCard.ParsedAttacks, parsedEffect)
			}
		}

		if cardHasUnknownEffect {
			unknownEffectCount++
		}

		enrichedCards = append(enrichedCards, enrichedCard)
	}

	// Write the enriched data to the output file
	fileData, err := json.MarshalIndent(enrichedCards, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling enriched card data to JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, fileData, 0644)
	if err != nil {
		fmt.Printf("Error writing to output file %s: %v\n", *outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed and saved enriched card data to %s\n", *outputFile)

	if unknownEffectCount > 0 {
		fmt.Printf("\n⚠️  Warning: Could not parse one or more effects for %d card(s).\n", unknownEffectCount)
		fmt.Println("   These have been marked with type UNKNOWN in the output file.")
	}
}
