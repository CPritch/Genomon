package core

// Condition is an interface representing a single, distinct parameter or trigger.
type Condition interface {
	isCondition() // Marker method
}

// --- Concrete Condition Structs ---
// Each struct is a small, self-contained piece of logic.

type TriggerCondition struct {
	Trigger TriggerType `json:"trigger"`
}

func (c TriggerCondition) isCondition() {}

type ScalingCondition struct {
	ScaleBy ScaleByType `json:"scaleBy"`
	Amount  int         `json:"amount"`
}

func (c ScalingCondition) isCondition() {}

type CoinFlipCondition struct {
	Result CoinFlipResult `json:"result"`
}

func (c CoinFlipCondition) isCondition() {}

type DurationCondition struct {
	Duration DurationType `json:"duration"`
}

func (c DurationCondition) isCondition() {}

type SearchCondition struct {
	Target      TargetType  `json:"target"`
	Random      bool        `json:"random"`
	PokemonType PokemonType `json:"pokemonType"`
}

func (c SearchCondition) isCondition() {}

type EnergyCondition struct {
	RequiredExtraEnergyCount int    `json:"requiredExtraEnergyCount"`
	RequiredEnergyType       string `json:"requiredEnergyType"`
}

func (c EnergyCondition) isCondition() {}

type EnergyAttachCondition struct {
	RequiredEnergyType string `json:"requiredEnergyType"`
}

func (c EnergyAttachCondition) isCondition() {}

// --- ENUMS for True Type Safety ---

type TriggerType string

const (
	TriggerEvolvedThisTurn   TriggerType = "EVOLVED_THIS_TURN"
	TriggerPlayedSupporter   TriggerType = "PLAYED_SUPPORTER"
	TriggerOpponentHasStatus TriggerType = "OPPONENT_HAS_STATUS"
	TriggerAttachEnergySelf  TriggerType = "ATTACH_ENERGY_SELF"
	TriggerOponentHpGreater  TriggerType = "OPPONENT_HP_GREATER"
)

type ScaleByType string

const (
	ScaleByOpponentRetreatCost ScaleByType = "OPPONENT_RETREAT_COST"
	ScaleByFriendlyBenched     ScaleByType = "FRIENDLY_BENCHED_COUNT"
)

type CoinFlipResult string

const (
	CoinFlipHeads CoinFlipResult = "HEADS"
	CoinFlipTails CoinFlipResult = "TAILS"
)

type DurationType string

const (
	DurationNextTurn DurationType = "NEXT_TURN"
)

type PokemonType string

// TODO: Fill in all the different pokemon (card?) types...
