package effects

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cpritch/genomon/internal/core"
)

// These are simple regular expressions to find keywords in effect text.
// We will expand this list significantly.
var (
	healRegex               = regexp.MustCompile(`Heal (\d+) damage from this Pokémon\.`)
	cantAttackNextTurnRegex = regexp.MustCompile(`During your next turn, this Pokémon can't attack\.`)
	forceSwitchRegex        = regexp.MustCompile(`Switch out your opponent's Active Pokémon to the Bench\.`)
	searchDeckRegex         = regexp.MustCompile(`Put (\d+) random {([A-Z])} Pokémon from your deck into your hand\.`)
	recoilDamageRegex       = regexp.MustCompile(`This Pokémon also does (\d+) damage to itself\.`)
	healAllFriendlyRegex    = regexp.MustCompile(`heal (\d+) damage from each of your Pokémon\.`)
	conditionalDamageRegex  = regexp.MustCompile(`If this Pokémon has at least (\d+) extra {([A-Z])} Energy attached, this attack does (\d+) more damage\.`)
)

// Parse takes the raw text of an effect and attempts to turn it into a structured Effect object.
func Parse(text string) core.Effect {
	// Trim whitespace for easier matching
	text = strings.TrimSpace(text)

	// --- AOE HEAL Effect ---
	if matches := healAllFriendlyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return core.Effect{
				Type:        core.EffectHeal,
				Target:      core.TargetAllFriendly,
				Amount:      amount,
				Description: text,
			}
		}
	}

	// --- HEAL Effect (Single Target) ---
	if matches := healRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return core.Effect{
				Type:        core.EffectHeal,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
			}
		}
	}

	// --- RESTRICTION: Can't attack next turn ---
	if cantAttackNextTurnRegex.MatchString(text) {
		return core.Effect{
			Type:        core.EffectRestrictionCantAttack,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"duration": "next_turn",
			},
		}
	}

	// --- FORCE SWITCH Effect ---
	if forceSwitchRegex.MatchString(text) {
		return core.Effect{
			Type:        core.EffectForceSwitch,
			Target:      core.TargetOpponentActive,
			Description: text,
		}
	}

	// --- RECOIL DAMAGE ---
	if matches := recoilDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return core.Effect{
				Type:        core.EffectRecoilDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
			}
		}
	}

	// --- SEARCH DECK ---
	if matches := searchDeckRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err1 := strconv.Atoi(matches[1])
		pokemonType := matches[2]
		if err1 == nil {
			return core.Effect{
				Type:        core.EffectSearchDeck,
				Target:      core.TargetDeck,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"pokemonType": pokemonType,
					"random":      true,
					"destination": "hand",
				},
			}
		}
	}

	// --- CONDITIONAL DAMAGE ---
	if matches := conditionalDamageRegex.FindStringSubmatch(text); len(matches) > 3 {
		requiredEnergyCount, err1 := strconv.Atoi(matches[1])
		energyType := matches[2]
		damageBonus, err2 := strconv.Atoi(matches[3])

		if err1 == nil && err2 == nil {
			return core.Effect{
				Type:        core.EffectConditionalDamage,
				Amount:      damageBonus,
				Description: text,
				Conditions: map[string]interface{}{
					"requiredExtraEnergyCount": requiredEnergyCount,
					"requiredEnergyType":       energyType,
				},
			}
		}
	}

	// --- Fallback for unknown effects ---
	// If no other rule matches, we return an UNKNOWN effect type. This allows us
	// to see which effects we still need to implement parsing for.
	return core.Effect{
		Type:        core.EffectUnknown,
		Description: text,
	}
}
