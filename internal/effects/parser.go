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
	healRegex                = regexp.MustCompile(`Heal (\d+) damage from this Pokémon\.`)
	cantAttackNextTurnRegex  = regexp.MustCompile(`During your next turn, this Pokémon can't attack\.`)
	forceSwitchRegex         = regexp.MustCompile(`Switch out your opponent’s? Active Pokémon to the Bench\.`)
	searchDeckRegex          = regexp.MustCompile(`Put (\d+) random {([A-Z])} Pokémon from your deck into your hand\.`)
	recoilDamageRegex        = regexp.MustCompile(`This Pokémon also does (\d+) damage to itself\.`)
	healAllFriendlyRegex     = regexp.MustCompile(`(?i)heal (\d+) damage from each of your Pokémon\.`)
	conditionalDamageRegex   = regexp.MustCompile(`If this Pokémon has at least (\d+) extra {([A-Z])} Energy attached, this attack does (\d+) more damage\.`)
	applyStatusOpponentRegex = regexp.MustCompile(`Your opponent's Active Pokémon is now (.*)\.`)
	applyStatusSelfRegex     = regexp.MustCompile(`This Pokémon is now (Poisoned|Asleep|Burned|Confused|Paralyzed)\.`)
	attachEnergyRegex        = regexp.MustCompile(`Take (?:a|1) {([A-Z])} Energy from your Energy Zone and attach it to this Pokémon\.`)
	searchEvolutionRegex     = regexp.MustCompile(`Put a random card that evolves from (\w+) from your deck into your hand\.`)
	//TODO: This regex for Komala is very specific, we can make it more generic later if needed
	triggeredSleepRegex                             = regexp.MustCompile(`whenever you attach an Energy from your Energy Zone to it, it is now Asleep\.`)
	conditionalDamageHPRatioRegex                   = regexp.MustCompile(`If your opponent's Active Pokémon has more remaining HP than this Pokémon, this attack does (\d+) more damage\.`)
	passiveDamageCheckupRegex                       = regexp.MustCompile(`During Pokémon Checkup, if this Pokémon is in the Active Spot, do (\d+) damage to your opponent's Active Pokémon\.`)
	scalingDamageRetreatCostRegex                   = regexp.MustCompile(`This attack does (\d+) more damage for each Energy in your opponent's Active Pokémon's Retreat Cost\.`)
	coinFlipTailsFailsRegex                         = regexp.MustCompile(`Flip a coin\. If tails, this attack does nothing\.`)
	conditionalDamageEvolvedTurnRegex               = regexp.MustCompile(`If this Pokémon evolved during this turn, this attack does (\d+) more damage\.`)
	passiveRetreatCostLatiasRegex                   = regexp.MustCompile(`If you have Latias in play, this Pokémon has no Retreat Cost\.`)
	discardAllEnergyRegex                           = regexp.MustCompile(`Discard all Energy from this Pokémon\.`)
	discardTypedEnergyRegex                         = regexp.MustCompile(`Discard (\d+) {([A-Z])} Energy from this Pokémon\.`)
	moveAllTypedEnergyRegex                         = regexp.MustCompile(`you may move all {([A-Z])} Energy from each of your Pokémon to this Pokémon\.`)
	reduceIncomingDamageRegex                       = regexp.MustCompile(`attacks used by the Defending Pokémon do −(\d+) damage\.`)
	discardFromHandRegex                            = regexp.MustCompile(`Discard a random (.*?) from your opponent's hand\.`)
	cantRetreatRegex                                = regexp.MustCompile(`the Defending Pokémon can't retreat\.`)
	multiHitRandomRegex                             = regexp.MustCompile(`1 of your opponent's Pokémon is chosen at random (\d+) times. For each time a Pokémon was chosen, do (\d+) damage to it\.`)
	statusOnCoinFlipRegex                           = regexp.MustCompile(`If heads, your opponent's Active Pokémon is now (.*)\.`)
	scalingDamageSelfDamageRegex                    = regexp.MustCompile(`This attack does more damage equal to the damage this Pokémon has on it\.`)
	splashDamageBenchedFriendlyRegex                = regexp.MustCompile(`This attack also does (\d+) damage to 1 of your Benched Pokémon\.`)
	conditionalDamageExistingDamageRegex            = regexp.MustCompile(`If your opponent's Active Pokémon has damage on it, this attack does (\d+) more damage\.`)
	conditionalDamageToolAttachedRegex              = regexp.MustCompile(`If your opponent's Active Pokémon has a Pokémon Tool attached, this attack does (\d+) more damage\.`)
	gutsAbilityRegex                                = regexp.MustCompile(`If this Pokémon would be Knocked Out by damage from an attack, flip a coin\. If heads, this Pokémon is not Knocked Out, and its remaining HP becomes (\d+)\.`)
	searchDeckByNameRegex                           = regexp.MustCompile(`Put 1 random (.*?) from your deck onto your Bench\.`)
	snipeAndDiscardTypedRegex                       = regexp.MustCompile(`Discard all {([A-Z])} Energy from this Pokémon\. This attack does (\d+) damage to 1 of your opponent's Pokémon\.`)
	discardRandomEnergySelfRegex                    = regexp.MustCompile(`Discard a random Energy from this Pokémon\.`)
	healOnEvolveRegex                               = regexp.MustCompile(`when you play this Pokémon from your hand to evolve .*?, you may heal (\d+) damage from 1 of your {([A-Z])} Pokémon\.`)
	forceSwitchOncePerTurnRegex                     = regexp.MustCompile(`Once during your turn, you may switch out your opponent's Active Pokémon to the Bench\.`)
	reduceDamageNextTurnSelfRegex                   = regexp.MustCompile(`During your opponent's next turn, this Pokémon takes −(\d+) damage from attacks\.`)
	copyAttackRegex                                 = regexp.MustCompile(`Choose 1 of your opponent's Active Pokémon's attacks and use it as this attack\.`)
	passiveDamageReductionTypedRegex                = regexp.MustCompile(`This Pokémon takes −(\d+) damage from attacks from (.*?) Pokémon\.`)
	drawAtEndOfTurnRegex                            = regexp.MustCompile(`At the end of your turn, if this Pokémon is in the Active Spot, draw a card\.`)
	applyStatusOncePerTurnRegex                     = regexp.MustCompile(`if this Pokémon is in the Active Spot, you may make your opponent's Active Pokémon (Poisoned)\.`)
	discardRandomEnergyGlobalRegex                  = regexp.MustCompile(`Discard a random Energy from among the Energy attached to all Pokémon`)
	switchSelfRegex                                 = regexp.MustCompile(`Switch this Pokémon with 1 of your Benched Pokémon\.`)
	snipeDamageRegex                                = regexp.MustCompile(`This attack does (\d+) damage to 1 of your opponent's Pokémon\.`)
	healOncePerTurnRegex                            = regexp.MustCompile(`Once during your turn, .*? you may heal (\d+) damage from 1 of your Pokémon\.`)
	attachEnergyEndsTurnRegex                       = regexp.MustCompile(`take a {([A-Z])} Energy .*? attach it to this Pokémon\. If you use this Ability, your turn ends\.`)
	discardEnergyOnEvolveRegex                      = regexp.MustCompile(`when you play this Pokémon from your hand to evolve .*? you may discard a random Energy from your opponent's Active Pokémon\.`)
	scalingDamageSelfEnergyRegex                    = regexp.MustCompile(`This attack does (\d+) more damage for each {([A-Z])} Energy attached to this Pokémon\.`)
	drawCardRegex                                   = regexp.MustCompile(`Draw a card\.`)
	attachEnergyMultiBenchedRegex                   = regexp.MustCompile(`Choose (\d+) of your Benched Pokémon\. For each of those Pokémon, take a {([A-Z])} Energy .*? attach it to that Pokémon\.`)
	conditionalDamageSupporterRegex                 = regexp.MustCompile(`If you played a Supporter card from your hand during this turn, this attack does (\d+) more damage\.`)
	drawOnEvolveRegex                               = regexp.MustCompile(`when you play this Pokémon from your hand to evolve .*?, you may draw (\d+) cards\.`)
	discardSingleTypedEnergyRegex                   = regexp.MustCompile(`Discard a {([A-Z])} Energy from this Pokémon\.`)
	shuffleIntoDeckOnCoinFlipRegex                  = regexp.MustCompile(`If heads, your opponent shuffles their Active Pokémon into their deck\.`)
	scalingDamageOpponentEnergyRegex                = regexp.MustCompile(`This attack does (\d+) more damage for each Energy attached to your opponent's Active Pokémon\.`)
	preventAllDamageOnCoinFlipRegex                 = regexp.MustCompile(`If heads, during your opponent's next turn, prevent all damage from—and effects of—attacks done to this Pokémon\.`)
	increaseOpponentCostsRegex                      = regexp.MustCompile(`attacks used by the Defending Pokémon cost (\d+) {([A-Z])} more, and its Retreat Cost is (\d+) {([A-Z])} more\.`)
	snipeDamagedPokemonRegex                        = regexp.MustCompile(`This attack does (\d+) damage to 1 of your opponent's Pokémon that have damage on them\.`)
	switchSelfFromBenchRegex                        = regexp.MustCompile(`if this Pokémon is on your Bench, you may switch it with your Active Pokémon\.`)
	conditionalDamageSelfToolRegex                  = regexp.MustCompile(`If this Pokémon has a Pokémon Tool attached, this attack does (\d+) more damage\.`)
	attachMultipleEnergyRegex                       = regexp.MustCompile(`Take (\d+) {([A-Z])} Energy from your Energy Zone and attach it to (?:this Pokémon|1 of your Benched Pokémon)\.`)
	passiveCostReductionInPlayRegex                 = regexp.MustCompile(`If you have (.*?) in play, attacks used by this Pokémon cost (\d+) less {([A-Z])} Energy\.`)
	conditionalDamageOpponentAbilityRegex           = regexp.MustCompile(`If your opponent's Active Pokémon has an Ability, this attack does (\d+) more damage\.`)
	scalingSnipeOpponentEnergyRegex                 = regexp.MustCompile(`This attack does (\d+) damage to 1 of your opponent's Pokémon for each Energy attached to that Pokémon\.`)
	snipeOncePerTurnRegex                           = regexp.MustCompile(`Once during your turn, you may do (\d+) damage to 1 of your opponent's Pokémon\.`)
	scalingDamageBenchedNameRegex                   = regexp.MustCompile(`This attack does (\d+) more damage for each of your Benched (\w+)\.`)
	damageBenchedOpponentAllRegex                   = regexp.MustCompile(`This attack also does (\d+) damage to each of your opponent's Benched Pokémon\.`)
	lifestealRegex                                  = regexp.MustCompile(`Heal from this Pokémon the same amount of damage you did to your opponent's Active Pokémon\.`)
	passiveEnergyValueRegex                         = regexp.MustCompile(`Each {([A-Z])} Energy attached to your {([A-Z])} Pokémon provides 2 {([A-Z])} Energy\.`)
	scalingDamageMultiCoinFlipRegex                 = regexp.MustCompile(`Flip (\d+) coins\. This attack does (\d+) damage for each heads\.`)
	conditionalDamageSwitchInRegex                  = regexp.MustCompile(`If this Pokémon moved from your Bench to the Active Spot this turn, this attack does (\d+) more damage\.`)
	conditionalDamageCoinFlipScalingRegex           = regexp.MustCompile(`(?i)Flip a coin(?: until you get tails)?\. If heads, this attack does (\d+) more damage(?: for each heads)?\.`)
	passiveImmunityRegex                            = regexp.MustCompile(`This Pokémon can't be affected by any Special Conditions\.`)
	scalingDamageBenchedCountRegex                  = regexp.MustCompile(`This attack does (\d+) more damage for each of your Benched Pokémon\.`)
	reactiveDamageRegex                             = regexp.MustCompile(`if this Pokémon is damaged by an attack, do (\d+) damage to the Attacking Pokémon\.`)
	attachEnergyToActiveTypedRegex                  = regexp.MustCompile(`take a {([A-Z])} Energy from your Energy Zone and attach it to the {([A-Z])} Pokémon in the Active Spot\.`)
	buffNextTurnRegex                               = regexp.MustCompile(`(?i)During your next turn, this Pokémon's ([\w\s]+) attack does \+(\d+) damage\.`)
	attachEnergyToBenchedRegex                      = regexp.MustCompile(`Take a {([A-Z])} Energy from your Energy Zone and attach it to 1 of your Benched\s+Pokémon\.`)
	scalingDamageBenchedTypeRegex                   = regexp.MustCompile(`This attack does (\d+) more damage for each Evolution Pokémon on your Bench\.`)
	passiveRestrictionRegex                         = regexp.MustCompile(`your opponent can't use any (.*?) cards from their hand\.`)
	modifyEnergyRegex                               = regexp.MustCompile(`Change the type of a random Energy attached to your opponent's Active Pokémon to 1 of the following at random: (.*?)\.`)
	triggeredDamageOnEnergyAttachRegex              = regexp.MustCompile(`Whenever you attach a {([A-Z])} Energy from your Energy Zone to this Pokémon, do (\d+) damage to your opponent's Active Pokémon\.`)
	switchSelfTypedRegex                            = regexp.MustCompile(`Switch this Pokémon with 1 of your Benched {([A-Z])} Pokémon\.`)
	conditionalDamageOnKORegex                      = regexp.MustCompile(`If any of your Pokémon were Knocked Out by damage from an attack during your opponent's last turn, this attack does (\d+) more damage\.`)
	splashDamageBenchedFriendlyAllRegex             = regexp.MustCompile(`This attack also does (\d+) damage to each of your Benched Pokémon\.`)
	scalingDamageBenchedNamesRegex                  = regexp.MustCompile(`This attack does (\d+) more damage for each of your Benched (.*?)\.`)
	discardMultipleTypedEnergyRegex                 = regexp.MustCompile(`Discard a {([A-Z])}, {([A-Z])}, and {([A-Z])} Energy from this Pokémon\.`)
	attachEnergyToBenchedStageRegex                 = regexp.MustCompile(`Take a {([A-Z])} Energy from your Energy Zone and attach it to 1 of your Benched (Basic) Pokémon\.`)
	passiveDamageBuffInPlayRegex                    = regexp.MustCompile(`If you have (.*?) in play, attacks used by this Pokémon do \+(\d+) damage to your opponent's Active Pokémon\.`)
	passiveDamagePreventionCoinFlipRegex            = regexp.MustCompile(`If any damage is done to this Pokémon by attacks, flip a coin\. If heads, prevent that damage\.`)
	healAtEndOfTurnRegex                            = regexp.MustCompile(`At the end of your turn, if this Pokémon is in the Active Spot, heal (\d+) damage from it\.`)
	passiveRetreatCostEnergyRegex                   = regexp.MustCompile(`If this Pokémon has any Energy attached, it has no Retreat Cost\.`)
	damageAllOpponentRegex                          = regexp.MustCompile(`This attack does (\d+) damage to each of your opponent's Pokémon\.`)
	passiveRetreatCostReductionBenchRegex           = regexp.MustCompile(`As long as this Pokémon is on your Bench, your Active (.*?) Pokémon's Retreat Cost is (\d+) less\.`)
	damageBenchedConditionalEnergyRegex             = regexp.MustCompile(`This attack also does (\d+) damage to each of your opponent's Benched Pokémon that has any Energy attached\.`)
	scalingDamagePerEnergyCoinFlipRegex             = regexp.MustCompile(`Flip a coin for each Energy attached to this Pokémon\. This attack does (\d+) damage for each heads\.`)
	alternateAttackCostRegex                        = regexp.MustCompile(`If this Pokémon has damage on it, this attack can be used for (\d+) {([A-Z])} Energy\.`)
	discardAllTypedEnergyRegex                      = regexp.MustCompile(`Discard all {([A-Z])} Energy from this Pokémon\.`)
	conditionalDamageOpponentStatusRegex            = regexp.MustCompile(`If your opponent's Active Pokémon is (Burned|Poisoned), this attack does (\d+) more damage\.`)
	preventAllDamageSimpleRegex                     = regexp.MustCompile(`If heads, during your opponent's next turn, prevent all damage done to this Pokémon by attacks\.`)
	scalingDamageBenchedTypeCountRegex              = regexp.MustCompile(`This attack does (\d+) damage for each of your Benched {([A-Z])} Pokémon\.`)
	passiveDamageReductionSimpleRegex               = regexp.MustCompile(`This Pokémon takes −(\d+) damage from attacks\.`)
	conditionalDamageNoDamageRegex                  = regexp.MustCompile(`If this Pokémon has no damage on it, this attack does (\d+) more damage\.`)
	passiveRetreatCostForOtherRegex                 = regexp.MustCompile(`Your Active (\w+) has no Retreat Cost\.`)
	recoilDamageOnCoinFlipRegex                     = regexp.MustCompile(`If tails, this Pokémon also does (\d+) damage to itself\.`)
	passiveImmunitySingleStatusRegex                = regexp.MustCompile(`This Pokémon can't be (Asleep)\.`)
	discardDeckRegex                                = regexp.MustCompile(`Discard the top card of your opponent's deck\.`)
	splashDamageSingleBenchedOpponentRegex          = regexp.MustCompile(`This attack (?:also )?does (\d+) damage to 1 of your opponent's Benched Pokémon\.`)
	scalingDamageAllBenchedRegex                    = regexp.MustCompile(`This attack does (\d+) damage for each Benched Pokémon \(both yours and your opponent's\)\.`)
	multiHitRandomGlobalRegex                       = regexp.MustCompile(`1 other Pokémon \(either yours or your opponent's\) is chosen at random (\d+) times. For each time a Pokémon was chosen, do (\d+) damage to it\.`)
	passiveSpecialConditionImmunityTypedEnergyRegex = regexp.MustCompile(`Each of your Pokémon that has any {([A-Z])} Energy attached recovers from all Special Conditions and can't be affected by any Special Conditions\.`)
	healAllFriendlyTypedRegex                       = regexp.MustCompile(`heal (\d+) damage from each of your {([A-Z])} Pokémon\.`)
	conditionalDamageOpponentPropertyRegex          = regexp.MustCompile(`(?i)If your opponent's Active Pokémon is (a {?[A-Z]}? Pokémon|an Evolution Pokémon|a Pokémon {?ex}?), this attack does (\d+) more damage\.`)
	scalingDamagePerTypedEnergyCoinFlipRegex        = regexp.MustCompile(`Flip a coin for each {([A-Z])} Energy attached to this Pokémon\. This attack does (\d+) damage for each heads\.`)
	healBenchedRegex                                = regexp.MustCompile(`Heal (\d+) damage from 1 of your Benched Pokémon\.`)
	conditionalDamageDifferentEnergyRegex           = regexp.MustCompile(`If this Pokémon has (\d+) or more different types of Energy attached, this attack does (\d+) more damage\.`)
	damageEqualsSelfDamageRegex                     = regexp.MustCompile(`This attack does damage to your opponent's Active Pokémon equal to the damage this Pokémon has on it\.`)
	setHPOnCoinFlipRegex                            = regexp.MustCompile(`If heads, your opponent's Active Pokémon's remaining HP is now (\d+)\.`)
	shuffleFromHandMultiCoinFlipRegex               = regexp.MustCompile(`Flip (\d+) coins\. For each heads, a card is chosen at random from your opponent's hand\. Your opponent reveals that card and shuffles it into their deck\.`)
	passiveReactiveDamageActiveRegex                = regexp.MustCompile(`If this Pokémon is in the Active Spot and is damaged by an attack .*?, do (\d+) damage to the Attacking Pokémon\.`)
	shuffleFromHandRevealRegex                      = regexp.MustCompile(`Your opponent reveals their hand\. Choose a card you find there and shuffle it into your opponent's deck\.`)
	restrictionCantUseAttackRegex                   = regexp.MustCompile(`During your next turn, this Pokémon can't use (.*?)\.`)
	scalingDamageUntilTailsRegex                    = regexp.MustCompile(`Flip a coin until you get tails\. This attack does (\d+) damage for each heads\.`)
	forceSwitchDamagedRegex                         = regexp.MustCompile(`you may switch in 1 of your opponent's Benched Pokémon that has damage on it to the Active Spot\.`)
	conditionalDamageDoubleHeadsRegex               = regexp.MustCompile(`Flip 2 coins\. If both of them are heads, this attack does (\d+) more damage\.`)
	scalingDamagePerPokemonInPlayRegex              = regexp.MustCompile(`Flip a coin for each Pokémon you have in play\. This attack does (\d+) damage for each heads\.`)
	modifyNextEnergyRegex                           = regexp.MustCompile(`Change the type of the next Energy that will be generated for your opponent to 1 of the following at random: (.*?)\.`)
	passiveOpponentDamageReductionRegex             = regexp.MustCompile(`As long as this Pokémon is in the Active Spot, attacks used by your opponent's Active Pokémon do −(\d+) damage\.`)
	conditionalDamageIfDamagedLastTurnRegex         = regexp.MustCompile(`If this Pokémon was damaged by an attack during your opponent's last turn .*?, this attack does (\d+) more damage\.`)
	lookAtTopCardRegex                              = regexp.MustCompile(`you may look at the top card of your deck\.`)
	passiveBuffStatusDamageRegex                    = regexp.MustCompile(`Your opponent's Active Pokémon takes \+(\d+) damage from being (Poisoned)\.`)
	attachEnergyOncePerTurnRegex                    = regexp.MustCompile(`Once during your turn, you may take a {([A-Z])} Energy from your Energy Zone and attach it to this Pokémon\.`)
	searchDeckRandomPokemonRegex                    = regexp.MustCompile(`Put a random Pokémon from your deck into your hand\.`)
	conditionalDamageSelfHasDamageRegex             = regexp.MustCompile(`If this Pokémon has damage on it, this attack does (\d+) more damage\.`)
	conditionalDamageOnAttackHistoryRegex           = regexp.MustCompile(`If 1 of your Pokémon used (.*?) during your last turn, this attack does (\d+) more damage\.`)
	discardEnergyUntilTailsRegex                    = regexp.MustCompile(`Flip a coin until you get tails\. For each heads, discard a random Energy from your opponent's Active Pokémon\.`)
	delayedDamageRegex                              = regexp.MustCompile(`At the end of your opponent's next turn, do (\d+) damage to the Defending Pokémon\.`)
	discardEnergyOpponentOnCoinFlipRegex            = regexp.MustCompile(`(?i)flip a coin\. If heads, discard a random Energy from your opponent's Active Pokémon\.`)
	passiveReactiveDamageOnKORegex                  = regexp.MustCompile(`If this Pokémon is .*? Knocked Out by damage from an attack .*?, do (\d+) damage to the Attacking Pokémon\.`)
	restrictOpponentHandNextTurnRegex               = regexp.MustCompile(`Your opponent can't use any (.*?) cards from their hand during their next turn\.`)
	passiveZeroRetreatForActiveRegex                = regexp.MustCompile(`Your Active Pokémon has no Retreat Cost\.`)
	scalingDamageAttackHistoryRegex                 = regexp.MustCompile(`This attack does (\d+) damage for each time your Pokémon used (.*?) during this game\.`)
	discardDeckBothPlayersRegex                     = regexp.MustCompile(`Discard the top (\d+) cards of each player's deck\.`)
	passiveDamageReductionOnCoinFlipRegex           = regexp.MustCompile(`If any damage is done to this Pokémon by attacks, flip a coin\. If heads, this Pokémon takes −(\d+) damage from that attack\.`)
	moveEnergyBenchedToActiveRegex                  = regexp.MustCompile(`move all {([A-Z])} Energy from 1 of your Benched {([A-Z])} Pokémon to your Active Pokémon\.`)
	restrictOpponentHandItemRegex                   = regexp.MustCompile(`During your opponent's next turn, they can't play any Item cards from their hand\.`)
	passiveDamageBuffEvolvesFromRegex               = regexp.MustCompile(`attacks used by your Pokémon that evolve from (\w+) do \+(\d+) damage to your opponent's Active Pokémon\.`)
	forceSwitchBenchedBasicRegex                    = regexp.MustCompile(`switch in 1 of your opponent's Benched Basic Pokémon to the Active Spot\.`)
	applyStatusBothActiveRegex                      = regexp.MustCompile(`Both Active Pokémon are now (Asleep)\.`)
	applyRestrictionCantAttackRegex                 = regexp.MustCompile(`the Defending Pokémon can't attack\.`)
	passiveDamageReductionInPlayRegex               = regexp.MustCompile(`If you have (.*?) in play, this Pokémon takes −(\d+) damage from attacks\.`)
	passiveDamageBuffTypedPokemonRegex              = regexp.MustCompile(`Attacks used by your {([A-Z])} Pokémon do \+(\d+) damage to your opponent's Active Pokémon\.`)
	drawWithDiscardCostRegex                        = regexp.MustCompile(`You must discard a card from your hand in order to use this Ability\. Once during your turn, you may draw a card\.`)
	knockoutOnCoinFlipRegex                         = regexp.MustCompile(`If both of them are heads, your opponent's Active Pokémon is Knocked Out\.`)
	passiveRestrictionEvolveRegex                   = regexp.MustCompile(`Your opponent can't play any Pokémon from their hand to evolve their Active Pokémon\.`)
	scalingDamageBenchedBaseRegex                   = regexp.MustCompile(`This attack does (\d+) damage for each of your Benched Pokémon\.`)
	scalingDamageAllOpponentEnergyRegex             = regexp.MustCompile(`This attack does (\d+) damage for each Energy attached to all of your opponent's Pokémon\.`)
	moveEnergyOnKORegex                             = regexp.MustCompile(`If this Pokémon is .*? Knocked Out .*?, move all {([A-Z])} Energy from this Pokémon to 1 of your Benched Pokémon\.`)
	recoilDamageOnKORegex                           = regexp.MustCompile(`If your opponent's Pokémon is Knocked Out by damage from this attack, this Pokémon also does (\d+) damage to itself\.`)
	revealHandRegex                                 = regexp.MustCompile(`Your opponent reveals their hand\.`)
	moveDamageRegex                                 = regexp.MustCompile(`As often as you like during your turn, you may choose 1 of your Pokémon that has damage on it, and move all of its damage to this Pokémon\.`)
	discardAllToolsRegex                            = regexp.MustCompile(`(Before doing damage, discard|Discard) all Pokémon Tools from your opponent's Active Pokémon\.`)
	passivePreventionFromEXRegex                    = regexp.MustCompile(`Prevent all damage done to this Pokémon by attacks from your opponent's Pokémon ex\.`)
	discardRandomEnergySelfMultipleRegex            = regexp.MustCompile(`Discard (\d+) random Energy from this Pokémon\.`)
	splashDamageAnyFriendlyRegex                    = regexp.MustCompile(`This attack also does (\d+) damage to 1 of your Pokémon\.`)
	conditionalRestrictionOnStageRegex              = regexp.MustCompile(`If the Defending Pokémon is a (Basic) Pokémon, it can't attack during your opponent's next turn\.`)
	attachEnergyScaledByCoinFlipsRegex              = regexp.MustCompile(`Flip (\d+) coins\. Take an amount of {([A-Z])} Energy .*? equal to the number of heads and attach it to your Benched {([A-Z])} Pokémon in any way you like\.`)
	attachMultipleSpecificEnergyRegex               = regexp.MustCompile(`Take a {([A-Z])}, {([A-Z])}, and {([A-Z])} Energy .*? and attach them to your Benched Basic Pokémon in any way you like\.`)
	conditionalDamageOpponentHasStatusRegex         = regexp.MustCompile(`If your opponent's Active Pokémon is affected by a Special Condition, this attack does (\d+) more damage\.`)
	applyAttackFailureChanceRegex                   = regexp.MustCompile(`if the Defending Pokémon tries to use an attack, your opponent flips a coin\. If tails, that attack doesn't happen\.`)
	discardFromHandOnCoinFlipRegex                  = regexp.MustCompile(`Flip a coin\. If heads, discard a random card from your opponent's hand\.`)
	moveEnergyAsOftenAsYouLikeRegex                 = regexp.MustCompile(`As often as you like during your turn, you may move a {([A-Z])} Energy from 1 of your Benched {([A-Z])} Pokémon to your Active {([A-Z])} Pokémon\.`)
	conditionalRestrictionOnCoinFlipRegex           = regexp.MustCompile(`If heads, the Defending Pokémon can't attack during your opponent's next turn\.`)
	switchSelfSubtypeRegex                          = regexp.MustCompile(`you may switch your Active (.*?) with 1 of your Benched (.*?)\.`)
	attachEnergyFromDiscardWithRecoilRegex          = regexp.MustCompile(`you may attach a {([A-Z])} Energy from your discard pile to this Pokémon\. If you do, do (\d+) damage to this Pokémon\.`)
	passiveEvolveToAnyRegex                         = regexp.MustCompile(`This Pokémon can evolve into any Pokémon that evolves from Eevee`)
	passiveCostReductionWithToolRegex               = regexp.MustCompile(`If this Pokémon has a Pokémon Tool attached, attacks used by this Pokémon cost (\d+) less {([A-Z])} Energy\.`)
	copyAttackWithEnergyCheckRegex                  = regexp.MustCompile(`Choose 1 of your opponent’s Pokémon’s attacks and use it as this attack\. If this Pokémon doesn’t have the necessary Energy to use that attack, this attack does nothing\.`)
	passiveIncreaseOpponentCostRegex                = regexp.MustCompile(`attacks used by your opponent's Active Pokémon cost (\d+) {([A-Z])} more\.`)
	attachEnergyToTypedBenchedRegex                 = regexp.MustCompile(`Take a {([A-Z])} Energy from your Energy Zone and attach it to 1 of your Benched {([A-Z])} Pokémon\.`)
	conditionalDamageOpponentStageRegex             = regexp.MustCompile(`If your opponent's Active Pokémon is a (Basic) Pokémon, this attack does (\d+) more damage\.`)
	lookAtEitherPlayerDeckRegex                     = regexp.MustCompile(`you may choose either player\. Look at the top card of that player's deck\.`)
	persistentAttackFailureRegex                    = regexp.MustCompile(`If the Defending Pokémon tries to use an attack, your opponent flips a coin\. If tails, that attack doesn't happen\. This effect lasts until the Defending Pokémon leaves the Active Spot`)
	healOnCoinFlipRegex                             = regexp.MustCompile(`Flip a coin\. If heads, heal (\d+) damage from this Pokémon\.`)
	discardOwnDeckAmountRegex                       = regexp.MustCompile(`Discard the top (\d+) cards of your deck\.`)
	discardDeckWithConditionalDamageRegex           = regexp.MustCompile(`Discard the top card of your deck\. If that card is a {([A-Z])} Pokémon, this attack does (\d+) more damage\.`)
	conditionalDamageOnBenchedNameRegex             = regexp.MustCompile(`If (\w+) is on your Bench, this attack does (\d+) more damage\.`)
	damageHalveHPRoundedDownRegex                   = regexp.MustCompile(`Halve your opponent's Active Pokémon's remaining HP, rounded down\.`)
	passiveGlobalDamageReductionUnownRegex          = regexp.MustCompile(`All of your Pokémon take −(\d+) damage from attacks from your opponent's Pokémon\.`)
	applyRandomStatusRegex                          = regexp.MustCompile(`1 Special Condition from among (.*?) is chosen at random, and your opponent's Active Pokémon is now affected by that Special Condition\.`)
	scalingDamageDiscardToolRegex                   = regexp.MustCompile(`Discard up to (\d+) Pokémon Tool cards from your hand\. This attack does (\d+) damage for each card you discarded in this way\.`)
	searchDeckToolRegex                             = regexp.MustCompile(`you may put a random Pokémon Tool card from your deck into your hand\.`)
	attachEnergyAtEndOfFirstTurnRegex               = regexp.MustCompile(`At the end of your first turn, take a {([A-Z])} Energy from your Energy Zone and attach it to this Pokémon\.`)
	evolveOnEnergyAttachRegex                       = regexp.MustCompile(`Whenever you attach an Energy from your Energy Zone to this Pokémon, put a random card from your deck that evolves from this Pokémon onto this Pokémon to evolve it\.`)
	healActiveOncePerTurnRegex                      = regexp.MustCompile(`Once during your turn, you may heal (\d+) damage from your Active Pokémon\.`)
	attachEnergyToAnyTypedFriendlyRegex             = regexp.MustCompile(`if this Pokémon is in the Active Spot, you may take a {([A-Z])} Energy from your Energy Zone and attach it to 1 of your {([A-Z])} Pokémon\.`)
	passiveZeroRetreatInPlayRegex                   = regexp.MustCompile(`If you have (.*?) in play, this Pokémon has no Retreat Cost\.`)
	scalingDamageSelfEnergyAllTypesRegex            = regexp.MustCompile(`This attack does (\d+) more damage for each Energy attached to this Pokémon\.`)
	returnToHandOnCoinFlipRegex                     = regexp.MustCompile(`If heads, put your opponent's Active Pokémon into their hand\.`)
	discardBenchedForScalingDamageRegex             = regexp.MustCompile(`You may discard any number of your Benched {([A-Z])} Pokémon\. This attack does (\d+) more damage for each Benched Pokémon you discarded in this way\.`)
	passiveGlobalHealBlockRegex                     = regexp.MustCompile(`Pokémon \(both yours and your opponent's\) can't be healed\.`)
	devolveOnConditionRegex                         = regexp.MustCompile(`If your opponent's Active Pokémon is an evolved Pokémon, devolve it by putting the highest Stage Evolution card on it into your opponent's hand\.`)
	shuffleHandAndDrawScaledRegex                   = regexp.MustCompile(`Shuffle your hand into your deck\. Draw a card for each card in your opponent's hand\.`)
	discardEnergyBothActiveRegex                    = regexp.MustCompile(`Discard a random Energy from both Active Pokémon\.`)
	conditionalDamageBenchedDamagedRegex            = regexp.MustCompile(`If any of your Benched Pokémon have damage on them, this attack does (\d+) more damage\.`)
	shuffleFromHandOnCoinFlipRegex                  = regexp.MustCompile(`If heads, your opponent reveals a random card from their hand and shuffles it into their deck\.`)
	drawUntilMatchHandSizeRegex                     = regexp.MustCompile(`Draw cards until you have the same number of cards in your hand as your opponent\.`)
	scalingDamageOpponentBenchedCountRegex          = regexp.MustCompile(`This attack does (\d+) more damage for each of your opponent's Benched Pokémon\.`)
	copyAttackOnCoinFlipRegex                       = regexp.MustCompile(`If heads, choose 1 of your opponent's Active Pokémon's attacks and use it as this attack\.`)
	snipeRandomOpponentRegex                        = regexp.MustCompile(`1 of your opponent's Pokémon is chosen at random\. Do (\d+) damage to it\.`)
	passiveGlobalDamageBuffUnownRegex               = regexp.MustCompile(`Attacks used by your Pokémon do \+(\d+) damage to your opponent's Active Pokémon\.`)
	preventionOnKORegex                             = regexp.MustCompile(`If your opponent's Pokémon is Knocked Out by damage from this Pokémon's attacks, during your opponent's next turn, prevent all damage from—and effects of—attacks done to this Pokémon\.`)
	moveAllEnergyToBenchedRegex                     = regexp.MustCompile(`Move all Energy from this Pokémon to 1 of your Benched Pokémon\.`)
	restrictEnergyAttachmentRegex                   = regexp.MustCompile(`they can't take any Energy from their Energy Zone to attach to their Active Pokémon\.`)
	shuffleFromHandSimpleRegex                      = regexp.MustCompile(`Your opponent reveals a random card from their hand and shuffles it into their deck\.`)
	voluntarySwitchRegex                            = regexp.MustCompile(`You may switch this Pokémon with 1 of your Benched Pokémon\.`)
	revealHandOnBenchPlayRegex                      = regexp.MustCompile(`when you put this Pokémon from your hand onto your Bench, you may have your opponent reveal their hand\.`)
	debuffIncomingDamageRegex                       = regexp.MustCompile(`During your opponent's next turn, this Pokémon takes \+(\d+) damage from attacks\.`)
	restrictionOnTailsRegex                         = regexp.MustCompile(`If tails, during your next turn, this Pokémon can't attack\.`)
	conditionalDamageOpponentNameRegex              = regexp.MustCompile(`If your opponent's Active Pokémon is (\w+), this attack does (\d+) more damage\.`)
	attachEnergyToSpecificPokemonRegex              = regexp.MustCompile(`Take a {([A-Z])} Energy from your Energy Zone and attach it to (\w+) or (\w+)\.`)
	passiveRetreatCostFirstTurnRegex                = regexp.MustCompile(`During your first turn, this Pokémon has no Retreat Cost\.`)
	passiveEffectPreventionRegex                    = regexp.MustCompile(`Prevent all effects of attacks used by your opponent's Pokémon done to this Pokémon\.`)
	applyStatusOncePerTurnAbilityRegex              = regexp.MustCompile(`Once during your turn, you may make your opponent's Active Pokémon (Burned)\.`)
	buffStackingRegex                               = regexp.MustCompile(`Until this Pokémon leaves the Active Spot, this Pokémon's (.*?) attack does \+(\d+) damage\. This effect stacks\.`)
	conditionalDamageIfEnergyAttachedRegex          = regexp.MustCompile(`If this Pokémon has any {([A-Z])} Energy attached, this attack does (\d+) more damage\.`)
	healOnEnergyAttachRegex                         = regexp.MustCompile(`Whenever you attach a {([A-Z])} Energy from your Energy Zone to this Pokémon, heal (\d+) damage from this Pokémon\.`)
	knockoutAttackerOnKORegex                       = regexp.MustCompile(`If this Pokémon is .*? Knocked Out .*?, flip a coin\. If heads, the Attacking Pokémon is Knocked Out\.`)
	damageAbilityInPlayRegex                        = regexp.MustCompile(`if you have (.*?) in play, you may do (\d+) damage to your opponent's Active Pokémon\.`)
	finalConditionalDamageMultiCoinFlipRegex        = regexp.MustCompile(`(?i)Flip (\d+) coins\. This attack does (\d+) more damage for each heads\.`)
	finalPassiveDamageReductionRegex                = regexp.MustCompile(`(?i)This Pokémon takes -(\d+) damage from attacks\.`)
	finalDiscardEnergyOpponentSimpleRegex           = regexp.MustCompile(`(?i)Discard a random Energy from your opponent's Active Pokémon\.`)
	finalConditionalDamageUntilTailsRegex           = regexp.MustCompile(`(?i)Flip a coin until you get tails\. This attack does (\d+) more damage for each heads\.`)
	finalScalingDamageOpponentEnergyBaseRegex       = regexp.MustCompile(`(?i)This attack does (\d+) damage for each Energy attached to your opponent's Active Pokémon\.`)
	finalAttachEnergyToActiveTypedRegex             = regexp.MustCompile(`(?i)you may take 1? {([A-Z])} Energy from your Energy Zone and attach it to the {([A-Z])} Pokémon in the Active Spot\.`)
	finalConditionalDamageOpponentIsEXRegex         = regexp.MustCompile(`(?i)If your opponent’s Active Pokémon is a Pokémon {ex}, this attack does (\d+) more damage\.`)
	finalDiscardSelfEnergyOnTailsRegex              = regexp.MustCompile(`(?i)Flip a coin\. If tails, discard (\d+) random Energy from this Pokémon\.`)
	finalReduceDamageNextTurnRegex                  = regexp.MustCompile(`(?i)During your opponent's next turn, this Pokémon takes -(\d+) damage from attacks\.`)
	finalForceSwitchOnHeadsRegex                    = regexp.MustCompile(`(?i)Flip a coin\. If heads, switch in 1 of your opponent's Benched Pokémon to the Active Spot\.`)
	finalIncreaseOpponentCostNextTurnRegex          = regexp.MustCompile(`(?i)During your opponent's next turn, attacks used by the Defending Pokémon cost (\d+) {([A-Z])} more\.`)
	finalPreventionOnHeadsRegex                     = regexp.MustCompile(`(?i)Flip a coin\. If heads, during your opponent’s next turn, prevent all damage from—and effects of—attacks done to this Pokémon\.`)
	finalSearchDeckGenericPokemonRegex              = regexp.MustCompile(`(?i)Once during your turn, you may put a random Pokémon from your deck into your hand\.`)
	finalForceSwitchSimpleRegex                     = regexp.MustCompile(`(?i)Switch out your opponent’s Active Pokémon to the Bench\.`)
	hariyamaPushOutRegex                            = regexp.MustCompile(`Switch out your opponent's Active Pokémon to the Bench\. \(Your opponent chooses the new Active Pokémon\.\)`)
)

