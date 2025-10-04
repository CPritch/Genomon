package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

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
	sampleSize := processCmd.Int("n", 0, "Number of random unknown effects to sample and print")

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
		handleProcessCommand(processInputFile, processOutputFile, sampleSize)
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

func handleProcessCommand(inputFile, outputFile *string, sampleSize *int) {
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
	var unknownCards []core.Card // Slice to store cards with unknown effects

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
				parsedEffects := effects.Parse(ability.Effect)
				for _, effect := range parsedEffects {
					effect.Name = ability.Name // Add name for mapping
					if effect.Type == core.EffectUnknown {
						cardHasUnknownEffect = true
					}
					enrichedCard.ParsedAbilities = append(enrichedCard.ParsedAbilities, effect)
				}
			}
		}

		// Parse attacks
		for _, attack := range rawCard.Attacks {
			if attack.Effect != "" {
				parsedEffects := effects.Parse(attack.Effect)
				for _, effect := range parsedEffects {
					effect.Name = attack.Name // Add name for mapping
					if effect.Type == core.EffectUnknown {
						cardHasUnknownEffect = true
					}
					enrichedCard.ParsedAttacks = append(enrichedCard.ParsedAttacks, effect)
				}
			}
		}

		if cardHasUnknownEffect {
			unknownCards = append(unknownCards, enrichedCard)
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

	if len(unknownCards) > 0 {
		fmt.Printf("\n⚠️  Warning: Could not parse one or more effects for %d card(s).\n", len(unknownCards))

		if sampleSize != nil && *sampleSize > 0 {
			fmt.Printf("--- Sampling %d random unknown effect(s) ---\n\n", *sampleSize)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			r.Shuffle(len(unknownCards), func(i, j int) {
				unknownCards[i], unknownCards[j] = unknownCards[j], unknownCards[i]
			})

			numToShow := *sampleSize
			if len(unknownCards) < numToShow {
				numToShow = len(unknownCards)
			}

			for i := 0; i < numToShow; i++ {
				card := unknownCards[i]
				fmt.Printf("Card: %s (%s)\n", card.Name, card.ID)
				for _, effect := range card.ParsedAbilities {
					if effect.Type == core.EffectUnknown {
						fmt.Printf("  └─ Ability '%s': %s\n", effect.Name, effect.Description)
					}
				}
				for _, effect := range card.ParsedAttacks {
					if effect.Type == core.EffectUnknown {
						fmt.Printf("  └─ Attack '%s': %s\n", effect.Name, effect.Description)
					}
				}
				fmt.Println()
			}
		}
	}
}
