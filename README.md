# genomon

A genetic learning system for discovering optimal deck strategies in the Pokémon Trading Card Game Pocket, written in Go.

## Project Status: Data Pipeline Complete

**Genomon** is currently in the foundational stage. The primary focus has been to build a robust and reliable data pipeline to source and understand the game's core components: the cards themselves.

We have successfully completed the first major phase:

1.  **Data Acquisition:** A client that syncs all official Pokémon TCG Pocket card data from the public [TCGdex API](https://www.google.com/search?q=https://tcgdex.net/docs/pocket). Many thanks to the great devs there.
2.  **Effect Parser:** A comprehensive, regex-based "compiler" that reads the natural language text of every card's abilities and attacks and converts it into a structured, machine-readable format. As of October 2025, this parser has **100% coverage** of all known card effects in the game. We are yet to step over and simplify this and confirm accuracy here.

The project now has a solid foundation of structured data, ready for the next phase: building the game simulation engine.

## Usage: The Data Pipeline

The application currently functions as a command-line tool to generate the card data needed for future steps.

### Prerequisites

  - Go (version 1.18 or later) installed on your system.

### Step 1: Sync Raw Card Data

First, you need to fetch the latest card data from the TCGdex API. This command contacts the API, retrieves all card sets for TCG Pocket, and saves the raw data.

```bash
go run ./cmd/genomon sync
```

This will create a file named `ptcgp-cards.json` in the project root. This file is the raw source of truth.

### Step 2: Process and Enrich Card Data

Next, run the effect parser. This command reads the raw `ptcgp-cards.json`, interprets every attack and ability, and saves a new, enriched file.

```bash
go run ./cmd/genomon process
```

This creates `genomon-cards.json`, which contains all the original card data plus the structured `parsedAbilities` and `parsedAttacks` fields. This is the core dataset that the future game simulation engine will use.

You can also sample the data for any effects the parser might have missed (currently none\!):

```bash
# Sample 5 random cards with effects that could not be parsed
go run ./cmd/genomon process -n 5
```

### ⚠️ Disclaimer on Effect Accuracy

The effect parser is a complex, hand-tuned system designed to cover all known card effects. While it has 100% coverage, the interpretation of nuanced effects may contain subtle inaccuracies. The logic is rule-based and has not yet been battle-tested in a live simulation. Verification and refinement of the parsed effects will be an ongoing process.

## Roadmap

The completion of the data pipeline marks the end of Phase 1. The next phases will focus on bringing the game and the genetic algorithm to life.

  - [x] **Phase 1: Data Foundation**

      - [x] Build a reliable client for the TCGdex API.
      - [x] Create Go structs to represent all card data.
      - [x] Develop a rule-based effect parser to convert card text to structured data.
      - [x] Achieve 100% parsing coverage for all known cards.
      - [ ] Simplify some of the regex parsing and confirm accuracy of card effects

  - [ ] **Phase 2: Game Simulation**

      - [ ] Implement the core game logic engine based on TCG Pocket rules.
      - [ ] Model game state (players, decks, hands, bench, active Pokémon, etc.).
      - [ ] Create a "headless" simulation environment where games can be played programmatically.
      - [ ] Build a simple "Random Agent" that can play legal moves randomly to test the engine.

  - [ ] **Phase 3: AI Player Development**

      - [ ] Develop a more intelligent Heuristic Agent that follows simple strategies.
      - [ ] Create a framework for pitting agents against each other.
      - [ ] (Optional) Explore a Reinforcement Learning agent for more advanced play.

  - [ ] **Phase 4: Genetic Algorithm for Deck Building**

      - [ ] Define the "genome" for a deck (a list of cards with constraints).
      - [ ] Create the initial population of random (but valid) decks.
      - [ ] Define a fitness function (e.g., win rate against a benchmark agent).
      - [ ] Implement genetic operators (crossover, mutation) for creating new decks.
      - [ ] Run large-scale simulations to evolve and discover optimal deck archetypes.