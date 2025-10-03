package core

import "github.com/cpritch/genomon/pkg/tcgdex"

// Card represents our internal, enriched representation of a card. It embeds
// the raw card data from the TCGdex API and adds our own parsed effects.
type Card struct {
	tcgdex.Card
	ParsedAbilities []Effect `json:"parsedAbilities"`
	ParsedAttacks   []Effect `json:"parsedAttacks"`
}

// Effect is our internal structured representation of an attack or ability's effect.
type Effect struct {
	Name        string                 `json:"name"`
	Type        EffectType             `json:"type"`
	Target      TargetType             `json:"target,omitempty"`
	Amount      int                    `json:"amount,omitempty"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
	Description string                 `json:"description"`
}

// EffectType is an enum for the different kinds of effects we can parse.
type EffectType string

const (
	EffectHeal                  EffectType = "HEAL"
	EffectDamage                EffectType = "DAMAGE"
	EffectApplyStatus           EffectType = "APPLY_STATUS"
	EffectRestrictionCantAttack EffectType = "RESTRICTION_CANT_ATTACK"
	EffectForceSwitch           EffectType = "FORCE_SWITCH"
	EffectSearchDeck            EffectType = "SEARCH_DECK"
	EffectRecoilDamage          EffectType = "RECOIL_DAMAGE"
	EffectConditionalDamage     EffectType = "CONDITIONAL_DAMAGE"
	EffectUnknown               EffectType = "UNKNOWN"
)

// TargetType defines who the effect applies to.
type TargetType string

const (
	TargetSelf           TargetType = "SELF"
	TargetOpponentActive TargetType = "OPPONENT_ACTIVE"
	TargetAllFriendly    TargetType = "ALL_FRIENDLY"
	TargetDeck           TargetType = "DECK"
)
