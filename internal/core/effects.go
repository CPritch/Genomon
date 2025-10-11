package core

import (
	"fmt"
	"strings"
)

// Effect is our internal structured representation of an attack or ability's effect.
type Effect struct {
	Name        string          `json:"name"`
	Type        EffectType      `json:"type"`
	Passive     bool            `json:"passive,omitempty"`
	Target      TargetType      `json:"target,omitempty"`
	Status      StatusCondition `json:"status,omitempty"`
	Amount      int             `json:"amount,omitempty"`
	Times       int             `json:"times,omitempty"`
	Conditions  []Condition     `json:"conditions,omitempty"`
	Description string          `json:"description"`
}

// EffectType is an enum for the different kinds of effects we can parse.
type EffectType string

const (
	EffectHeal                     EffectType = "HEAL"
	EffectDraw                     EffectType = "DRAW"
	EffectDamage                   EffectType = "DAMAGE"
	EffectCopyAttack               EffectType = "COPY_ATTACK"
	EffectApplyStatus              EffectType = "APPLY_STATUS"
	EffectRestrictionCantAttack    EffectType = "RESTRICTION_CANT_ATTACK"
	EffectForceSwitch              EffectType = "FORCE_SWITCH"
	EffectSearchDeck               EffectType = "SEARCH_DECK"
	EffectRecoilDamage             EffectType = "RECOIL_DAMAGE"
	EffectConditionalDamage        EffectType = "CONDITIONAL_DAMAGE"
	EffectAttachEnergy             EffectType = "ATTACH_ENERGY"
	EffectTriggeredAbility         EffectType = "TRIGGERED_ABILITY"
	EffectScalingDamage            EffectType = "SCALING_DAMAGE"
	EffectAttackMayFail            EffectType = "ATTACK_MAY_FAIL"
	EffectDiscardEnergy            EffectType = "DISCARD_ENERGY"
	EffectMoveEnergy               EffectType = "MOVE_ENERGY"
	EffectReduceIncomingDamage     EffectType = "REDUCE_INCOMING_DAMAGE"
	EffectDiscardCard              EffectType = "DISCARD_CARD"
	EffectApplyRestriction         EffectType = "APPLY_RESTRICTION"
	EffectMultiHitRandomDamage     EffectType = "MULTI_HIT_RANDOM_DAMAGE"
	EffectDamageBenchedFriendly    EffectType = "DAMAGE_BENCHED_FRIENDLY"
	EffectSnipeDamage              EffectType = "SNIPE_DAMAGE"
	EffectSwitchSelf               EffectType = "SWITCH_SELF"
	EffectShuffleIntoDeck          EffectType = "SHUFFLE_INTO_DECK"
	EffectApplyPrevention          EffectType = "APPLY_PREVENTION"
	EffectScalingSnipeDamage       EffectType = "SCALING_SNIPE_DAMAGE"
	EffectDamageBenchedOpponentAll EffectType = "DAMAGE_BENCHED_OPPONENT_ALL"
	EffectLifesteal                EffectType = "LIFESTEAL"
	EffectApplyReactiveDamage      EffectType = "APPLY_REACTIVE_DAMAGE"
	EffectBuffNextTurn             EffectType = "BUFF_NEXT_TURN"
	EffectModifyEnergy             EffectType = "MODIFY_ENERGY"
	EffectDamageAllOpponent        EffectType = "DAMAGE_ALL_OPPONENT"
	EffectUnknown                  EffectType = "UNKNOWN"
	EffectSetHP                    EffectType = "SET_HP"
	EffectShuffleFromHand          EffectType = "SHUFFLE_FROM_HAND"
	EffectLookAtDeck               EffectType = "LOOK_AT_DECK"
	EffectDelayedDamage            EffectType = "DELAYED_DAMAGE"
	EffectKnockout                 EffectType = "KNOCKOUT"
	EffectMoveDamage               EffectType = "MOVE_DAMAGE"
	EffectRevealHand               EffectType = "REVEAL_HAND"
	EffectDamageHalveHP            EffectType = "DAMAGE_HALVE_HP"
	EffectReturnToHand             EffectType = "RETURN_TO_HAND"
	EffectDevolve                  EffectType = "DEVOLVE"
	EffectDebuffIncomingDamage     EffectType = "DEBUFF_INCOMING_DAMAGE"
	EffectRemoveRetreatCost        EffectType = "REMOVE_RETREAT_COST"
)

// TargetType defines who the effect applies to.
type TargetType string

const (
	TargetSelf               TargetType = "SELF"
	TargetOpponentActive     TargetType = "OPPONENT_ACTIVE"
	TargetOpponentHand       TargetType = "OPPONENT_HAND"
	TargetAllFriendly        TargetType = "ALL_FRIENDLY"
	TargetAllOpponent        TargetType = "ALL_OPPONENT"
	TargetAllPokemonInPlay   TargetType = "ALL_POKEMON_IN_PLAY"
	TargetBenchedFriendly    TargetType = "BENCHED_FRIENDLY"
	TargetBenchedOpponent    TargetType = "BENCHED_OPPONENT"
	TargetBenchedOpponentAll TargetType = "BENCHED_OPPONENT_ALL"
	TargetDeck               TargetType = "DECK"
	TargetHand               TargetType = "HAND"
	TargetEnergyZone         TargetType = "ENERGY_ZONE"
	TargetEnemyEnergyZone    TargetType = "ENEMY_ENERGY_ZONE"
	TargetGraveyard          TargetType = "GRAVEYARD"
	TargetOpponentGraveyard  TargetType = "OPPONENT_GRAVEYARD"
)

// StatusCondition represents the special conditions in the game.
type StatusCondition string

const (
	StatusPoisoned  StatusCondition = "POISONED"
	StatusConfused  StatusCondition = "CONFUSED"
	StatusAsleep    StatusCondition = "ASLEEP"
	StatusBurned    StatusCondition = "BURNED"
	StatusParalyzed StatusCondition = "PARALYZED"
)

func ParseStatusCondition(s string) (StatusCondition, error) {
	// Normalize the input string by trimming whitespace and converting to lowercase
	// for case-insensitive comparison.
	normalizedInput := strings.ToLower(strings.TrimSpace(s))

	switch normalizedInput {
	case "poisoned":
		return StatusPoisoned, nil
	case "confused":
		return StatusConfused, nil
	case "asleep":
		return StatusAsleep, nil
	case "burned":
		return StatusBurned, nil
	case "paralyzed":
		return StatusParalyzed, nil
	default:
		// If no case matches, return an empty StatusCondition and an error.
		return "", fmt.Errorf("invalid status condition: %q", s)
	}
}