// Parse takes the raw text of an effect and attempts to turn it into a structured Effect object.
func Parse(text string) []core.Effect {
	// Trim whitespace for easier matching
	text = strings.TrimSpace(text)

	// --- AOE HEAL Effect ---
	if matches := healAllFriendlyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Target:      core.TargetAllFriendly,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- HEAL Effect (Single Target) ---
	if matches := healRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- RESTRICTION: Can't attack next turn ---
	if cantAttackNextTurnRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectRestrictionCantAttack,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"duration": "next_turn",
			},
		}}
	}

	// --- FORCE SWITCH Effect ---
	if forceSwitchRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectForceSwitch,
			Target:      core.TargetOpponentActive,
			Description: text,
		}}
	}

	// --- RECOIL DAMAGE ---
	if matches := recoilDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectRecoilDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- SEARCH DECK ---
	if matches := searchDeckRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err1 := strconv.Atoi(matches[1])
		pokemonType := matches[2]
		if err1 == nil {
			return []core.Effect{{
				Type:        core.EffectSearchDeck,
				Target:      core.TargetDeck,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"pokemonType": pokemonType,
					"random":      true,
					"destination": "hand",
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE ---
	if matches := conditionalDamageRegex.FindStringSubmatch(text); len(matches) > 3 {
		requiredEnergyCount, err1 := strconv.Atoi(matches[1])
		energyType := matches[2]
		damageBonus, err2 := strconv.Atoi(matches[3])

		if err1 == nil && err2 == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      damageBonus,
				Description: text,
				Conditions: map[string]interface{}{
					"requiredExtraEnergyCount": requiredEnergyCount,
					"requiredEnergyType":       energyType,
				},
			}}
		}
	}

	// --- APPLY STATUS (Opponent) ---
	if matches := applyStatusOpponentRegex.FindStringSubmatch(text); len(matches) > 1 {
		// Handles "Poisoned", "Poisoned and Burned", etc.
		statusesStr := strings.ReplaceAll(matches[1], " and ", ", ")
		statusSlice := strings.Split(statusesStr, ", ")

		var effects []core.Effect
		for _, status := range statusSlice {
			effects = append(effects, core.Effect{
				Type:        core.EffectApplyStatus,
				Target:      core.TargetOpponentActive,
				Status:      core.StatusCondition(strings.ToUpper(status)),
				Description: text,
			})
		}
		return effects
	}

	// --- APPLY STATUS (Self) ---
	if matches := applyStatusSelfRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectApplyStatus,
			Target:      core.TargetSelf,
			Status:      core.StatusCondition(strings.ToUpper(matches[1])),
			Description: text,
		}}
	}

	// --- ATTACH ENERGY ---
	if matches := attachEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"source":     "EnergyZone",
				"energyType": matches[1],
			},
		}}
	}

	// --- SEARCH DECK (Evolution) ---
	if matches := searchEvolutionRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectSearchDeck,
			Target:      core.TargetDeck,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"evolvesFrom": matches[1],
				"random":      true,
				"destination": "hand",
			},
		}}
	}

	// --- TRIGGERED ABILITY (Komala's Comatose) ---
	if triggeredSleepRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectTriggeredAbility,
			Target:      core.TargetSelf,
			Status:      core.StatusAsleep,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger": "ATTACH_ENERGY_TO_SELF",
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (HP Ratio) ---
	if matches := conditionalDamageHPRatioRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_HP_GREATER",
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Evolved this turn) ---
	if matches := conditionalDamageEvolvedTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "EVOLVED_THIS_TURN",
				},
			}}
		}
	}

	// --- PASSIVE DAMAGE (Checkup) ---
	if matches := passiveDamageCheckupRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveDamage,
				Target:      core.TargetOpponentActive,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"phase":    "CHECKUP",
					"location": "ACTIVE",
				},
			}}
		}
	}

	// --- SCALING DAMAGE (Opponent Retreat Cost) ---
	if matches := scalingDamageRetreatCostRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "OPPONENT_RETREAT_COST",
				},
			}}
		}
	}

	// --- ATTACK MAY FAIL (Coin Flip) ---
	if coinFlipTailsFailsRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectAttackMayFail,
			Description: text,
			Conditions: map[string]interface{}{
				"chance": 0.5,
				"on":     "TAILS",
			},
		}}
	}

	// --- PASSIVE ABILITY (Retreat Cost) ---
	if passiveRetreatCostLatiasRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":           "ZERO_RETREAT_COST",
				"requires_in_play": "Latias",
			},
		}}
	}

	// --- DISCARD ENERGY (All) ---
	if discardAllEnergyRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"amount": "ALL",
			},
		}}
	}

	// --- DISCARD ENERGY (Typed) ---
	if matches := discardTypedEnergyRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDiscardEnergy,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"energyType": matches[2],
				},
			}}
		}
	}

	// --- MOVE ENERGY ---
	if matches := moveAllTypedEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectMoveEnergy,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"source":     "ALL_FRIENDLY",
				"energyType": matches[1],
			},
		}}
	}

	// --- REDUCE INCOMING DAMAGE ---
	if matches := reduceIncomingDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectReduceIncomingDamage,
				Target:      core.TargetOpponentActive,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"duration": "next_turn",
				},
			}}
		}
	}

	// --- DISCARD FROM HAND ---
	if matches := discardFromHandRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectDiscardFromHand,
			Target:      core.TargetOpponentHand,
			Amount:      1, // "a random" implies 1
			Description: text,
			Conditions: map[string]interface{}{
				"card_type": matches[1],
				"random":    true,
			},
		}}
	}

	// --- APPLY RESTRICTION (Can't Retreat) ---
	if cantRetreatRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "CANT_RETREAT",
				"duration":    "next_turn",
			},
		}}
	}

	// --- MULTI-HIT RANDOM DAMAGE ---
	if matches := multiHitRandomRegex.FindStringSubmatch(text); len(matches) > 2 {
		hits, err1 := strconv.Atoi(matches[1])
		damage, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			return []core.Effect{{
				Type:        core.EffectMultiHitRandomDamage,
				Target:      core.TargetOpponentActive, // Target pool is opponent's field
				Amount:      damage,
				Description: text,
				Conditions: map[string]interface{}{
					"hits": hits,
				},
			}}
		}
	}

	// --- APPLY STATUS (On Coin Flip) ---
	if matches := statusOnCoinFlipRegex.FindStringSubmatch(text); len(matches) > 1 {
		// Handles "Poisoned and Paralyzed"
		statuses := strings.Split(matches[1], " and ")
		return []core.Effect{{
			Type:        core.EffectApplyStatus,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
				"statuses":     statuses,
			},
		}}
	}

	// --- SCALING DAMAGE (Self Damage) ---
	if scalingDamageSelfDamageRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectScalingDamage,
			Description: text,
			Conditions: map[string]interface{}{
				"scale_by": "SELF_DAMAGE_COUNTERS",
			},
		}}
	}

	// --- SPLASH DAMAGE (Friendly Bench) ---
	if matches := splashDamageBenchedFriendlyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDamageBenchedFriendly,
				Target:      core.TargetBenchedFriendly,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Opponent has damage) ---
	if matches := conditionalDamageExistingDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_HAS_DAMAGE",
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Opponent has tool) ---
	if matches := conditionalDamageToolAttachedRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_HAS_TOOL",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Guts) ---
	if matches := gutsAbilityRegex.FindStringSubmatch(text); len(matches) > 1 {
		hp, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Target:      core.TargetSelf,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":       "PREVENT_KNOCKOUT",
					"on_coin_flip": "HEADS",
					"remaining_hp": hp,
				},
			}}
		}
	}

	// --- SEARCH DECK (By Name) ---
	if matches := searchDeckByNameRegex.FindStringSubmatch(text); len(matches) > 1 {
		// Handles "Wishiwashi or Wishiwashi ex"
		names := strings.Split(matches[1], " or ")
		return []core.Effect{{
			Type:        core.EffectSearchDeck,
			Target:      core.TargetDeck,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"pokemonNames": names,
				"random":       true,
				"destination":  "bench",
			},
		}}
	}

	// --- SNIPE and DISCARD (Luxray) ---
	if matches := snipeAndDiscardTypedRegex.FindStringSubmatch(text); len(matches) > 2 {
		energyType := matches[1]
		damage, err := strconv.Atoi(matches[2])
		if err == nil {
			discardEffect := core.Effect{
				Type:        core.EffectDiscardEnergy,
				Target:      core.TargetSelf,
				Description: "Discard all {L} Energy from this Pokémon.",
				Conditions: map[string]interface{}{
					"amount":     "ALL",
					"energyType": energyType,
				},
			}
			snipeEffect := core.Effect{
				Type:        core.EffectSnipeDamage,
				Target:      core.TargetBenchedOpponent, // Can target any opponent's Pokémon
				Amount:      damage,
				Description: "This attack does 120 damage to 1 of your opponent's Pokémon.",
			}
			return []core.Effect{discardEffect, snipeEffect}
		}
	}

	// --- DISCARD ENERGY (Random Self) ---
	if discardRandomEnergySelfRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetSelf,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"random": true,
			},
		}}
	}

	// --- HEAL (On Evolve) ---
	if matches := healOnEvolveRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		targetType := matches[2]
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":     "ON_EVOLVE",
					"target_type": targetType,
				},
			}}
		}
	}

	// --- FORCE SWITCH (Once per turn ability) ---
	if forceSwitchOncePerTurnRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectForceSwitch,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger": "ONCE_PER_TURN",
			},
		}}
	}

	// --- REDUCE INCOMING DAMAGE (Self, next turn) ---
	if matches := reduceDamageNextTurnSelfRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectReduceIncomingDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"duration": "opponent_next_turn",
				},
			}}
		}
	}

	// --- COPY ATTACK ---
	if copyAttackRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectCopyAttack,
			Target:      core.TargetOpponentActive,
			Description: text,
		}}
	}

	// --- PASSIVE ABILITY (Damage Reduction by Type) ---
	if matches := passiveDamageReductionTypedRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		// Handles "{R}" or "{R} or {W}"
		typesStr := strings.ReplaceAll(matches[2], "{", "")
		typesStr = strings.ReplaceAll(typesStr, "}", "")
		types := strings.Split(typesStr, " or ")
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Target:      core.TargetSelf,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":     "REDUCE_INCOMING_DAMAGE",
					"amount":     amount,
					"from_types": types,
				},
			}}
		}
	}

	// --- DRAW (End of turn) ---
	if drawAtEndOfTurnRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDraw,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":  "END_OF_TURN",
				"location": "ACTIVE",
			},
		}}
	}

	// --- APPLY STATUS (Once per turn ability) ---
	if matches := applyStatusOncePerTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectApplyStatus,
			Target:      core.TargetOpponentActive,
			Status:      core.StatusCondition(strings.ToUpper(matches[1])),
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":  "ONCE_PER_TURN",
				"location": "ACTIVE",
			},
		}}
	}

	// --- DISCARD ENERGY (Random Global) ---
	if discardRandomEnergyGlobalRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetAllPokemonInPlay,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"random": true,
			},
		}}
	}

	// --- SWITCH SELF ---
	if switchSelfRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectSwitchSelf,
			Target:      core.TargetBenchedFriendly,
			Description: text,
		}}
	}

	// --- SNIPE DAMAGE (Standalone) ---
	if matches := snipeDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectSnipeDamage,
				Target:      core.TargetBenchedOpponent, // User chooses which one
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- HEAL (Once per turn ability) ---
	if matches := healOncePerTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":  "ONCE_PER_TURN",
					"location": "ACTIVE",
				},
			}}
		}
	}

	// --- ATTACH ENERGY (Ends turn) ---
	if matches := attachEnergyEndsTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":    "ONCE_PER_TURN",
				"energyType": matches[1],
				"source":     "EnergyZone",
				"effect":     "ENDS_TURN",
			},
		}}
	}

	// --- DISCARD ENERGY (On Evolve) ---
	if discardEnergyOnEvolveRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetOpponentActive,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger": "ON_EVOLVE",
				"random":  true,
			},
		}}
	}

	// --- SCALING DAMAGE (Self Attached Energy) ---
	if matches := scalingDamageSelfEnergyRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		energyType := matches[2]
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":      "SELF_ATTACHED_ENERGY",
					"scale_by_type": energyType,
				},
			}}
		}
	}

	// --- DRAW (Simple) ---
	if drawCardRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDraw,
			Amount:      1,
			Description: text,
		}}
	}

	// --- ATTACH ENERGY (Multi-Benched) ---
	if matches := attachEnergyMultiBenchedRegex.FindStringSubmatch(text); len(matches) > 2 {
		count, err1 := strconv.Atoi(matches[1])
		energyType := matches[2]
		if err1 == nil {
			return []core.Effect{{
				Type:        core.EffectAttachEnergy,
				Target:      core.TargetBenchedFriendly,
				Amount:      count,
				Description: text,
				Conditions: map[string]interface{}{
					"source":     "EnergyZone",
					"energyType": energyType,
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Played Supporter) ---
	if matches := conditionalDamageSupporterRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "PLAYED_SUPPORTER_THIS_TURN",
				},
			}}
		}
	}

	// --- DRAW (On Evolve) ---
	if matches := drawOnEvolveRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDraw,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "ON_EVOLVE",
				},
			}}
		}
	}

	// --- DISCARD ENERGY (Single Typed) ---
	if matches := discardSingleTypedEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetSelf,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"energyType": matches[1],
			},
		}}
	}

	// --- SHUFFLE INTO DECK (On Coin Flip) ---
	if shuffleIntoDeckOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectShuffleIntoDeck,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
			},
		}}
	}

	// --- SCALING DAMAGE (Opponent Attached Energy) ---
	if matches := scalingDamageOpponentEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "OPPONENT_ATTACHED_ENERGY",
				},
			}}
		}
	}

	// --- APPLY PREVENTION (On Coin Flip) ---
	if preventAllDamageOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyPrevention,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
				"duration":     "opponent_next_turn",
				"prevent":      "ALL_DAMAGE_AND_EFFECTS",
			},
		}}
	}

	// --- APPLY RESTRICTION (Increase opponent costs) ---
	if matches := increaseOpponentCostsRegex.FindStringSubmatch(text); len(matches) > 4 {
		attackCost, err1 := strconv.Atoi(matches[1])
		retreatCost, err2 := strconv.Atoi(matches[3])
		if err1 == nil && err2 == nil {
			attackCostEffect := core.Effect{
				Type:        core.EffectApplyRestriction,
				Target:      core.TargetOpponentActive,
				Description: "Increase opponent's attack cost",
				Conditions: map[string]interface{}{
					"restriction": "INCREASE_ATTACK_COST",
					"amount":      attackCost,
					"energyType":  matches[2],
					"duration":    "opponent_next_turn",
				},
			}
			retreatCostEffect := core.Effect{
				Type:        core.EffectApplyRestriction,
				Target:      core.TargetOpponentActive,
				Description: "Increase opponent's retreat cost",
				Conditions: map[string]interface{}{
					"restriction": "INCREASE_RETREAT_COST",
					"amount":      retreatCost,
					"energyType":  matches[4],
					"duration":    "opponent_next_turn",
				},
			}
			return []core.Effect{attackCostEffect, retreatCostEffect}
		}
	}

	// --- SNIPE DAMAGE (Damaged Pokémon) ---
	if matches := snipeDamagedPokemonRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectSnipeDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"target_condition": "HAS_DAMAGE",
				},
			}}
		}
	}

	// --- SWITCH SELF (From Bench) ---
	if switchSelfFromBenchRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectSwitchSelf,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":  "ONCE_PER_TURN",
				"location": "BENCH",
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Self has tool) ---
	if matches := conditionalDamageSelfToolRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "SELF_HAS_TOOL",
				},
			}}
		}
	}

	// --- ATTACH ENERGY (Multiple) ---
	if matches := attachMultipleEnergyRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		energyType := matches[2]
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectAttachEnergy,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"source":     "EnergyZone",
					"energyType": energyType,
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Cost Reduction) ---
	if matches := passiveCostReductionInPlayRegex.FindStringSubmatch(text); len(matches) > 3 {
		// Handles "Arceus or Arceus ex"
		names := strings.Split(matches[1], " or ")
		amount, err1 := strconv.Atoi(matches[2])
		energyType := matches[3]
		if err1 == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":           "REDUCE_ATTACK_COST",
					"requires_in_play": names,
					"amount":           amount,
					"energyType":       energyType,
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Opponent has Ability) ---
	if matches := conditionalDamageOpponentAbilityRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_HAS_ABILITY",
				},
			}}
		}
	}

	// --- SCALING SNIPE DAMAGE (Opponent Energy) ---
	if matches := scalingSnipeOpponentEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingSnipeDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "TARGET_ATTACHED_ENERGY",
				},
			}}
		}
	}

	// --- SNIPE DAMAGE (Once per turn ability) ---
	if matches := snipeOncePerTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectSnipeDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "ONCE_PER_TURN",
				},
			}}
		}
	}

	// --- SCALING DAMAGE (Benched Pokémon Name) ---
	if matches := scalingDamageBenchedNameRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		name := matches[2]
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":      "BENCHED_POKEMON_NAME",
					"scale_by_name": name,
				},
			}}
		}
	}

	// --- DAMAGE BENCHED OPPONENT (All) ---
	if matches := damageBenchedOpponentAllRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDamageBenchedOpponentAll,
				Target:      core.TargetBenchedOpponentAll,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- LIFESTEAL ---
	if lifestealRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectLifesteal,
			Target:      core.TargetSelf,
			Description: text,
		}}
	}

	// --- PASSIVE ABILITY (Energy Value) ---
	if matches := passiveEnergyValueRegex.FindStringSubmatch(text); len(matches) > 3 {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":       "ENERGY_VALUE_DOUBLED",
				"energy_type":  matches[1],
				"pokemon_type": matches[2],
			},
		}}
	}

	// --- SCALING DAMAGE (Multi-coin flip) ---
	if matches := scalingDamageMultiCoinFlipRegex.FindStringSubmatch(text); len(matches) > 2 {
		flips, err1 := strconv.Atoi(matches[1])
		amount, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":  "COIN_FLIP_HEADS",
					"num_flips": flips,
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Switched in) ---
	if matches := conditionalDamageSwitchInRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "SWITCHED_IN_THIS_TURN",
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Coin Flip Scaling) ---
	if matches := conditionalDamageCoinFlipScalingRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "COIN_FLIP_HEADS",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Immunity) ---
	if passiveImmunityRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect": "IMMUNE_TO_SPECIAL_CONDITIONS",
			},
		}}
	}

	// --- SCALING DAMAGE (Benched Pokémon count) ---
	if matches := scalingDamageBenchedCountRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "BENCHED_POKEMON_COUNT",
				},
			}}
		}
	}

	// --- APPLY REACTIVE DAMAGE ---
	if matches := reactiveDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectApplyReactiveDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"duration": "opponent_next_turn",
				},
			}}
		}
	}

	// --- ATTACH ENERGY (To Active Typed Pokémon) ---
	if matches := attachEnergyToActiveTypedRegex.FindStringSubmatch(text); len(matches) > 2 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":         "ONCE_PER_TURN",
				"source":          "EnergyZone",
				"energyType":      matches[1],
				"target_location": "ACTIVE",
				"target_type":     matches[2],
			},
		}}
	}

	// --- BUFF NEXT TURN ---
	if matches := buffNextTurnRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectBuffNextTurn,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"attack_name": matches[1],
				},
			}}
		}
	}

	// --- ATTACH ENERGY (To Benched) ---
	if matches := attachEnergyToBenchedRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetBenchedFriendly,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"source":     "EnergyZone",
				"energyType": matches[1],
			},
		}}
	}

	// --- SCALING DAMAGE (Benched Pokémon Type) ---
	if matches := scalingDamageBenchedTypeRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":      "BENCHED_POKEMON_TYPE",
					"scale_by_type": "EVOLUTION",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Restriction) ---
	if matches := passiveRestrictionRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":    "RESTRICT_OPPONENT_PLAY",
				"card_type": matches[1],
				"location":  "ACTIVE",
			},
		}}
	}

	// --- MODIFY ENERGY ---
	if matches := modifyEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		typesStr := strings.ReplaceAll(matches[1], "{", "")
		types := strings.Split(typesStr, "}, ")
		types[len(types)-1] = strings.TrimSuffix(types[len(types)-1], "}")

		return []core.Effect{{
			Type:        core.EffectModifyEnergy,
			Target:      core.TargetOpponentActive,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"random_energy":  true,
				"random_type":    true,
				"possible_types": types,
			},
		}}
	}

	// --- TRIGGERED DAMAGE (On Energy Attach) ---
	if matches := triggeredDamageOnEnergyAttachRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveDamage,
				Target:      core.TargetOpponentActive,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":     "ATTACH_ENERGY_TO_SELF",
					"energy_type": matches[1],
				},
			}}
		}
	}

	// --- SWITCH SELF (Typed) ---
	if matches := switchSelfTypedRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectSwitchSelf,
			Target:      core.TargetBenchedFriendly,
			Description: text,
			Conditions: map[string]interface{}{
				"target_type": matches[1],
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Opponent has specific property) ---
	if matches := conditionalDamageOpponentPropertyRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":  "OPPONENT_HAS_PROPERTY",
					"property": strings.ToUpper(matches[1]),
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Friendly KO'd last turn) ---
	if matches := conditionalDamageOnKORegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "FRIENDLY_KO_LAST_TURN",
				},
			}}
		}
	}

	// --- SPLASH DAMAGE (Friendly Bench All) ---
	if matches := splashDamageBenchedFriendlyAllRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDamageBenchedFriendly,
				Target:      core.TargetBenchedFriendly,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"target_all": true,
				},
			}}
		}
	}

	// --- SCALING DAMAGE (Benched Pokémon Names) ---
	if matches := scalingDamageBenchedNamesRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		// Handles "Wishiwashi and Wishiwashi ex"
		names := strings.Split(matches[2], " and ")
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":       "BENCHED_POKEMON_NAMES",
					"scale_by_names": names,
				},
			}}
		}
	}

	// --- DISCARD ENERGY (Multiple specific types) ---
	if matches := discardMultipleTypedEnergyRegex.FindStringSubmatch(text); len(matches) > 3 {
		types := []string{matches[1], matches[2], matches[3]}
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetSelf,
			Amount:      3,
			Description: text,
			Conditions: map[string]interface{}{
				"energyTypes": types,
			},
		}}
	}

	// --- ATTACH ENERGY (To Benched Stage) ---
	if matches := attachEnergyToBenchedStageRegex.FindStringSubmatch(text); len(matches) > 2 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetBenchedFriendly,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"source":       "EnergyZone",
				"energyType":   matches[1],
				"target_stage": matches[2],
			},
		}}
	}

	// --- PASSIVE ABILITY (Damage Buff) ---
	if matches := passiveDamageBuffInPlayRegex.FindStringSubmatch(text); len(matches) > 2 {
		names := strings.Split(matches[1], " or ")
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":           "BUFF_DAMAGE_OUTPUT",
					"amount":           amount,
					"requires_in_play": names,
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Damage Prevention Coin Flip) ---
	if passiveDamagePreventionCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":       "PREVENT_INCOMING_DAMAGE",
				"on_coin_flip": "HEADS",
			},
		}}
	}

	// --- HEAL (End of turn) ---
	if matches := healAtEndOfTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":  "END_OF_TURN",
					"location": "ACTIVE",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Zero Retreat Cost) ---
	if passiveRetreatCostEnergyRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":  "ZERO_RETREAT_COST",
				"trigger": "HAS_ENERGY_ATTACHED",
			},
		}}
	}

	// --- DAMAGE ALL OPPONENT POKEMON ---
	if matches := damageAllOpponentRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDamageAllOpponent,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- PASSIVE ABILITY (Retreat Cost Reduction for others) ---
	if matches := passiveRetreatCostReductionBenchRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":       "REDUCE_RETREAT_COST",
					"location":     "BENCH",
					"target":       "ACTIVE",
					"target_stage": matches[1],
					"amount":       amount,
				},
			}}
		}
	}

	// --- DAMAGE BENCHED OPPONENT (Conditional on Energy) ---
	if matches := damageBenchedConditionalEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDamageBenchedOpponentAll,
				Target:      core.TargetBenchedOpponentAll,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"target_condition": "HAS_ENERGY_ATTACHED",
				},
			}}
		}
	}

	// --- SCALING DAMAGE (Coin flips per energy) ---
	if matches := scalingDamagePerEnergyCoinFlipRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":            "COIN_FLIP_HEADS",
					"num_flips_scales_by": "SELF_ATTACHED_ENERGY",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Alternate Attack Cost) ---
	if matches := alternateAttackCostRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":      "ALTERNATE_ATTACK_COST",
					"trigger":     "SELF_HAS_DAMAGE",
					"cost_amount": amount,
					"cost_type":   matches[2],
				},
			}}
		}
	}

	// --- DISCARD ENERGY (All of a specific type) ---
	if matches := discardAllTypedEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"amount":     "ALL",
				"energyType": matches[1],
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Opponent has status) ---
	if matches := conditionalDamageOpponentStatusRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_HAS_STATUS",
					"status":  strings.ToUpper(matches[1]),
				},
			}}
		}
	}

	// --- APPLY PREVENTION (Damage only) ---
	if preventAllDamageSimpleRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyPrevention,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
				"duration":     "opponent_next_turn",
				"prevent":      "ALL_DAMAGE",
			},
		}}
	}

	// --- SCALING DAMAGE (Benched Pokémon Type Count) ---
	if matches := scalingDamageBenchedTypeCountRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":      "BENCHED_POKEMON_TYPE_COUNT",
					"scale_by_type": matches[2],
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Simple Damage Reduction) ---
	if matches := passiveDamageReductionSimpleRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Target:      core.TargetSelf,
				Description: text,
				Conditions: map[string]interface{}{
					"effect": "REDUCE_INCOMING_DAMAGE",
					"amount": amount,
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Self has no damage) ---
	if matches := conditionalDamageNoDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "SELF_HAS_NO_DAMAGE",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Zero Retreat for another Pokémon) ---
	if matches := passiveRetreatCostForOtherRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":          "ZERO_RETREAT_COST",
				"target_name":     matches[1],
				"target_location": "ACTIVE",
			},
		}}
	}

	// --- RECOIL DAMAGE (On Coin Flip) ---
	if matches := recoilDamageOnCoinFlipRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectRecoilDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"on_coin_flip": "TAILS",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Single Status Immunity) ---
	if matches := passiveImmunitySingleStatusRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect": "IMMUNE_TO_STATUS",
				"status": strings.ToUpper(matches[1]),
			},
		}}
	}

	// --- DISCARD DECK ---
	if discardDeckRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardDeck,
			Target:      core.TargetOpponentActive, // Implied target is opponent
			Amount:      1,
			Description: text,
		}}
	}

	// --- SPLASH DAMAGE (Single Benched Opponent) ---
	if matches := splashDamageSingleBenchedOpponentRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectSnipeDamage,
				Target:      core.TargetBenchedOpponent,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- SCALING DAMAGE (All Benched Pokémon) ---
	if matches := scalingDamageAllBenchedRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "ALL_BENCHED_POKEMON_COUNT",
				},
			}}
		}
	}

	// --- MULTI-HIT RANDOM DAMAGE (Global Pool) ---
	if matches := multiHitRandomGlobalRegex.FindStringSubmatch(text); len(matches) > 2 {
		hits, err1 := strconv.Atoi(matches[1])
		damage, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			return []core.Effect{{
				Type:        core.EffectMultiHitRandomDamage,
				Amount:      damage,
				Description: text,
				Conditions: map[string]interface{}{
					"hits":        hits,
					"target_pool": "GLOBAL_OTHER",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Special Condition Immunity via Energy) ---
	if matches := passiveSpecialConditionImmunityTypedEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":      "IMMUNE_TO_SPECIAL_CONDITIONS",
				"target":      "ALL_FRIENDLY",
				"trigger":     "HAS_ENERGY_ATTACHED",
				"energy_type": matches[1],
			},
		}}
	}

	// --- HEAL (All friendly of a type) ---
	if matches := healAllFriendlyTypedRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":     "ONCE_PER_TURN",
					"target_all":  true,
					"target_type": matches[2],
				},
			}}
		}
	}

	// --- SCALING DAMAGE (Coin flips per typed energy) ---
	if matches := scalingDamagePerTypedEnergyCoinFlipRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":            "COIN_FLIP_HEADS",
					"num_flips_scales_by": "SELF_ATTACHED_ENERGY_TYPED",
					"energy_type":         matches[1],
				},
			}}
		}
	}

	// --- HEAL (Benched) ---
	if matches := healBenchedRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Target:      core.TargetBenchedFriendly,
				Amount:      amount,
				Description: text,
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Different Energy Types) ---
	if matches := conditionalDamageDifferentEnergyRegex.FindStringSubmatch(text); len(matches) > 2 {
		count, err1 := strconv.Atoi(matches[1])
		amount, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":        "DIFFERENT_ENERGY_TYPES_ATTACHED",
					"required_count": count,
				},
			}}
		}
	}

	// --- DAMAGE (Equals self damage) ---
	if damageEqualsSelfDamageRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDamage, // This is a direct damage effect, not a bonus
			Description: text,
			Conditions: map[string]interface{}{
				"amount_equals": "SELF_DAMAGE_COUNTERS",
			},
		}}
	}

	// --- SET HP (On Coin Flip) ---
	if matches := setHPOnCoinFlipRegex.FindStringSubmatch(text); len(matches) > 1 {
		hp, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectSetHP,
				Target:      core.TargetOpponentActive,
				Amount:      hp,
				Description: text,
				Conditions: map[string]interface{}{
					"on_coin_flip": "HEADS",
				},
			}}
		}
	}

	// --- SHUFFLE FROM HAND (Multi-coin flip) ---
	if matches := shuffleFromHandMultiCoinFlipRegex.FindStringSubmatch(text); len(matches) > 1 {
		flips, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectShuffleFromHand,
				Target:      core.TargetOpponentHand,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":  "COIN_FLIP_HEADS",
					"num_flips": flips,
					"random":    true,
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Reactive Damage) ---
	if matches := passiveReactiveDamageActiveRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":   "REACTIVE_DAMAGE",
					"amount":   amount,
					"location": "ACTIVE",
				},
			}}
		}
	}

	// --- SHUFFLE FROM HAND (Reveal and choose) ---
	if shuffleFromHandRevealRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectShuffleFromHand,
			Target:      core.TargetOpponentHand,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"reveal_hand":    true,
				"player_chooses": true,
			},
		}}
	}

	// --- APPLY RESTRICTION (Can't use same attack) ---
	if matches := restrictionCantUseAttackRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "CANT_USE_ATTACK",
				"attack_name": matches[1],
				"duration":    "next_turn",
			},
		}}
	}

	// --- SCALING DAMAGE (Flip until tails) ---
	if matches := scalingDamageUntilTailsRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "COIN_FLIP_HEADS_UNTIL_TAILS",
				},
			}}
		}
	}

	// --- FORCE SWITCH (Opponent must choose damaged) ---
	if forceSwitchDamagedRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectForceSwitch,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":          "ONCE_PER_TURN",
				"location":         "ACTIVE",
				"target_condition": "HAS_DAMAGE",
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Double Heads) ---
	if matches := conditionalDamageDoubleHeadsRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"on_coin_flip": "DOUBLE_HEADS",
				},
			}}
		}
	}

	// --- SCALING DAMAGE (Per Pokémon in Play) ---
	if matches := scalingDamagePerPokemonInPlayRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":            "COIN_FLIP_HEADS",
					"num_flips_scales_by": "ALL_POKEMON_IN_PLAY",
				},
			}}
		}
	}

	// --- MODIFY ENERGY (Next generated) ---
	if matches := modifyNextEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		typesStr := strings.ReplaceAll(matches[1], "{", "")
		types := strings.Split(typesStr, "}, ")
		types[len(types)-1] = strings.TrimSuffix(types[len(types)-1], "}")

		return []core.Effect{{
			Type:        core.EffectModifyEnergy,
			Target:      core.TargetOpponentActive, // Implied target is opponent
			Description: text,
			Conditions: map[string]interface{}{
				"target_energy":  "NEXT_GENERATED",
				"random_type":    true,
				"possible_types": types,
			},
		}}
	}

	// --- PASSIVE ABILITY (Opponent Damage Reduction) ---
	if matches := passiveOpponentDamageReductionRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":   "REDUCE_OPPONENT_DAMAGE_OUTPUT",
					"amount":   amount,
					"location": "ACTIVE",
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (If damaged last turn) ---
	if matches := conditionalDamageIfDamagedLastTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "DAMAGED_LAST_TURN",
				},
			}}
		}
	}

	// --- LOOK AT DECK ---
	if lookAtTopCardRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectLookAtDeck,
			Target:      core.TargetDeck,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger": "ONCE_PER_TURN",
			},
		}}
	}

	// --- PASSIVE ABILITY (Buff Status Damage) ---
	if matches := passiveBuffStatusDamageRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect": "BUFF_STATUS_DAMAGE",
					"amount": amount,
					"status": strings.ToUpper(matches[2]),
				},
			}}
		}
	}

	// --- ATTACH ENERGY (Once per turn) ---
	if matches := attachEnergyOncePerTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetSelf,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":    "ONCE_PER_TURN",
				"source":     "EnergyZone",
				"energyType": matches[1],
			},
		}}
	}

	// --- SEARCH DECK (Generic Pokémon) ---
	if searchDeckRandomPokemonRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectSearchDeck,
			Target:      core.TargetDeck,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"pokemonType": "ANY",
				"random":      true,
				"destination": "hand",
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Self has damage) ---
	if matches := conditionalDamageSelfHasDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "SELF_HAS_DAMAGE",
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Attack history) ---
	if matches := conditionalDamageOnAttackHistoryRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":     "ATTACK_USED_LAST_TURN",
					"attack_name": matches[1],
				},
			}}
		}
	}

	// --- DISCARD ENERGY (Until tails) ---
	if discardEnergyUntilTailsRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"scale_by": "COIN_FLIP_HEADS_UNTIL_TAILS",
				"random":   true,
			},
		}}
	}

	// --- DELAYED DAMAGE ---
	if matches := delayedDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDelayedDamage,
				Target:      core.TargetOpponentActive,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "END_OF_OPPONENT_NEXT_TURN",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Reactive Damage on KO) ---
	if matches := passiveReactiveDamageOnKORegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":   "REACTIVE_DAMAGE_ON_KO",
					"amount":   amount,
					"location": "ACTIVE",
				},
			}}
		}
	}

	// --- APPLY RESTRICTION (Opponent's hand) ---
	if matches := restrictOpponentHandNextTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetOpponentActive, // Effect is on the opponent player
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "CANT_PLAY_CARD_TYPE",
				"card_type":   matches[1],
				"duration":    "opponent_next_turn",
			},
		}}
	}

	// --- PASSIVE ABILITY (Zero Retreat for Active) ---
	if passiveZeroRetreatForActiveRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect": "ZERO_RETREAT_COST",
				"target": "ACTIVE",
			},
		}}
	}

	// --- SCALING DAMAGE (Attack history count) ---
	if matches := scalingDamageAttackHistoryRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":    "ATTACK_USAGE_COUNT",
					"attack_name": matches[2],
				},
			}}
		}
	}

	// --- DISCARD DECK (Both Players) ---
	if matches := discardDeckBothPlayersRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDiscardDeck,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"target": "BOTH_PLAYERS",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Damage Reduction on Coin Flip) ---
	if matches := passiveDamageReductionOnCoinFlipRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":       "REDUCE_INCOMING_DAMAGE",
					"amount":       amount,
					"on_coin_flip": "HEADS",
				},
			}}
		}
	}

	// --- MOVE ENERGY (Benched to Active) ---
	if matches := moveEnergyBenchedToActiveRegex.FindStringSubmatch(text); len(matches) > 2 {
		return []core.Effect{{
			Type:        core.EffectMoveEnergy,
			Target:      core.TargetBenchedFriendly, // Source
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":     "ONCE_PER_TURN",
				"amount":      "ALL",
				"energyType":  matches[1],
				"source_type": matches[2],
				"destination": "ACTIVE",
			},
		}}
	}

	// --- APPLY RESTRICTION (Opponent cannot play Item) ---
	if restrictOpponentHandItemRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "CANT_PLAY_CARD_TYPE",
				"card_type":   "Item",
				"duration":    "opponent_next_turn",
			},
		}}
	}

	// --- PASSIVE ABILITY (Damage Buff for other Pokémon) ---
	if matches := passiveDamageBuffEvolvesFromRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":              "BUFF_DAMAGE_OUTPUT",
					"amount":              amount,
					"location":            "BENCH",
					"target_evolves_from": matches[1],
				},
			}}
		}
	}

	// --- FORCE SWITCH (Player chooses Benched Basic) ---
	if forceSwitchBenchedBasicRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectForceSwitch,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":        "ONCE_PER_TURN",
				"location":       "ACTIVE",
				"player_chooses": true,
				"target_pool":    "BENCHED",
				"target_stage":   "Basic",
			},
		}}
	}

	// --- APPLY STATUS (Both Active) ---
	if matches := applyStatusBothActiveRegex.FindStringSubmatch(text); len(matches) > 1 {
		status := core.StatusCondition(strings.ToUpper(matches[1]))
		effectSelf := core.Effect{
			Type:   core.EffectApplyStatus,
			Target: core.TargetSelf,
			Status: status,
		}
		effectOpponent := core.Effect{
			Type:   core.EffectApplyStatus,
			Target: core.TargetOpponentActive,
			Status: status,
		}
		return []core.Effect{effectSelf, effectOpponent}
	}

	// --- APPLY RESTRICTION (Can't Attack) ---
	if applyRestrictionCantAttackRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "CANT_ATTACK",
				"duration":    "opponent_next_turn",
			},
		}}
	}

	// --- PASSIVE ABILITY (Damage Reduction via other Pokémon) ---
	if matches := passiveDamageReductionInPlayRegex.FindStringSubmatch(text); len(matches) > 2 {
		names := strings.Split(matches[1], " or ")
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":           "REDUCE_INCOMING_DAMAGE",
					"amount":           amount,
					"requires_in_play": names,
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Global Damage Buff by Type) ---
	if matches := passiveDamageBuffTypedPokemonRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":      "BUFF_DAMAGE_OUTPUT",
					"amount":      amount,
					"target_type": matches[1],
				},
			}}
		}
	}

	// --- DRAW (With Discard Cost) ---
	if drawWithDiscardCostRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDraw,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":           "ONCE_PER_TURN",
				"cost_discard_hand": 1,
			},
		}}
	}

	// --- KNOCKOUT (On Coin Flip) ---
	if knockoutOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectKnockout,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "DOUBLE_HEADS",
			},
		}}
	}

	// --- PASSIVE ABILITY (Restrict Evolution) ---
	if passiveRestrictionEvolveRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect": "RESTRICT_OPPONENT_EVOLVE",
				"target": "ACTIVE",
			},
		}}
	}

	// --- SCALING DAMAGE (Benched count, base damage) ---
	if matches := scalingDamageBenchedBaseRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"is_base_damage": true, // Differentiates from bonus damage
					"scale_by":       "BENCHED_POKEMON_COUNT",
				},
			}}
		}
	}

	// --- SCALING DAMAGE (All Opponent's Energy) ---
	if matches := scalingDamageAllOpponentEnergyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"is_base_damage": true,
					"scale_by":       "ALL_OPPONENT_POKEMON_ENERGY",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Move Energy on KO) ---
	if matches := moveEnergyOnKORegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":      "MOVE_ENERGY_ON_KO",
				"energy_type": matches[1],
				"source":      "SELF",
				"destination": "BENCHED",
			},
		}}
	}

	// --- RECOIL DAMAGE (On KO) ---
	if matches := recoilDamageOnKORegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectRecoilDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_KO",
				},
			}}
		}
	}

	// --- REVEAL HAND ---
	if revealHandRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectRevealHand,
			Target:      core.TargetOpponentHand,
			Description: text,
		}}
	}

	// --- MOVE DAMAGE ---
	if moveDamageRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectMoveDamage,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":     "AS_OFTEN_AS_YOU_LIKE",
				"amount":      "ALL",
				"source":      "ANY_FRIENDLY_DAMAGED",
				"destination": "SELF",
			},
		}}
	}

	// --- DISCARD TOOL ---
	if discardAllToolsRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardTool,
			Target:      core.TargetOpponentActive,
			Amount:      99, // Represents "all"
			Description: text,
		}}
	}

	// --- PASSIVE ABILITY (Prevention from Pokémon ex) ---
	if passivePreventionFromEXRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":            "PREVENT_INCOMING_DAMAGE",
				"from_pokemon_type": "EX",
			},
		}}
	}

	// --- DISCARD ENERGY (Random Self, Multiple) ---
	if matches := discardRandomEnergySelfMultipleRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDiscardEnergy,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"random": true,
				},
			}}
		}
	}

	// --- SPLASH DAMAGE (Any Friendly) ---
	if matches := splashDamageAnyFriendlyRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDamageBenchedFriendly, // Re-using this type, but target pool is wider
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"target_pool": "ANY_FRIENDLY",
				},
			}}
		}
	}

	// --- APPLY RESTRICTION (Conditional on Stage) ---
	if matches := conditionalRestrictionOnStageRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction":     "CANT_ATTACK",
				"duration":        "opponent_next_turn",
				"target_if_stage": strings.ToUpper(matches[1]),
			},
		}}
	}

	// --- ATTACH ENERGY (Scaled by coin flips) ---
	if matches := attachEnergyScaledByCoinFlipsRegex.FindStringSubmatch(text); len(matches) > 3 {
		flips, err1 := strconv.Atoi(matches[1])
		if err1 == nil {
			return []core.Effect{{
				Type:        core.EffectAttachEnergy,
				Target:      core.TargetBenchedFriendly,
				Description: text,
				Conditions: map[string]interface{}{
					"source":            "EnergyZone",
					"energyType":        matches[2],
					"target_type":       matches[3],
					"scale_by":          "COIN_FLIP_HEADS",
					"num_flips":         flips,
					"distribute_freely": true,
				},
			}}
		}
	}

	// --- ATTACH ENERGY (Multiple Specific Types) ---
	if matches := attachMultipleSpecificEnergyRegex.FindStringSubmatch(text); len(matches) > 3 {
		types := []string{matches[1], matches[2], matches[3]}
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetBenchedFriendly,
			Description: text,
			Conditions: map[string]interface{}{
				"source":            "EnergyZone",
				"energyTypes":       types,
				"target_stage":      "Basic",
				"distribute_freely": true,
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Opponent has Special Condition) ---
	if matches := conditionalDamageOpponentHasStatusRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_HAS_SPECIAL_CONDITION",
				},
			}}
		}
	}

	// --- APPLY RESTRICTION (Attack may fail) ---
	if applyAttackFailureChanceRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "ATTACK_MAY_FAIL",
				"chance":      0.5,
				"on":          "TAILS",
				"duration":    "opponent_next_turn",
			},
		}}
	}

	// --- DISCARD FROM HAND (On Coin Flip) ---
	if discardFromHandOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardFromHand,
			Target:      core.TargetOpponentHand,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
				"random":       true,
			},
		}}
	}

	// --- MOVE ENERGY (As often as you like) ---
	if matches := moveEnergyAsOftenAsYouLikeRegex.FindStringSubmatch(text); len(matches) > 3 {
		return []core.Effect{{
			Type:        core.EffectMoveEnergy,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":          "AS_OFTEN_AS_YOU_LIKE",
				"amount":           1,
				"energyType":       matches[1],
				"source_type":      matches[2],
				"destination_type": matches[3],
				"source":           "BENCHED",
				"destination":      "ACTIVE",
			},
		}}
	}

	// --- APPLY RESTRICTION (Can't Attack, on coin flip) ---
	if conditionalRestrictionOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction":  "CANT_ATTACK",
				"duration":     "opponent_next_turn",
				"on_coin_flip": "HEADS",
			},
		}}
	}

	// --- SWITCH SELF (Subtype) ---
	if matches := switchSelfSubtypeRegex.FindStringSubmatch(text); len(matches) > 2 {
		return []core.Effect{{
			Type:        core.EffectSwitchSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":        "ONCE_PER_TURN",
				"source_subtype": matches[1],
				"target_subtype": matches[2],
			},
		}}
	}

	// --- ATTACH ENERGY (From Discard with Recoil) ---
	if matches := attachEnergyFromDiscardWithRecoilRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			attachEffect := core.Effect{
				Type:        core.EffectAttachEnergy,
				Description: "Attach energy from discard.",
				Conditions: map[string]interface{}{
					"trigger":    "ONCE_PER_TURN",
					"source":     "DISCARD_PILE",
					"energyType": matches[1],
				},
			}
			recoilEffect := core.Effect{
				Type:        core.EffectRecoilDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: "Take recoil damage.",
			}
			return []core.Effect{attachEffect, recoilEffect}
		}
	}

	// --- PASSIVE ABILITY (Evolve to any) ---
	if passiveEvolveToAnyRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect": "CAN_EVOLVE_INTO_ANY",
			},
		}}
	}

	// --- PASSIVE ABILITY (Cost Reduction with Tool) ---
	if matches := passiveCostReductionWithToolRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":     "REDUCE_ATTACK_COST",
					"trigger":    "SELF_HAS_TOOL",
					"amount":     amount,
					"energyType": matches[2],
				},
			}}
		}
	}

	// --- COPY ATTACK (With energy check) ---
	if copyAttackWithEnergyCheckRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectCopyAttack,
			Target:      core.TargetOpponentActive, // Text implies any opponent's Pokémon
			Description: text,
			Conditions: map[string]interface{}{
				"target_pool":     "ANY_OPPONENT",
				"requires_energy": true,
			},
		}}
	}

	// --- PASSIVE ABILITY (Increase Opponent Cost) ---
	if matches := passiveIncreaseOpponentCostRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":     "INCREASE_OPPONENT_ATTACK_COST",
					"location":   "ACTIVE",
					"amount":     amount,
					"energyType": matches[2],
				},
			}}
		}
	}

	// --- ATTACH ENERGY (To Typed Benched) ---
	if matches := attachEnergyToTypedBenchedRegex.FindStringSubmatch(text); len(matches) > 2 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetBenchedFriendly,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"source":      "EnergyZone",
				"energyType":  matches[1],
				"target_type": matches[2],
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Opponent is Stage) ---
	if matches := conditionalDamageOpponentStageRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":        "OPPONENT_IS_STAGE",
					"opponent_stage": strings.ToUpper(matches[1]),
				},
			}}
		}
	}

	// --- LOOK AT DECK (Either player) ---
	if lookAtEitherPlayerDeckRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectLookAtDeck,
			Target:      core.TargetDeck,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":       "ONCE_PER_TURN",
				"target_player": "EITHER",
			},
		}}
	}

	// --- APPLY RESTRICTION (Persistent attack failure) ---
	if persistentAttackFailureRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "ATTACK_MAY_FAIL",
				"chance":      0.5,
				"on":          "TAILS",
				"duration":    "PERSISTENT",
			},
		}}
	}

	// --- HEAL (On Coin Flip) ---
	if matches := healOnCoinFlipRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"on_coin_flip": "HEADS",
				},
			}}
		}
	}

	// --- DISCARD DECK (Self) ---
	if matches := discardOwnDeckAmountRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDiscardDeck,
				Target:      core.TargetDeck,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"target_player": "SELF",
				},
			}}
		}
	}

	// --- DISCARD DECK & CONDITIONAL DAMAGE ---
	if matches := discardDeckWithConditionalDamageRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			discardEffect := core.Effect{
				Type:        core.EffectDiscardDeck,
				Target:      core.TargetDeck,
				Amount:      1,
				Description: "Discard top card of your deck.",
				Conditions: map[string]interface{}{
					"target_player": "SELF",
				},
			}
			damageEffect := core.Effect{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: "Do more damage if discarded card is a Pokémon of a certain type.",
				Conditions: map[string]interface{}{
					"trigger":        "DISCARDED_CARD_IS_TYPE",
					"discarded_type": matches[1],
				},
			}
			return []core.Effect{discardEffect, damageEffect}
		}
	}

	// --- CONDITIONAL DAMAGE (Pokémon on Bench by name) ---
	if matches := conditionalDamageOnBenchedNameRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":      "POKEMON_ON_BENCH",
					"pokemon_name": matches[1],
				},
			}}
		}
	}

	// --- DAMAGE (Halve HP) ---
	if damageHalveHPRoundedDownRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDamageHalveHP,
			Target:      core.TargetOpponentActive,
			Description: text,
		}}
	}

	// --- PASSIVE ABILITY (Global Damage Reduction - Unown) ---
	if matches := passiveGlobalDamageReductionUnownRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":           "REDUCE_INCOMING_DAMAGE",
					"target":           "ALL_FRIENDLY",
					"amount":           amount,
					"requires_in_play": []string{"Unown"}, // Special condition for this ability
				},
			}}
		}
	}

	// --- APPLY STATUS (Random) ---
	if matches := applyRandomStatusRegex.FindStringSubmatch(text); len(matches) > 1 {
		// "Asleep, Burned, Confused, Paralyzed, and Poisoned"
		statusListStr := strings.ReplaceAll(matches[1], ", and ", ", ")
		possibleStatuses := strings.Split(statusListStr, ", ")
		return []core.Effect{{
			Type:        core.EffectApplyStatus,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"random":            true,
				"possible_statuses": possibleStatuses,
			},
		}}
	}

	// --- SCALING DAMAGE (By Discarding Tools from Hand) ---
	if matches := scalingDamageDiscardToolRegex.FindStringSubmatch(text); len(matches) > 2 {
		maxDiscard, err1 := strconv.Atoi(matches[1])
		damagePer, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      damagePer,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":    "DISCARD_TOOL_FROM_HAND",
					"max_discard": maxDiscard,
				},
			}}
		}
	}

	// --- SEARCH DECK (Tool) ---
	if searchDeckToolRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectSearchDeck,
			Target:      core.TargetDeck,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":     "ONCE_PER_TURN",
				"card_type":   "Pokémon Tool",
				"random":      true,
				"destination": "hand",
			},
		}}
	}

	// --- ATTACH ENERGY (End of first turn) ---
	if matches := attachEnergyAtEndOfFirstTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Target:      core.TargetSelf,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":    "END_OF_FIRST_TURN",
				"source":     "EnergyZone",
				"energyType": matches[1],
			},
		}}
	}

	// --- EVOLVE (On Energy Attach) ---
	if evolveOnEnergyAttachRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":  "EVOLVE",
				"trigger": "ATTACH_ENERGY_TO_SELF",
				"source":  "DECK",
				"random":  true,
			},
		}}
	}

	// --- HEAL (Active, once per turn) ---
	if matches := healActiveOncePerTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Target:      core.TargetSelf, // Target is Active, which is self in this context
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "ONCE_PER_TURN",
				},
			}}
		}
	}

	// --- ATTACH ENERGY (To any typed friendly) ---
	if matches := attachEnergyToAnyTypedFriendlyRegex.FindStringSubmatch(text); len(matches) > 2 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":     "ONCE_PER_TURN",
				"location":    "ACTIVE",
				"source":      "EnergyZone",
				"energyType":  matches[1],
				"target_type": matches[2],
				"target_pool": "ANY_FRIENDLY",
			},
		}}
	}

	// --- PASSIVE ABILITY (Zero Retreat via other Pokémon) ---
	if matches := passiveZeroRetreatInPlayRegex.FindStringSubmatch(text); len(matches) > 1 {
		names := strings.Split(matches[1], " or ")
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":           "ZERO_RETREAT_COST",
				"requires_in_play": names,
			},
		}}
	}

	// --- SCALING DAMAGE (Self Attached Energy - All Types) ---
	if matches := scalingDamageSelfEnergyAllTypesRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "SELF_ATTACHED_ENERGY",
				},
			}}
		}
	}

	// --- RETURN TO HAND (On Coin Flip) ---
	if returnToHandOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectReturnToHand,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
			},
		}}
	}

	// --- DISCARD BENCHED for SCALING DAMAGE ---
	if matches := discardBenchedForScalingDamageRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDiscardBenched,
				Target:      core.TargetBenchedFriendly, // Player chooses which/how many
				Description: "You may discard any number of your Benched {W} Pokémon.",
				Conditions: map[string]interface{}{
					"target_type": matches[1],
				},
			}, {
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: "This attack does 40 more damage for each Benched Pokémon you discarded in this way.",
				Conditions: map[string]interface{}{
					"scale_by": "DISCARDED_BENCHED_COUNT",
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Global Heal Block) ---
	if passiveGlobalHealBlockRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect": "PREVENT_HEALING",
				"target": "GLOBAL",
			},
		}}
	}

	// --- DEVOLVE ---
	if devolveOnConditionRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDevolve,
			Target:      core.TargetOpponentActive,
			Amount:      1, // Highest stage = 1 level
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":     "OPPONENT_IS_EVOLVED",
				"destination": "HAND",
			},
		}}
	}

	// --- SHUFFLE HAND AND DRAW SCALED ---
	if shuffleHandAndDrawScaledRegex.MatchString(text) {
		shuffleEffect := core.Effect{
			Type:        core.EffectShuffleFromHand,
			Target:      core.TargetSelf, // Shuffle own hand
			Description: "Shuffle your hand into your deck.",
			Conditions: map[string]interface{}{
				"destination": "DECK",
			},
		}
		drawEffect := core.Effect{
			Type:        core.EffectDraw,
			Description: "Draw a card for each card in your opponent's hand.",
			Conditions: map[string]interface{}{
				"scale_by": "OPPONENT_HAND_SIZE",
			},
		}
		return []core.Effect{shuffleEffect, drawEffect}
	}

	// --- DISCARD ENERGY (Both Active) ---
	if discardEnergyBothActiveRegex.MatchString(text) {
		discardSelf := core.Effect{
			Type:       core.EffectDiscardEnergy,
			Target:     core.TargetSelf,
			Amount:     1,
			Conditions: map[string]interface{}{"random": true},
		}
		discardOpponent := core.Effect{
			Type:       core.EffectDiscardEnergy,
			Target:     core.TargetOpponentActive,
			Amount:     1,
			Conditions: map[string]interface{}{"random": true},
		}
		return []core.Effect{discardSelf, discardOpponent}
	}

	// --- CONDITIONAL DAMAGE (Benched are damaged) ---
	if matches := conditionalDamageBenchedDamagedRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "ANY_BENCHED_FRIENDLY_HAS_DAMAGE",
				},
			}}
		}
	}

	// --- SHUFFLE FROM HAND (On Coin Flip) ---
	if shuffleFromHandOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectShuffleFromHand,
			Target:      core.TargetOpponentHand,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
				"reveal":       true,
				"random":       true,
				"destination":  "DECK",
			},
		}}
	}

	// --- DISCARD ENERGY (Opponent, on Coin Flip) ---
	if discardEnergyOpponentOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetOpponentActive,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
				"random":       true,
			},
		}}
	}

	// --- DRAW (Until hand size matches opponent) ---
	if drawUntilMatchHandSizeRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDraw,
			Description: text,
			Conditions: map[string]interface{}{
				"draw_until": "MATCH_OPPONENT_HAND_SIZE",
			},
		}}
	}

	// --- SCALING DAMAGE (Opponent Benched Count) ---
	if matches := scalingDamageOpponentBenchedCountRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by": "OPPONENT_BENCHED_POKEMON_COUNT",
				},
			}}
		}
	}

	// --- COPY ATTACK (On Coin Flip) ---
	if copyAttackOnCoinFlipRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectCopyAttack,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
			},
		}}
	}

	// --- SNIPE DAMAGE (Random Opponent Pokémon) ---
	if matches := snipeRandomOpponentRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectSnipeDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"target_pool": "ANY_OPPONENT",
					"random":      true,
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Global Damage Buff - Unown) ---
	if matches := passiveGlobalDamageBuffUnownRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect":           "BUFF_DAMAGE_OUTPUT",
					"target":           "ALL_FRIENDLY",
					"amount":           amount,
					"requires_in_play": []string{"Unown"}, // Special condition
				},
			}}
		}
	}

	// --- APPLY PREVENTION (On KO) ---
	if preventionOnKORegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility, // This is a passive trigger
			Description: text,
			Conditions: map[string]interface{}{
				"effect":   "APPLY_PREVENTION_ON_KO",
				"prevent":  "ALL_DAMAGE_AND_EFFECTS",
				"duration": "opponent_next_turn",
			},
		}}
	}

	// --- MOVE ENERGY (All to Benched) ---
	if moveAllEnergyToBenchedRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectMoveEnergy,
			Target:      core.TargetBenchedFriendly,
			Description: text,
			Conditions: map[string]interface{}{
				"amount":      "ALL",
				"source":      "SELF",
				"destination": "BENCHED",
			},
		}}
	}

	// --- APPLY RESTRICTION (Energy Attachment) ---
	if restrictEnergyAttachmentRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction": "CANT_ATTACH_ENERGY",
				"target":      "ACTIVE",
				"duration":    "opponent_next_turn",
			},
		}}
	}

	// --- SHUFFLE FROM HAND (Simple) ---
	if shuffleFromHandSimpleRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectShuffleFromHand,
			Target:      core.TargetOpponentHand,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"reveal":      true,
				"random":      true,
				"destination": "DECK",
			},
		}}
	}

	// --- SWITCH SELF (Voluntary) ---
	if voluntarySwitchRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectSwitchSelf,
			Target:      core.TargetBenchedFriendly,
			Description: text,
			Conditions: map[string]interface{}{
				"voluntary": true,
			},
		}}
	}

	// --- REVEAL HAND (On Bench Play) ---
	if revealHandOnBenchPlayRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectRevealHand,
			Target:      core.TargetOpponentHand,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger": "ON_PLAY_TO_BENCH",
			},
		}}
	}

	// --- DEBUFF INCOMING DAMAGE ---
	if matches := debuffIncomingDamageRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDebuffIncomingDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"duration": "opponent_next_turn",
				},
			}}
		}
	}

	// --- APPLY RESTRICTION (Can't attack on TAILS) ---
	if restrictionOnTailsRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyRestriction,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"restriction":  "CANT_ATTACK",
				"on_coin_flip": "TAILS",
				"duration":     "next_turn",
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Opponent Name) ---
	if matches := conditionalDamageOpponentNameRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":       "OPPONENT_IS_NAME",
					"opponent_name": matches[1],
				},
			}}
		}
	}

	// --- ATTACH ENERGY (To other specific Pokémon) ---
	if matches := attachEnergyToSpecificPokemonRegex.FindStringSubmatch(text); len(matches) > 3 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"source":       "EnergyZone",
				"energyType":   matches[1],
				"target_names": []string{matches[2], matches[3]},
			},
		}}
	}

	// --- PASSIVE ABILITY (Retreat Cost on First Turn) ---
	if passiveRetreatCostFirstTurnRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":   "ZERO_RETREAT_COST",
				"duration": "FIRST_TURN",
			},
		}}
	}

	// --- PASSIVE ABILITY (Effect Prevention) ---
	if passiveEffectPreventionRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect": "PREVENT_INCOMING_EFFECTS",
			},
		}}
	}

	// --- APPLY STATUS (Once per turn ability) ---
	if matches := applyStatusOncePerTurnAbilityRegex.FindStringSubmatch(text); len(matches) > 1 {
		return []core.Effect{{
			Type:        core.EffectApplyStatus,
			Target:      core.TargetOpponentActive,
			Status:      core.StatusCondition(strings.ToUpper(matches[1])),
			Description: text,
			Conditions: map[string]interface{}{
				"trigger": "ONCE_PER_TURN",
			},
		}}
	}

	// --- BUFF (Stacking) ---
	if matches := buffStackingRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectBuffNextTurn,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"attack_name": matches[1],
					"stacking":    true,
					"duration":    "PERSISTENT_ACTIVE",
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (If any typed energy attached) ---
	if matches := conditionalDamageIfEnergyAttachedRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":     "SELF_HAS_TYPED_ENERGY",
					"energy_type": matches[1],
				},
			}}
		}
	}

	// --- HEAL (On Energy Attach) ---
	if matches := healOnEnergyAttachRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectHeal,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":     "ATTACH_ENERGY_TO_SELF",
					"energy_type": matches[1],
				},
			}}
		}
	}

	// --- KNOCKOUT (Attacker, on coin flip when KO'd) ---
	if knockoutAttackerOnKORegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectPassiveAbility,
			Description: text,
			Conditions: map[string]interface{}{
				"effect":       "KO_ATTACKER_ON_KO",
				"on_coin_flip": "HEADS",
			},
		}}
	}

	// --- DAMAGE (Ability, requires other Pokémon in play) ---
	if matches := damageAbilityInPlayRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[2])
		names := strings.Split(matches[1], " or ")
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDamage,
				Target:      core.TargetOpponentActive,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger":          "ONCE_PER_TURN",
					"requires_in_play": names,
				},
			}}
		}
	}

	// --- CONDITIONAL DAMAGE (Multi-coin flip) ---
	if matches := finalConditionalDamageMultiCoinFlipRegex.FindStringSubmatch(text); len(matches) > 2 {
		flips, err1 := strconv.Atoi(matches[1])
		amount, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"scale_by":  "COIN_FLIP_HEADS",
					"num_flips": flips,
				},
			}}
		}
	}

	// --- PASSIVE ABILITY (Simple Damage Reduction) ---
	if matches := finalPassiveDamageReductionRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectPassiveAbility,
				Description: text,
				Conditions: map[string]interface{}{
					"effect": "REDUCE_INCOMING_DAMAGE",
					"amount": amount,
				},
			}}
		}
	}

	// --- DISCARD ENERGY (Opponent, simple) ---
	if finalDiscardEnergyOpponentSimpleRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectDiscardEnergy,
			Target:      core.TargetOpponentActive,
			Amount:      1,
			Description: text,
			Conditions:  map[string]interface{}{"random": true},
		}}
	}

	// --- CONDITIONAL DAMAGE (Flip until tails) ---
	if matches := finalConditionalDamageUntilTailsRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions:  map[string]interface{}{"scale_by": "COIN_FLIP_HEADS_UNTIL_TAILS"},
			}}
		}
	}

	// --- SCALING DAMAGE (Opponent Energy, Base Damage) ---
	if matches := finalScalingDamageOpponentEnergyBaseRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectScalingDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"is_base_damage": true,
					"scale_by":       "OPPONENT_ATTACHED_ENERGY",
				},
			}}
		}
	}

	// --- ATTACH ENERGY (To Active Typed Pokémon) ---
	if matches := finalAttachEnergyToActiveTypedRegex.FindStringSubmatch(text); len(matches) > 2 {
		return []core.Effect{{
			Type:        core.EffectAttachEnergy,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":         "ONCE_PER_TURN",
				"source":          "EnergyZone",
				"energyType":      matches[1],
				"target_location": "ACTIVE",
				"target_type":     matches[2],
			},
		}}
	}

	// --- CONDITIONAL DAMAGE (Opponent is EX) ---
	if matches := finalConditionalDamageOpponentIsEXRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectConditionalDamage,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"trigger": "OPPONENT_IS_EX",
				},
			}}
		}
	}

	// --- DISCARD ENERGY (Self, on tails) ---
	if matches := finalDiscardSelfEnergyOnTailsRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectDiscardEnergy,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions: map[string]interface{}{
					"on_coin_flip": "TAILS",
					"random":       true,
				},
			}}
		}
	}

	// --- REDUCE INCOMING DAMAGE (Next Turn) ---
	if matches := finalReduceDamageNextTurnRegex.FindStringSubmatch(text); len(matches) > 1 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectReduceIncomingDamage,
				Target:      core.TargetSelf,
				Amount:      amount,
				Description: text,
				Conditions:  map[string]interface{}{"duration": "opponent_next_turn"},
			}}
		}
	}

	// --- FORCE SWITCH (On Heads) ---
	if finalForceSwitchOnHeadsRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectForceSwitch,
			Target:      core.TargetOpponentActive,
			Description: text,
			Conditions:  map[string]interface{}{"on_coin_flip": "HEADS"},
		}}
	}

	// --- INCREASE OPPONENT COST (Next Turn) ---
	if matches := finalIncreaseOpponentCostNextTurnRegex.FindStringSubmatch(text); len(matches) > 2 {
		amount, err := strconv.Atoi(matches[1])
		if err == nil {
			return []core.Effect{{
				Type:        core.EffectApplyRestriction,
				Target:      core.TargetOpponentActive,
				Description: text,
				Conditions: map[string]interface{}{
					"restriction": "INCREASE_ATTACK_COST",
					"amount":      amount,
					"energyType":  matches[2],
					"duration":    "opponent_next_turn",
				},
			}}
		}
	}

	// --- PREVENTION (On Heads) ---
	if finalPreventionOnHeadsRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectApplyPrevention,
			Target:      core.TargetSelf,
			Description: text,
			Conditions: map[string]interface{}{
				"on_coin_flip": "HEADS",
				"duration":     "opponent_next_turn",
				"prevent":      "ALL_DAMAGE_AND_EFFECTS",
			},
		}}
	}

	// --- SEARCH DECK (Generic Pokémon) ---
	if finalSearchDeckGenericPokemonRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectSearchDeck,
			Target:      core.TargetDeck,
			Amount:      1,
			Description: text,
			Conditions: map[string]interface{}{
				"trigger":     "ONCE_PER_TURN",
				"pokemonType": "ANY",
				"random":      true,
				"destination": "hand",
			},
		}}
	}

	// --- FORCE SWITCH (Simple) ---
	if finalForceSwitchSimpleRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectForceSwitch,
			Target:      core.TargetOpponentActive,
			Description: text,
		}}
	}

	// --- FINAL HARIYAMA FIX ---
	if hariyamaPushOutRegex.MatchString(text) {
		return []core.Effect{{
			Type:        core.EffectForceSwitch,
			Target:      core.TargetOpponentActive,
			Description: text,
		}}
	}

	// --- Fallback for unknown effects ---
	// If no other rule matches, we return an UNKNOWN effect type. This allows us
	// to see which effects we still need to implement parsing for.
	return []core.Effect{{
		Type:        core.EffectUnknown,
		Description: text,
	}}
}
