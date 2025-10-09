package effects

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cpritch/genomon/internal/core"
)

// REGEX DEFINITIONS
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

// effectParser pairs a regular expression with a function that can parse its matches.
type effectParser struct {
	regex   *regexp.Regexp
	handler func(matches []string, text string) []core.Effect
}

// effectParsers holds our list of all known effect parsers.
var effectParsers = []effectParser{
	{regex: healAllFriendlyRegex, handler: parseHealAllFriendly},
	{regex: healRegex, handler: parseHeal},
	{regex: cantAttackNextTurnRegex, handler: parseCantAttackNextTurn},
	{regex: forceSwitchRegex, handler: parseForceSwitch},
	{regex: recoilDamageRegex, handler: parseRecoilDamage},
	{regex: searchDeckRegex, handler: parseSearchDeck},
	{regex: conditionalDamageRegex, handler: parseConditionalDamage},
	{regex: applyStatusOpponentRegex, handler: parseApplyStatusOpponent},
	{regex: applyStatusSelfRegex, handler: parseApplyStatusSelf},
	{regex: attachEnergyRegex, handler: parseAttachEnergy},
	{regex: searchEvolutionRegex, handler: parseSearchEvolution},
	{regex: triggeredSleepRegex, handler: parseTriggeredSleep},
	{regex: conditionalDamageHPRatioRegex, handler: parseConditionalDamageHPRatio},
	{regex: conditionalDamageEvolvedTurnRegex, handler: parseConditionalDamageEvolvedTurn},
	{regex: passiveDamageCheckupRegex, handler: parsePassiveDamageCheckup},
	{regex: scalingDamageRetreatCostRegex, handler: parseScalingDamageRetreatCost},
	{regex: coinFlipTailsFailsRegex, handler: parseCoinFlipTailsFails},
	{regex: passiveRetreatCostLatiasRegex, handler: parsePassiveRetreatCostLatias},
	{regex: discardAllEnergyRegex, handler: parseDiscardAllEnergy},
	{regex: discardTypedEnergyRegex, handler: parseDiscardTypedEnergy},
	{regex: moveAllTypedEnergyRegex, handler: parseMoveAllTypedEnergy},
	{regex: reduceIncomingDamageRegex, handler: parseReduceIncomingDamage},
	{regex: discardFromHandRegex, handler: parseDiscardFromHand},
	{regex: cantRetreatRegex, handler: parseCantRetreat},
	{regex: multiHitRandomRegex, handler: parseMultiHitRandom},
	{regex: statusOnCoinFlipRegex, handler: parseStatusOnCoinFlip},
	{regex: scalingDamageSelfDamageRegex, handler: parseScalingDamageSelfDamage},
	{regex: splashDamageBenchedFriendlyRegex, handler: parseSplashDamageBenchedFriendly},
	{regex: conditionalDamageExistingDamageRegex, handler: parseConditionalDamageExistingDamage},
	{regex: conditionalDamageToolAttachedRegex, handler: parseConditionalDamageToolAttached},
	{regex: gutsAbilityRegex, handler: parseGutsAbility},
	{regex: searchDeckByNameRegex, handler: parseSearchDeckByName},
	{regex: snipeAndDiscardTypedRegex, handler: parseSnipeAndDiscardTyped},
	{regex: discardRandomEnergySelfRegex, handler: parseDiscardRandomEnergySelf},
	{regex: healOnEvolveRegex, handler: parseHealOnEvolve},
	{regex: forceSwitchOncePerTurnRegex, handler: parseForceSwitchOncePerTurn},
	{regex: reduceDamageNextTurnSelfRegex, handler: parseReduceDamageNextTurnSelf},
	{regex: copyAttackRegex, handler: parseCopyAttack},
	{regex: passiveDamageReductionTypedRegex, handler: parsePassiveDamageReductionTyped},
	{regex: drawAtEndOfTurnRegex, handler: parseDrawAtEndOfTurn},
	{regex: applyStatusOncePerTurnRegex, handler: parseApplyStatusOncePerTurn},
	{regex: discardRandomEnergyGlobalRegex, handler: parseDiscardRandomEnergyGlobal},
	{regex: switchSelfRegex, handler: parseSwitchSelf},
	{regex: snipeDamageRegex, handler: parseSnipeDamage},
	{regex: healOncePerTurnRegex, handler: parseHealOncePerTurn},
	{regex: attachEnergyEndsTurnRegex, handler: parseAttachEnergyEndsTurn},
	{regex: discardEnergyOnEvolveRegex, handler: parseDiscardEnergyOnEvolve},
	{regex: scalingDamageSelfEnergyRegex, handler: parseScalingDamageSelfEnergy},
	{regex: drawCardRegex, handler: parseDrawCard},
	{regex: attachEnergyMultiBenchedRegex, handler: parseAttachEnergyMultiBenched},
	{regex: conditionalDamageSupporterRegex, handler: parseConditionalDamageSupporter},
	{regex: drawOnEvolveRegex, handler: parseDrawOnEvolve},
	{regex: discardSingleTypedEnergyRegex, handler: parseDiscardSingleTypedEnergy},
	{regex: shuffleIntoDeckOnCoinFlipRegex, handler: parseShuffleIntoDeckOnCoinFlip},
	{regex: scalingDamageOpponentEnergyRegex, handler: parseScalingDamageOpponentEnergy},
	{regex: preventAllDamageOnCoinFlipRegex, handler: parsePreventAllDamageOnCoinFlip},
	{regex: increaseOpponentCostsRegex, handler: parseIncreaseOpponentCosts},
	{regex: snipeDamagedPokemonRegex, handler: parseSnipeDamagedPokemon},
	{regex: switchSelfFromBenchRegex, handler: parseSwitchSelfFromBench},
	{regex: conditionalDamageSelfToolRegex, handler: parseConditionalDamageSelfTool},
	{regex: attachMultipleEnergyRegex, handler: parseAttachMultipleEnergy},
	{regex: passiveCostReductionInPlayRegex, handler: parsePassiveCostReductionInPlay},
	{regex: conditionalDamageOpponentAbilityRegex, handler: parseConditionalDamageOpponentAbility},
	{regex: scalingSnipeOpponentEnergyRegex, handler: parseScalingSnipeOpponentEnergy},
	{regex: snipeOncePerTurnRegex, handler: parseSnipeOncePerTurn},
	{regex: scalingDamageBenchedNameRegex, handler: parseScalingDamageBenchedName},
	{regex: damageBenchedOpponentAllRegex, handler: parseDamageBenchedOpponentAll},
	{regex: lifestealRegex, handler: parseLifesteal},
	{regex: passiveEnergyValueRegex, handler: parsePassiveEnergyValue},
	{regex: scalingDamageMultiCoinFlipRegex, handler: parseScalingDamageMultiCoinFlip},
	{regex: conditionalDamageSwitchInRegex, handler: parseConditionalDamageSwitchIn},
	{regex: conditionalDamageCoinFlipScalingRegex, handler: parseConditionalDamageCoinFlipScaling},
	{regex: passiveImmunityRegex, handler: parsePassiveImmunity},
	{regex: scalingDamageBenchedCountRegex, handler: parseScalingDamageBenchedCount},
	{regex: reactiveDamageRegex, handler: parseReactiveDamage},
	{regex: attachEnergyToActiveTypedRegex, handler: parseAttachEnergyToActiveTyped},
	{regex: buffNextTurnRegex, handler: parseBuffNextTurn},
	{regex: attachEnergyToBenchedRegex, handler: parseAttachEnergyToBenched},
	{regex: scalingDamageBenchedTypeRegex, handler: parseScalingDamageBenchedType},
	{regex: passiveRestrictionRegex, handler: parsePassiveRestriction},
	{regex: modifyEnergyRegex, handler: parseModifyEnergy},
	{regex: triggeredDamageOnEnergyAttachRegex, handler: parseTriggeredDamageOnEnergyAttach},
	{regex: switchSelfTypedRegex, handler: parseSwitchSelfTyped},
	{regex: conditionalDamageOpponentPropertyRegex, handler: parseConditionalDamageOpponentProperty},
	{regex: conditionalDamageOnKORegex, handler: parseConditionalDamageOnKO},
	{regex: splashDamageBenchedFriendlyAllRegex, handler: parseSplashDamageBenchedFriendlyAll},
	{regex: scalingDamageBenchedNamesRegex, handler: parseScalingDamageBenchedNames},
	{regex: discardMultipleTypedEnergyRegex, handler: parseDiscardMultipleTypedEnergy},
	{regex: attachEnergyToBenchedStageRegex, handler: parseAttachEnergyToBenchedStage},
	{regex: passiveDamageBuffInPlayRegex, handler: parsePassiveDamageBuffInPlay},
	{regex: passiveDamagePreventionCoinFlipRegex, handler: parsePassiveDamagePreventionCoinFlip},
	{regex: healAtEndOfTurnRegex, handler: parseHealAtEndOfTurn},
	{regex: passiveRetreatCostEnergyRegex, handler: parsePassiveRetreatCostEnergy},
	{regex: damageAllOpponentRegex, handler: parseDamageAllOpponent},
	{regex: passiveRetreatCostReductionBenchRegex, handler: parsePassiveRetreatCostReductionBench},
	{regex: damageBenchedConditionalEnergyRegex, handler: parseDamageBenchedConditionalEnergy},
	{regex: scalingDamagePerEnergyCoinFlipRegex, handler: parseScalingDamagePerEnergyCoinFlip},
	{regex: alternateAttackCostRegex, handler: parseAlternateAttackCost},
	{regex: discardAllTypedEnergyRegex, handler: parseDiscardAllTypedEnergy},
	{regex: conditionalDamageOpponentStatusRegex, handler: parseConditionalDamageOpponentStatus},
	{regex: preventAllDamageSimpleRegex, handler: parsePreventAllDamageSimple},
	{regex: scalingDamageBenchedTypeCountRegex, handler: parseScalingDamageBenchedTypeCount},
	{regex: passiveDamageReductionSimpleRegex, handler: parsePassiveDamageReductionSimple},
	{regex: conditionalDamageNoDamageRegex, handler: parseConditionalDamageNoDamage},
	{regex: passiveRetreatCostForOtherRegex, handler: parsePassiveRetreatCostForOther},
	{regex: recoilDamageOnCoinFlipRegex, handler: parseRecoilDamageOnCoinFlip},
	{regex: passiveImmunitySingleStatusRegex, handler: parsePassiveImmunitySingleStatus},
	{regex: discardDeckRegex, handler: parseDiscardDeck},
	{regex: splashDamageSingleBenchedOpponentRegex, handler: parseSplashDamageSingleBenchedOpponent},
	{regex: scalingDamageAllBenchedRegex, handler: parseScalingDamageAllBenched},
	{regex: multiHitRandomGlobalRegex, handler: parseMultiHitRandomGlobal},
	{regex: passiveSpecialConditionImmunityTypedEnergyRegex, handler: parsePassiveSpecialConditionImmunityTypedEnergy},
	{regex: healAllFriendlyTypedRegex, handler: parseHealAllFriendlyTyped},
	{regex: scalingDamagePerTypedEnergyCoinFlipRegex, handler: parseScalingDamagePerTypedEnergyCoinFlip},
	{regex: healBenchedRegex, handler: parseHealBenched},
	{regex: conditionalDamageDifferentEnergyRegex, handler: parseConditionalDamageDifferentEnergy},
	{regex: damageEqualsSelfDamageRegex, handler: parseDamageEqualsSelfDamage},
	{regex: setHPOnCoinFlipRegex, handler: parseSetHPOnCoinFlip},
	{regex: shuffleFromHandMultiCoinFlipRegex, handler: parseShuffleFromHandMultiCoinFlip},
	{regex: passiveReactiveDamageActiveRegex, handler: parsePassiveReactiveDamageActive},
	{regex: shuffleFromHandRevealRegex, handler: parseShuffleFromHandReveal},
	{regex: restrictionCantUseAttackRegex, handler: parseRestrictionCantUseAttack},
	{regex: scalingDamageUntilTailsRegex, handler: parseScalingDamageUntilTails},
	{regex: forceSwitchDamagedRegex, handler: parseForceSwitchDamaged},
	{regex: conditionalDamageDoubleHeadsRegex, handler: parseConditionalDamageDoubleHeads},
	{regex: scalingDamagePerPokemonInPlayRegex, handler: parseScalingDamagePerPokemonInPlay},
	{regex: modifyNextEnergyRegex, handler: parseModifyNextEnergy},
	{regex: passiveOpponentDamageReductionRegex, handler: parsePassiveOpponentDamageReduction},
	{regex: conditionalDamageIfDamagedLastTurnRegex, handler: parseConditionalDamageIfDamagedLastTurn},
	{regex: lookAtTopCardRegex, handler: parseLookAtTopCard},
	{regex: passiveBuffStatusDamageRegex, handler: parsePassiveBuffStatusDamage},
	{regex: attachEnergyOncePerTurnRegex, handler: parseAttachEnergyOncePerTurn},
	{regex: searchDeckRandomPokemonRegex, handler: parseSearchDeckRandomPokemon},
	{regex: conditionalDamageSelfHasDamageRegex, handler: parseConditionalDamageSelfHasDamage},
	{regex: conditionalDamageOnAttackHistoryRegex, handler: parseConditionalDamageOnAttackHistory},
	{regex: discardEnergyUntilTailsRegex, handler: parseDiscardEnergyUntilTails},
	{regex: delayedDamageRegex, handler: parseDelayedDamage},
	{regex: passiveReactiveDamageOnKORegex, handler: parsePassiveReactiveDamageOnKO},
	{regex: restrictOpponentHandNextTurnRegex, handler: parseRestrictOpponentHandNextTurn},
	{regex: passiveZeroRetreatForActiveRegex, handler: parsePassiveZeroRetreatForActive},
	{regex: scalingDamageAttackHistoryRegex, handler: parseScalingDamageAttackHistory},
	{regex: discardDeckBothPlayersRegex, handler: parseDiscardDeckBothPlayers},
	{regex: passiveDamageReductionOnCoinFlipRegex, handler: parsePassiveDamageReductionOnCoinFlip},
	{regex: moveEnergyBenchedToActiveRegex, handler: parseMoveEnergyBenchedToActive},
	{regex: restrictOpponentHandItemRegex, handler: parseRestrictOpponentHandItem},
	{regex: passiveDamageBuffEvolvesFromRegex, handler: parsePassiveDamageBuffEvolvesFrom},
	{regex: forceSwitchBenchedBasicRegex, handler: parseForceSwitchBenchedBasic},
	{regex: applyStatusBothActiveRegex, handler: parseApplyStatusBothActive},
	{regex: applyRestrictionCantAttackRegex, handler: parseApplyRestrictionCantAttack},
	{regex: passiveDamageReductionInPlayRegex, handler: parsePassiveDamageReductionInPlay},
	{regex: passiveDamageBuffTypedPokemonRegex, handler: parsePassiveDamageBuffTypedPokemon},
	{regex: drawWithDiscardCostRegex, handler: parseDrawWithDiscardCost},
	{regex: knockoutOnCoinFlipRegex, handler: parseKnockoutOnCoinFlip},
	{regex: passiveRestrictionEvolveRegex, handler: parsePassiveRestrictionEvolve},
	{regex: scalingDamageBenchedBaseRegex, handler: parseScalingDamageBenchedBase},
	{regex: scalingDamageAllOpponentEnergyRegex, handler: parseScalingDamageAllOpponentEnergy},
	{regex: moveEnergyOnKORegex, handler: parseMoveEnergyOnKO},
	{regex: recoilDamageOnKORegex, handler: parseRecoilDamageOnKO},
	{regex: revealHandRegex, handler: parseRevealHand},
	{regex: moveDamageRegex, handler: parseMoveDamage},
	{regex: discardAllToolsRegex, handler: parseDiscardAllTools},
	{regex: passivePreventionFromEXRegex, handler: parsePassivePreventionFromEX},
	{regex: discardRandomEnergySelfMultipleRegex, handler: parseDiscardRandomEnergySelfMultiple},
	{regex: splashDamageAnyFriendlyRegex, handler: parseSplashDamageAnyFriendly},
	{regex: conditionalRestrictionOnStageRegex, handler: parseConditionalRestrictionOnStage},
	{regex: attachEnergyScaledByCoinFlipsRegex, handler: parseAttachEnergyScaledByCoinFlips},
	{regex: attachMultipleSpecificEnergyRegex, handler: parseAttachMultipleSpecificEnergy},
	{regex: conditionalDamageOpponentHasStatusRegex, handler: parseConditionalDamageOpponentHasStatus},
	{regex: applyAttackFailureChanceRegex, handler: parseApplyAttackFailureChance},
	{regex: discardFromHandOnCoinFlipRegex, handler: parseDiscardFromHandOnCoinFlip},
	{regex: moveEnergyAsOftenAsYouLikeRegex, handler: parseMoveEnergyAsOftenAsYouLike},
	{regex: conditionalRestrictionOnCoinFlipRegex, handler: parseConditionalRestrictionOnCoinFlip},
	{regex: switchSelfSubtypeRegex, handler: parseSwitchSelfSubtype},
	{regex: attachEnergyFromDiscardWithRecoilRegex, handler: parseAttachEnergyFromDiscardWithRecoil},
	{regex: passiveEvolveToAnyRegex, handler: parsePassiveEvolveToAny},
	{regex: passiveCostReductionWithToolRegex, handler: parsePassiveCostReductionWithTool},
	{regex: copyAttackWithEnergyCheckRegex, handler: parseCopyAttackWithEnergyCheck},
	{regex: passiveIncreaseOpponentCostRegex, handler: parsePassiveIncreaseOpponentCost},
	{regex: attachEnergyToTypedBenchedRegex, handler: parseAttachEnergyToTypedBenched},
	{regex: conditionalDamageOpponentStageRegex, handler: parseConditionalDamageOpponentStage},
	{regex: lookAtEitherPlayerDeckRegex, handler: parseLookAtEitherPlayerDeck},
	{regex: persistentAttackFailureRegex, handler: parsePersistentAttackFailure},
	{regex: healOnCoinFlipRegex, handler: parseHealOnCoinFlip},
	{regex: discardOwnDeckAmountRegex, handler: parseDiscardOwnDeckAmount},
	{regex: discardDeckWithConditionalDamageRegex, handler: parseDiscardDeckWithConditionalDamage},
	{regex: conditionalDamageOnBenchedNameRegex, handler: parseConditionalDamageOnBenchedName},
	{regex: damageHalveHPRoundedDownRegex, handler: parseDamageHalveHPRoundedDown},
	{regex: passiveGlobalDamageReductionUnownRegex, handler: parsePassiveGlobalDamageReductionUnown},
	{regex: applyRandomStatusRegex, handler: parseApplyRandomStatus},
	{regex: scalingDamageDiscardToolRegex, handler: parseScalingDamageDiscardTool},
	{regex: searchDeckToolRegex, handler: parseSearchDeckTool},
	{regex: attachEnergyAtEndOfFirstTurnRegex, handler: parseAttachEnergyAtEndOfFirstTurn},
	{regex: evolveOnEnergyAttachRegex, handler: parseEvolveOnEnergyAttach},
	{regex: healActiveOncePerTurnRegex, handler: parseHealActiveOncePerTurn},
	{regex: attachEnergyToAnyTypedFriendlyRegex, handler: parseAttachEnergyToAnyTypedFriendly},
	{regex: passiveZeroRetreatInPlayRegex, handler: parsePassiveZeroRetreatInPlay},
	{regex: scalingDamageSelfEnergyAllTypesRegex, handler: parseScalingDamageSelfEnergyAllTypes},
	{regex: returnToHandOnCoinFlipRegex, handler: parseReturnToHandOnCoinFlip},
	{regex: discardBenchedForScalingDamageRegex, handler: parseDiscardBenchedForScalingDamage},
	{regex: passiveGlobalHealBlockRegex, handler: parsePassiveGlobalHealBlock},
	{regex: devolveOnConditionRegex, handler: parseDevolveOnCondition},
	{regex: shuffleHandAndDrawScaledRegex, handler: parseShuffleHandAndDrawScaled},
	{regex: discardEnergyBothActiveRegex, handler: parseDiscardEnergyBothActive},
	{regex: conditionalDamageBenchedDamagedRegex, handler: parseConditionalDamageBenchedDamaged},
	{regex: shuffleFromHandOnCoinFlipRegex, handler: parseShuffleFromHandOnCoinFlip},
	{regex: discardEnergyOpponentOnCoinFlipRegex, handler: parseDiscardEnergyOpponentOnCoinFlip},
	{regex: drawUntilMatchHandSizeRegex, handler: parseDrawUntilMatchHandSize},
	{regex: scalingDamageOpponentBenchedCountRegex, handler: parseScalingDamageOpponentBenchedCount},
	{regex: copyAttackOnCoinFlipRegex, handler: parseCopyAttackOnCoinFlip},
	{regex: snipeRandomOpponentRegex, handler: parseSnipeRandomOpponent},
	{regex: passiveGlobalDamageBuffUnownRegex, handler: parsePassiveGlobalDamageBuffUnown},
	{regex: preventionOnKORegex, handler: parsePreventionOnKO},
	{regex: moveAllEnergyToBenchedRegex, handler: parseMoveAllEnergyToBenched},
	{regex: restrictEnergyAttachmentRegex, handler: parseRestrictEnergyAttachment},
	{regex: shuffleFromHandSimpleRegex, handler: parseShuffleFromHandSimple},
	{regex: voluntarySwitchRegex, handler: parseVoluntarySwitch},
	{regex: revealHandOnBenchPlayRegex, handler: parseRevealHandOnBenchPlay},
	{regex: debuffIncomingDamageRegex, handler: parseDebuffIncomingDamage},
	{regex: restrictionOnTailsRegex, handler: parseRestrictionOnTails},
	{regex: conditionalDamageOpponentNameRegex, handler: parseConditionalDamageOpponentName},
	{regex: attachEnergyToSpecificPokemonRegex, handler: parseAttachEnergyToSpecificPokemon},
	{regex: passiveRetreatCostFirstTurnRegex, handler: parsePassiveRetreatCostFirstTurn},
	{regex: passiveEffectPreventionRegex, handler: parsePassiveEffectPrevention},
	{regex: applyStatusOncePerTurnAbilityRegex, handler: parseApplyStatusOncePerTurnAbility},
	{regex: buffStackingRegex, handler: parseBuffStacking},
	{regex: conditionalDamageIfEnergyAttachedRegex, handler: parseConditionalDamageIfEnergyAttached},
	{regex: healOnEnergyAttachRegex, handler: parseHealOnEnergyAttach},
	{regex: knockoutAttackerOnKORegex, handler: parseKnockoutAttackerOnKO},
	{regex: damageAbilityInPlayRegex, handler: parseDamageAbilityInPlay},
	{regex: finalConditionalDamageMultiCoinFlipRegex, handler: parseFinalConditionalDamageMultiCoinFlip},
	{regex: finalPassiveDamageReductionRegex, handler: parseFinalPassiveDamageReduction},
	{regex: finalDiscardEnergyOpponentSimpleRegex, handler: parseFinalDiscardEnergyOpponentSimple},
	{regex: finalConditionalDamageUntilTailsRegex, handler: parseFinalConditionalDamageUntilTails},
	{regex: finalScalingDamageOpponentEnergyBaseRegex, handler: parseFinalScalingDamageOpponentEnergyBase},
	{regex: finalAttachEnergyToActiveTypedRegex, handler: parseFinalAttachEnergyToActiveTyped},
	{regex: finalConditionalDamageOpponentIsEXRegex, handler: parseFinalConditionalDamageOpponentIsEX},
	{regex: finalDiscardSelfEnergyOnTailsRegex, handler: parseFinalDiscardSelfEnergyOnTails},
	{regex: finalReduceDamageNextTurnRegex, handler: parseFinalReduceDamageNextTurn},
	{regex: finalForceSwitchOnHeadsRegex, handler: parseFinalForceSwitchOnHeads},
	{regex: finalIncreaseOpponentCostNextTurnRegex, handler: parseFinalIncreaseOpponentCostNextTurn},
	{regex: finalPreventionOnHeadsRegex, handler: parseFinalPreventionOnHeads},
	{regex: finalSearchDeckGenericPokemonRegex, handler: parseFinalSearchDeckGenericPokemon},
	{regex: finalForceSwitchSimpleRegex, handler: parseFinalForceSwitchSimple},
	{regex: hariyamaPushOutRegex, handler: parseHariyamaPushOut},
}

// Parse iterates through a list of registered effect parsers and returns
// a structured Effect object from the first one that matches the input text.
func Parse(text string) []core.Effect {
	text = strings.TrimSpace(text)

	for _, p := range effectParsers {
		if matches := p.regex.FindStringSubmatch(text); len(matches) > 0 {
			if effects := p.handler(matches, text); effects != nil {
				return effects
			}
		}
	}

	// If no other rule matches, return an UNKNOWN effect type.
	return []core.Effect{{
		Type:        core.EffectUnknown,
		Description: text,
	}}
}

// --- Handler Functions ---

func parseHealAllFriendly(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Target:      core.TargetAllFriendly,
		Amount:      amount,
		Description: text},
	}
}
func parseHeal(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text},
	}
}
func parseCantAttackNextTurn(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectRestrictionCantAttack,
		Target:      core.TargetSelf,
		Description: text,
		Conditions: []core.Condition{
			core.DurationCondition{
				Duration: core.DurationNextTurn,
			},
		},
	}}
}
func parseForceSwitch(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectForceSwitch,
		Target:      core.TargetOpponentActive,
		Description: text,
	}}
}
func parseRecoilDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectRecoilDamage,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text},
	}
}
func parseSearchDeck(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err1 := strconv.Atoi(matches[1])
	pokemonType := matches[2]
	if err1 == nil {
		return []core.Effect{{
			Type:        core.EffectSearchDeck,
			Target:      core.TargetDeck,
			Amount:      amount,
			Description: text,
			Conditions: []core.Condition{
				core.SearchCondition{
					PokemonType: core.PokemonType(pokemonType),
					Random:      true,
					Target:      core.TargetHand,
				},
			},
		}}
	}
	return nil
}
func parseConditionalDamage(matches []string, text string) []core.Effect {
	if len(matches) < 4 {
		return nil
	}
	requiredEnergyCount, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	energyType := matches[2]
	damageBonus, err2 := strconv.Atoi(matches[3])
	if err2 != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      damageBonus,
		Description: text,
		Conditions: []core.Condition{
			core.EnergyCondition{
				RequiredExtraEnergyCount: requiredEnergyCount,
				RequiredEnergyType:       energyType,
			},
		}},
	}
}
func parseApplyStatusOpponent(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseApplyStatusSelf(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectApplyStatus,
		Target:      core.TargetSelf,
		Status:      core.StatusCondition(strings.ToUpper(matches[1])),
		Description: text,
	}}
}
func parseAttachEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectAttachEnergy,
		Target:      core.TargetSelf,
		Description: text,
		Conditions: []core.Condition{
			core.EnergyAttachCondition{
				RequiredEnergyType: matches[1],
			},
		},
	}}
}
func parseSearchEvolution(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSearchDeck,
		Target:      core.TargetDeck,
		Amount:      1,
		Description: text,
		Conditions: []core.Condition{
			core.SearchCondition{
				PokemonType: core.PokemonType("EVOLVES_FROM"),
				Random:      true,
				Target:      core.TargetHand,
			},
		},
	}}
}
func parseTriggeredSleep(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectTriggeredAbility,
		Target:      core.TargetSelf,
		Status:      core.StatusAsleep,
		Description: text,
		Conditions: []core.Condition{
			core.TriggerCondition{
				Trigger: core.TriggerAttachEnergySelf,
			},
		},
	}}
}
func parseConditionalDamageHPRatio(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: []core.Condition{
			core.TriggerCondition{
				Trigger: core.TriggerOponentHpGreater,
			},
		}},
	}
}
func parseConditionalDamageEvolvedTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "EVOLVED_THIS_TURN",
		}},
	}
}
func parsePassiveDamageCheckup(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveDamage,
		Target:      core.TargetOpponentActive,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"phase":    "CHECKUP",
			"location": "ACTIVE",
		}},
	}
}
func parseScalingDamageRetreatCost(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "OPPONENT_RETREAT_COST",
		}},
	}
}
func parseCoinFlipTailsFails(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectAttackMayFail,
		Description: text,
		Conditions: map[string]interface{}{
			"chance": 0.5,
			"on":     "TAILS",
		},
	}}
}
func parsePassiveRetreatCostLatias(matches []string, text string) []core.Effect {
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
func parseDiscardAllEnergy(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDiscardEnergy,
		Target:      core.TargetSelf,
		Description: text,
		Conditions: map[string]interface{}{
			"amount": "ALL",
		},
	}}
}
func parseDiscardTypedEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDiscardEnergy,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"energyType": matches[2],
		}},
	}
}
func parseMoveAllTypedEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseReduceIncomingDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectReduceIncomingDamage,
		Target:      core.TargetOpponentActive,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"duration": "next_turn",
		}},
	}
}
func parseDiscardFromHand(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseCantRetreat(matches []string, text string) []core.Effect {
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
func parseMultiHitRandom(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	hits, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	damage, err2 := strconv.Atoi(matches[2])
	if err2 != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectMultiHitRandomDamage,
		Target:      core.TargetOpponentActive, // Target pool is opponent's field
		Amount:      damage,
		Description: text,
		Conditions: map[string]interface{}{
			"hits": hits,
		}},
	}
}
func parseStatusOnCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseScalingDamageSelfDamage(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "SELF_DAMAGE_COUNTERS",
		},
	}}
}
func parseSplashDamageBenchedFriendly(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDamageBenchedFriendly,
		Target:      core.TargetBenchedFriendly,
		Amount:      amount,
		Description: text},
	}
}
func parseConditionalDamageExistingDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "OPPONENT_HAS_DAMAGE",
		}},
	}
}
func parseConditionalDamageToolAttached(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "OPPONENT_HAS_TOOL",
		}},
	}
}
func parseGutsAbility(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	hp, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Target:      core.TargetSelf,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":       "PREVENT_KNOCKOUT",
			"on_coin_flip": "HEADS",
			"remaining_hp": hp,
		}},
	}
}
func parseSearchDeckByName(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseSnipeAndDiscardTyped(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	energyType := matches[1]
	damage, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	discardEffect := core.Effect{
		Type:        core.EffectDiscardEnergy,
		Target:      core.TargetSelf,
		Description: "Discard all {L} Energy from this Pokémon.",
		Conditions: map[string]interface{}{
			"amount":     "ALL",
			"energyType": energyType,
		}}
	snipeEffect := core.Effect{
		Type:        core.EffectSnipeDamage,
		Target:      core.TargetBenchedOpponent, // Can target any opponent's Pokémon
		Amount:      damage,
		Description: "This attack does 120 damage to 1 of your opponent's Pokémon.",
	}
	return []core.Effect{discardEffect, snipeEffect}
}
func parseDiscardRandomEnergySelf(matches []string, text string) []core.Effect {
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
func parseHealOnEvolve(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parseForceSwitchOncePerTurn(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectForceSwitch,
		Target:      core.TargetOpponentActive,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "ONCE_PER_TURN",
		},
	}}
}
func parseReduceDamageNextTurnSelf(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectReduceIncomingDamage,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"duration": "opponent_next_turn",
		}},
	}
}
func parseCopyAttack(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectCopyAttack,
		Target:      core.TargetOpponentActive,
		Description: text,
	}}
}
func parsePassiveDamageReductionTyped(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parseDrawAtEndOfTurn(matches []string, text string) []core.Effect {
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
func parseApplyStatusOncePerTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseDiscardRandomEnergyGlobal(matches []string, text string) []core.Effect {
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
func parseSwitchSelf(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectSwitchSelf,
		Target:      core.TargetBenchedFriendly,
		Description: text,
	}}
}
func parseSnipeDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSnipeDamage,
		Target:      core.TargetBenchedOpponent, // User chooses which one
		Amount:      amount,
		Description: text},
	}
}
func parseHealOncePerTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":  "ONCE_PER_TURN",
			"location": "ACTIVE",
		}},
	}
}
func parseAttachEnergyEndsTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseDiscardEnergyOnEvolve(matches []string, text string) []core.Effect {
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
func parseScalingDamageSelfEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parseDrawCard(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDraw,
		Amount:      1,
		Description: text,
	}}
}
func parseAttachEnergyMultiBenched(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parseConditionalDamageSupporter(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "PLAYED_SUPPORTER_THIS_TURN",
		}},
	}
}
func parseDrawOnEvolve(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDraw,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "ON_EVOLVE",
		}},
	}
}
func parseDiscardSingleTypedEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseShuffleIntoDeckOnCoinFlip(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectShuffleIntoDeck,
		Target:      core.TargetOpponentActive,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "HEADS",
		},
	}}
}
func parseScalingDamageOpponentEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "OPPONENT_ATTACHED_ENERGY",
		}},
	}
}
func parsePreventAllDamageOnCoinFlip(matches []string, text string) []core.Effect {
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
func parseIncreaseOpponentCosts(matches []string, text string) []core.Effect {
	if len(matches) < 5 {
		return nil
	}
	attackCost, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	retreatCost, err2 := strconv.Atoi(matches[3])
	if err2 != nil {
		return nil
	}
	attackCostEffect := core.Effect{
		Type:        core.EffectApplyRestriction,
		Target:      core.TargetOpponentActive,
		Description: "Increase opponent's attack cost",
		Conditions: map[string]interface{}{
			"restriction": "INCREASE_ATTACK_COST",
			"amount":      attackCost,
			"energyType":  matches[2],
			"duration":    "opponent_next_turn",
		}}
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
func parseSnipeDamagedPokemon(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSnipeDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"target_condition": "HAS_DAMAGE",
		}},
	}
}
func parseSwitchSelfFromBench(matches []string, text string) []core.Effect {
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
func parseConditionalDamageSelfTool(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "SELF_HAS_TOOL",
		}},
	}
}
func parseAttachMultipleEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parsePassiveCostReductionInPlay(matches []string, text string) []core.Effect {
	if len(matches) < 4 {
		return nil
	}
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
	return nil
}
func parseConditionalDamageOpponentAbility(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "OPPONENT_HAS_ABILITY",
		}},
	}
}
func parseScalingSnipeOpponentEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingSnipeDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "TARGET_ATTACHED_ENERGY",
		}},
	}
}
func parseSnipeOncePerTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSnipeDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "ONCE_PER_TURN",
		}},
	}
}
func parseScalingDamageBenchedName(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parseDamageBenchedOpponentAll(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDamageBenchedOpponentAll,
		Target:      core.TargetBenchedOpponentAll,
		Amount:      amount,
		Description: text},
	}
}
func parseLifesteal(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectLifesteal,
		Target:      core.TargetSelf,
		Description: text,
	}}
}
func parsePassiveEnergyValue(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parseScalingDamageMultiCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	flips, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	amount, err2 := strconv.Atoi(matches[2])
	if err2 != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":  "COIN_FLIP_HEADS",
			"num_flips": flips,
		}},
	}
}
func parseConditionalDamageSwitchIn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "SWITCHED_IN_THIS_TURN",
		}},
	}
}
func parseConditionalDamageCoinFlipScaling(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "COIN_FLIP_HEADS",
		}},
	}
}
func parsePassiveImmunity(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "IMMUNE_TO_SPECIAL_CONDITIONS",
		},
	}}
}
func parseScalingDamageBenchedCount(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "BENCHED_POKEMON_COUNT",
		}},
	}
}
func parseReactiveDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectApplyReactiveDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"duration": "opponent_next_turn",
		}},
	}
}
func parseAttachEnergyToActiveTyped(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parseBuffNextTurn(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectBuffNextTurn,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"attack_name": matches[1],
		}},
	}
}
func parseAttachEnergyToBenched(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseScalingDamageBenchedType(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":      "BENCHED_POKEMON_TYPE",
			"scale_by_type": "EVOLUTION",
		}},
	}
}
func parsePassiveRestriction(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseModifyEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseTriggeredDamageOnEnergyAttach(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveDamage,
		Target:      core.TargetOpponentActive,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":     "ATTACH_ENERGY_TO_SELF",
			"energy_type": matches[1],
		}},
	}
}
func parseSwitchSelfTyped(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSwitchSelf,
		Target:      core.TargetBenchedFriendly,
		Description: text,
		Conditions: map[string]interface{}{
			"target_type": matches[1],
		},
	}}
}
func parseConditionalDamageOpponentProperty(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":  "OPPONENT_HAS_PROPERTY",
			"property": strings.ToUpper(matches[1]),
		}},
	}
}
func parseConditionalDamageOnKO(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "FRIENDLY_KO_LAST_TURN",
		}},
	}
}
func parseSplashDamageBenchedFriendlyAll(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDamageBenchedFriendly,
		Target:      core.TargetBenchedFriendly,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"target_all": true,
		}},
	}
}
func parseScalingDamageBenchedNames(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parseDiscardMultipleTypedEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 4 {
		return nil
	}
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
func parseAttachEnergyToBenchedStage(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parsePassiveDamageBuffInPlay(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	names := strings.Split(matches[1], " or ")
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":           "BUFF_DAMAGE_OUTPUT",
			"amount":           amount,
			"requires_in_play": names,
		}},
	}
}
func parsePassiveDamagePreventionCoinFlip(matches []string, text string) []core.Effect {
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
func parseHealAtEndOfTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":  "END_OF_TURN",
			"location": "ACTIVE",
		}},
	}
}
func parsePassiveRetreatCostEnergy(matches []string, text string) []core.Effect {
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
func parseDamageAllOpponent(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDamageAllOpponent,
		Amount:      amount,
		Description: text},
	}
}
func parsePassiveRetreatCostReductionBench(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":       "REDUCE_RETREAT_COST",
			"location":     "BENCH",
			"target":       "ACTIVE",
			"target_stage": matches[1],
			"amount":       amount,
		}},
	}
}
func parseDamageBenchedConditionalEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDamageBenchedOpponentAll,
		Target:      core.TargetBenchedOpponentAll,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"target_condition": "HAS_ENERGY_ATTACHED",
		}},
	}
}
func parseScalingDamagePerEnergyCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":            "COIN_FLIP_HEADS",
			"num_flips_scales_by": "SELF_ATTACHED_ENERGY",
		}},
	}
}
func parseAlternateAttackCost(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":      "ALTERNATE_ATTACK_COST",
			"trigger":     "SELF_HAS_DAMAGE",
			"cost_amount": amount,
			"cost_type":   matches[2],
		}},
	}
}
func parseDiscardAllTypedEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseConditionalDamageOpponentStatus(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "OPPONENT_HAS_STATUS",
			"status":  strings.ToUpper(matches[1]),
		}},
	}
}
func parsePreventAllDamageSimple(matches []string, text string) []core.Effect {
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
func parseScalingDamageBenchedTypeCount(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":      "BENCHED_POKEMON_TYPE_COUNT",
			"scale_by_type": matches[2],
		}},
	}
}
func parsePassiveDamageReductionSimple(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Target:      core.TargetSelf,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "REDUCE_INCOMING_DAMAGE",
			"amount": amount,
		}},
	}
}
func parseConditionalDamageNoDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "SELF_HAS_NO_DAMAGE",
		}},
	}
}
func parsePassiveRetreatCostForOther(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseRecoilDamageOnCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectRecoilDamage,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "TAILS",
		}},
	}
}
func parsePassiveImmunitySingleStatus(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "IMMUNE_TO_STATUS",
			"status": strings.ToUpper(matches[1]),
		},
	}}
}
func parseDiscardDeck(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDiscardDeck,
		Target:      core.TargetOpponentActive, // Implied target is opponent
		Amount:      1,
		Description: text,
	}}
}
func parseSplashDamageSingleBenchedOpponent(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSnipeDamage,
		Target:      core.TargetBenchedOpponent,
		Amount:      amount,
		Description: text},
	}
}
func parseScalingDamageAllBenched(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "ALL_BENCHED_POKEMON_COUNT",
		}},
	}
}
func parseMultiHitRandomGlobal(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	hits, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	damage, err2 := strconv.Atoi(matches[2])
	if err2 != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectMultiHitRandomDamage,
		Amount:      damage,
		Description: text,
		Conditions: map[string]interface{}{
			"hits":        hits,
			"target_pool": "GLOBAL_OTHER",
		}},
	}
}
func parsePassiveSpecialConditionImmunityTypedEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseHealAllFriendlyTyped(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":     "ONCE_PER_TURN",
			"target_all":  true,
			"target_type": matches[2],
		}},
	}
}
func parseScalingDamagePerTypedEnergyCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":            "COIN_FLIP_HEADS",
			"num_flips_scales_by": "SELF_ATTACHED_ENERGY_TYPED",
			"energy_type":         matches[1],
		}},
	}
}
func parseHealBenched(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Target:      core.TargetBenchedFriendly,
		Amount:      amount,
		Description: text},
	}
}
func parseConditionalDamageDifferentEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	count, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	amount, err2 := strconv.Atoi(matches[2])
	if err2 != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":        "DIFFERENT_ENERGY_TYPES_ATTACHED",
			"required_count": count,
		}},
	}
}
func parseDamageEqualsSelfDamage(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDamage, // This is a direct damage effect, not a bonus
		Description: text,
		Conditions: map[string]interface{}{
			"amount_equals": "SELF_DAMAGE_COUNTERS",
		},
	}}
}
func parseSetHPOnCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	hp, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSetHP,
		Target:      core.TargetOpponentActive,
		Amount:      hp,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "HEADS",
		}},
	}
}
func parseShuffleFromHandMultiCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	flips, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectShuffleFromHand,
		Target:      core.TargetOpponentHand,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":  "COIN_FLIP_HEADS",
			"num_flips": flips,
			"random":    true,
		}},
	}
}
func parsePassiveReactiveDamageActive(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":   "REACTIVE_DAMAGE",
			"amount":   amount,
			"location": "ACTIVE",
		}},
	}
}
func parseShuffleFromHandReveal(matches []string, text string) []core.Effect {
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
func parseRestrictionCantUseAttack(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseScalingDamageUntilTails(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "COIN_FLIP_HEADS_UNTIL_TAILS",
		}},
	}
}
func parseForceSwitchDamaged(matches []string, text string) []core.Effect {
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
func parseConditionalDamageDoubleHeads(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "DOUBLE_HEADS",
		}},
	}
}
func parseScalingDamagePerPokemonInPlay(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":            "COIN_FLIP_HEADS",
			"num_flips_scales_by": "ALL_POKEMON_IN_PLAY",
		}},
	}
}
func parseModifyNextEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parsePassiveOpponentDamageReduction(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":   "REDUCE_OPPONENT_DAMAGE_OUTPUT",
			"amount":   amount,
			"location": "ACTIVE",
		}},
	}
}
func parseConditionalDamageIfDamagedLastTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "DAMAGED_LAST_TURN",
		}},
	}
}
func parseLookAtTopCard(matches []string, text string) []core.Effect {
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
func parsePassiveBuffStatusDamage(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "BUFF_STATUS_DAMAGE",
			"amount": amount,
			"status": strings.ToUpper(matches[2]),
		}},
	}
}
func parseAttachEnergyOncePerTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseSearchDeckRandomPokemon(matches []string, text string) []core.Effect {
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
func parseConditionalDamageSelfHasDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "SELF_HAS_DAMAGE",
		}},
	}
}
func parseConditionalDamageOnAttackHistory(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":     "ATTACK_USED_LAST_TURN",
			"attack_name": matches[1],
		}},
	}
}
func parseDiscardEnergyUntilTails(matches []string, text string) []core.Effect {
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
func parseDelayedDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDelayedDamage,
		Target:      core.TargetOpponentActive,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "END_OF_OPPONENT_NEXT_TURN",
		}},
	}
}
func parsePassiveReactiveDamageOnKO(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":   "REACTIVE_DAMAGE_ON_KO",
			"amount":   amount,
			"location": "ACTIVE",
		}},
	}
}
func parseRestrictOpponentHandNextTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parsePassiveZeroRetreatForActive(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "ZERO_RETREAT_COST",
			"target": "ACTIVE",
		},
	}}
}
func parseScalingDamageAttackHistory(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":    "ATTACK_USAGE_COUNT",
			"attack_name": matches[2],
		}},
	}
}
func parseDiscardDeckBothPlayers(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDiscardDeck,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"target": "BOTH_PLAYERS",
		}},
	}
}
func parsePassiveDamageReductionOnCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":       "REDUCE_INCOMING_DAMAGE",
			"amount":       amount,
			"on_coin_flip": "HEADS",
		}},
	}
}
func parseMoveEnergyBenchedToActive(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parseRestrictOpponentHandItem(matches []string, text string) []core.Effect {
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
func parsePassiveDamageBuffEvolvesFrom(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":              "BUFF_DAMAGE_OUTPUT",
			"amount":              amount,
			"location":            "BENCH",
			"target_evolves_from": matches[1],
		}},
	}
}
func parseForceSwitchBenchedBasic(matches []string, text string) []core.Effect {
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
func parseApplyStatusBothActive(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseApplyRestrictionCantAttack(matches []string, text string) []core.Effect {
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
func parsePassiveDamageReductionInPlay(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	names := strings.Split(matches[1], " or ")
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":           "REDUCE_INCOMING_DAMAGE",
			"amount":           amount,
			"requires_in_play": names,
		}},
	}
}
func parsePassiveDamageBuffTypedPokemon(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":      "BUFF_DAMAGE_OUTPUT",
			"amount":      amount,
			"target_type": matches[1],
		}},
	}
}
func parseDrawWithDiscardCost(matches []string, text string) []core.Effect {
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
func parseKnockoutOnCoinFlip(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectKnockout,
		Target:      core.TargetOpponentActive,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "DOUBLE_HEADS",
		},
	}}
}
func parsePassiveRestrictionEvolve(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "RESTRICT_OPPONENT_EVOLVE",
			"target": "ACTIVE",
		},
	}}
}
func parseScalingDamageBenchedBase(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"is_base_damage": true, // Differentiates from bonus damage
			"scale_by":       "BENCHED_POKEMON_COUNT",
		}},
	}
}
func parseScalingDamageAllOpponentEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"is_base_damage": true,
			"scale_by":       "ALL_OPPONENT_POKEMON_ENERGY",
		}},
	}
}
func parseMoveEnergyOnKO(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseRecoilDamageOnKO(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectRecoilDamage,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "OPPONENT_KO",
		}},
	}
}
func parseRevealHand(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectRevealHand,
		Target:      core.TargetOpponentHand,
		Description: text,
	}}
}
func parseMoveDamage(matches []string, text string) []core.Effect {
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
func parseDiscardAllTools(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDiscardTool,
		Target:      core.TargetOpponentActive,
		Amount:      99, // Represents "all"
		Description: text,
	}}
}
func parsePassivePreventionFromEX(matches []string, text string) []core.Effect {
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
func parseDiscardRandomEnergySelfMultiple(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDiscardEnergy,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"random": true,
		}},
	}
}
func parseSplashDamageAnyFriendly(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDamageBenchedFriendly, // Re-using this type, but target pool is wider
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"target_pool": "ANY_FRIENDLY",
		}},
	}
}
func parseConditionalRestrictionOnStage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseAttachEnergyScaledByCoinFlips(matches []string, text string) []core.Effect {
	if len(matches) < 4 {
		return nil
	}
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
	return nil
}
func parseAttachMultipleSpecificEnergy(matches []string, text string) []core.Effect {
	if len(matches) < 4 {
		return nil
	}
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
func parseConditionalDamageOpponentHasStatus(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "OPPONENT_HAS_SPECIAL_CONDITION",
		}},
	}
}
func parseApplyAttackFailureChance(matches []string, text string) []core.Effect {
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
func parseDiscardFromHandOnCoinFlip(matches []string, text string) []core.Effect {
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
func parseMoveEnergyAsOftenAsYouLike(matches []string, text string) []core.Effect {
	if len(matches) < 4 {
		return nil
	}
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
func parseConditionalRestrictionOnCoinFlip(matches []string, text string) []core.Effect {
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
func parseSwitchSelfSubtype(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parseAttachEnergyFromDiscardWithRecoil(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	attachEffect := core.Effect{
		Type:        core.EffectAttachEnergy,
		Description: "Attach energy from discard.",
		Conditions: map[string]interface{}{
			"trigger":    "ONCE_PER_TURN",
			"source":     "DISCARD_PILE",
			"energyType": matches[1],
		}}
	recoilEffect := core.Effect{
		Type:        core.EffectRecoilDamage,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: "Take recoil damage.",
	}
	return []core.Effect{attachEffect, recoilEffect}
}
func parsePassiveEvolveToAny(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "CAN_EVOLVE_INTO_ANY",
		},
	}}
}
func parsePassiveCostReductionWithTool(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":     "REDUCE_ATTACK_COST",
			"trigger":    "SELF_HAS_TOOL",
			"amount":     amount,
			"energyType": matches[2],
		}},
	}
}
func parseCopyAttackWithEnergyCheck(matches []string, text string) []core.Effect {
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
func parsePassiveIncreaseOpponentCost(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":     "INCREASE_OPPONENT_ATTACK_COST",
			"location":   "ACTIVE",
			"amount":     amount,
			"energyType": matches[2],
		}},
	}
}
func parseAttachEnergyToTypedBenched(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parseConditionalDamageOpponentStage(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":        "OPPONENT_IS_STAGE",
			"opponent_stage": strings.ToUpper(matches[1]),
		}},
	}
}
func parseLookAtEitherPlayerDeck(matches []string, text string) []core.Effect {
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
func parsePersistentAttackFailure(matches []string, text string) []core.Effect {
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
func parseHealOnCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "HEADS",
		}},
	}
}
func parseDiscardOwnDeckAmount(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDiscardDeck,
		Target:      core.TargetDeck,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"target_player": "SELF",
		}},
	}
}
func parseDiscardDeckWithConditionalDamage(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	discardEffect := core.Effect{
		Type:        core.EffectDiscardDeck,
		Target:      core.TargetDeck,
		Amount:      1,
		Description: "Discard top card of your deck.",
		Conditions: map[string]interface{}{
			"target_player": "SELF",
		}}
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
func parseConditionalDamageOnBenchedName(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":      "POKEMON_ON_BENCH",
			"pokemon_name": matches[1],
		}},
	}
}
func parseDamageHalveHPRoundedDown(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDamageHalveHP,
		Target:      core.TargetOpponentActive,
		Description: text,
	}}
}
func parsePassiveGlobalDamageReductionUnown(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":           "REDUCE_INCOMING_DAMAGE",
			"target":           "ALL_FRIENDLY",
			"amount":           amount,
			"requires_in_play": []string{"Unown"}, // Special condition for this ability,
		}},
	}
}
func parseApplyRandomStatus(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseScalingDamageDiscardTool(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	maxDiscard, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	damagePer, err2 := strconv.Atoi(matches[2])
	if err2 != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      damagePer,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":    "DISCARD_TOOL_FROM_HAND",
			"max_discard": maxDiscard,
		}},
	}
}
func parseSearchDeckTool(matches []string, text string) []core.Effect {
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
func parseAttachEnergyAtEndOfFirstTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseEvolveOnEnergyAttach(matches []string, text string) []core.Effect {
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
func parseHealActiveOncePerTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Target:      core.TargetSelf, // Target is Active, which is self in this context
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "ONCE_PER_TURN",
		}},
	}
}
func parseAttachEnergyToAnyTypedFriendly(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parsePassiveZeroRetreatInPlay(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseScalingDamageSelfEnergyAllTypes(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "SELF_ATTACHED_ENERGY",
		}},
	}
}
func parseReturnToHandOnCoinFlip(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectReturnToHand,
		Target:      core.TargetOpponentActive,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "HEADS",
		},
	}}
}
func parseDiscardBenchedForScalingDamage(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDiscardBenched,
		Target:      core.TargetBenchedFriendly, // Player chooses which/how many
		Description: "You may discard any number of your Benched {W} Pokémon.",
		Conditions: map[string]interface{}{
			"target_type": matches[1],
		}},
		{
			Type:        core.EffectScalingDamage,
			Amount:      amount,
			Description: "This attack does 40 more damage for each Benched Pokémon you discarded in this way.",
			Conditions: map[string]interface{}{
				"scale_by": "DISCARDED_BENCHED_COUNT",
			},
		}}
}
func parsePassiveGlobalHealBlock(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "PREVENT_HEALING",
			"target": "GLOBAL",
		},
	}}
}
func parseDevolveOnCondition(matches []string, text string) []core.Effect {
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
func parseShuffleHandAndDrawScaled(matches []string, text string) []core.Effect {
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
func parseDiscardEnergyBothActive(matches []string, text string) []core.Effect {
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
func parseConditionalDamageBenchedDamaged(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "ANY_BENCHED_FRIENDLY_HAS_DAMAGE",
		}},
	}
}
func parseShuffleFromHandOnCoinFlip(matches []string, text string) []core.Effect {
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
func parseDiscardEnergyOpponentOnCoinFlip(matches []string, text string) []core.Effect {
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
func parseDrawUntilMatchHandSize(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDraw,
		Description: text,
		Conditions: map[string]interface{}{
			"draw_until": "MATCH_OPPONENT_HAND_SIZE",
		},
	}}
}
func parseScalingDamageOpponentBenchedCount(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by": "OPPONENT_BENCHED_POKEMON_COUNT",
		}},
	}
}
func parseCopyAttackOnCoinFlip(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectCopyAttack,
		Target:      core.TargetOpponentActive,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "HEADS",
		},
	}}
}
func parseSnipeRandomOpponent(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectSnipeDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"target_pool": "ANY_OPPONENT",
			"random":      true,
		}},
	}
}
func parsePassiveGlobalDamageBuffUnown(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":           "BUFF_DAMAGE_OUTPUT",
			"target":           "ALL_FRIENDLY",
			"amount":           amount,
			"requires_in_play": []string{"Unown"}, // Special condition,
		}},
	}
}
func parsePreventionOnKO(matches []string, text string) []core.Effect {
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
func parseMoveAllEnergyToBenched(matches []string, text string) []core.Effect {
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
func parseRestrictEnergyAttachment(matches []string, text string) []core.Effect {
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
func parseShuffleFromHandSimple(matches []string, text string) []core.Effect {
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
func parseVoluntarySwitch(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectSwitchSelf,
		Target:      core.TargetBenchedFriendly,
		Description: text,
		Conditions: map[string]interface{}{
			"voluntary": true,
		},
	}}
}
func parseRevealHandOnBenchPlay(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectRevealHand,
		Target:      core.TargetOpponentHand,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "ON_PLAY_TO_BENCH",
		},
	}}
}
func parseDebuffIncomingDamage(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDebuffIncomingDamage,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"duration": "opponent_next_turn",
		}},
	}
}
func parseRestrictionOnTails(matches []string, text string) []core.Effect {
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
func parseConditionalDamageOpponentName(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":       "OPPONENT_IS_NAME",
			"opponent_name": matches[1],
		}},
	}
}
func parseAttachEnergyToSpecificPokemon(matches []string, text string) []core.Effect {
	if len(matches) < 4 {
		return nil
	}
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
func parsePassiveRetreatCostFirstTurn(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":   "ZERO_RETREAT_COST",
			"duration": "FIRST_TURN",
		},
	}}
}
func parsePassiveEffectPrevention(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "PREVENT_INCOMING_EFFECTS",
		},
	}}
}
func parseApplyStatusOncePerTurnAbility(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
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
func parseBuffStacking(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectBuffNextTurn,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"attack_name": matches[1],
			"stacking":    true,
			"duration":    "PERSISTENT_ACTIVE",
		}},
	}
}
func parseConditionalDamageIfEnergyAttached(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":     "SELF_HAS_TYPED_ENERGY",
			"energy_type": matches[1],
		}},
	}
}
func parseHealOnEnergyAttach(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectHeal,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger":     "ATTACH_ENERGY_TO_SELF",
			"energy_type": matches[1],
		}},
	}
}
func parseKnockoutAttackerOnKO(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect":       "KO_ATTACKER_ON_KO",
			"on_coin_flip": "HEADS",
		},
	}}
}
func parseDamageAbilityInPlay(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
	return nil
}
func parseFinalConditionalDamageMultiCoinFlip(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	flips, err1 := strconv.Atoi(matches[1])
	if err1 != nil {
		return nil
	}
	amount, err2 := strconv.Atoi(matches[2])
	if err2 != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"scale_by":  "COIN_FLIP_HEADS",
			"num_flips": flips,
		}},
	}
}
func parseFinalPassiveDamageReduction(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectPassiveAbility,
		Description: text,
		Conditions: map[string]interface{}{
			"effect": "REDUCE_INCOMING_DAMAGE",
			"amount": amount,
		}},
	}
}
func parseFinalDiscardEnergyOpponentSimple(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectDiscardEnergy,
		Target:      core.TargetOpponentActive,
		Amount:      1,
		Description: text,
		Conditions:  map[string]interface{}{"random": true},
	}}
}
func parseFinalConditionalDamageUntilTails(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions:  map[string]interface{}{"scale_by": "COIN_FLIP_HEADS_UNTIL_TAILS"}},
	}
}
func parseFinalScalingDamageOpponentEnergyBase(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectScalingDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"is_base_damage": true,
			"scale_by":       "OPPONENT_ATTACHED_ENERGY",
		}},
	}
}
func parseFinalAttachEnergyToActiveTyped(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
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
func parseFinalConditionalDamageOpponentIsEX(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectConditionalDamage,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"trigger": "OPPONENT_IS_EX",
		}},
	}
}
func parseFinalDiscardSelfEnergyOnTails(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectDiscardEnergy,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions: map[string]interface{}{
			"on_coin_flip": "TAILS",
			"random":       true,
		}},
	}
}
func parseFinalReduceDamageNextTurn(matches []string, text string) []core.Effect {
	if len(matches) < 2 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectReduceIncomingDamage,
		Target:      core.TargetSelf,
		Amount:      amount,
		Description: text,
		Conditions:  map[string]interface{}{"duration": "opponent_next_turn"}},
	}
}
func parseFinalForceSwitchOnHeads(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectForceSwitch,
		Target:      core.TargetOpponentActive,
		Description: text,
		Conditions:  map[string]interface{}{"on_coin_flip": "HEADS"},
	}}
}
func parseFinalIncreaseOpponentCostNextTurn(matches []string, text string) []core.Effect {
	if len(matches) < 3 {
		return nil
	}
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil
	}
	return []core.Effect{{
		Type:        core.EffectApplyRestriction,
		Target:      core.TargetOpponentActive,
		Description: text,
		Conditions: map[string]interface{}{
			"restriction": "INCREASE_ATTACK_COST",
			"amount":      amount,
			"energyType":  matches[2],
			"duration":    "opponent_next_turn",
		}},
	}
}
func parseFinalPreventionOnHeads(matches []string, text string) []core.Effect {
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
func parseFinalSearchDeckGenericPokemon(matches []string, text string) []core.Effect {
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
func parseFinalForceSwitchSimple(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectForceSwitch,
		Target:      core.TargetOpponentActive,
		Description: text,
	}}
}
func parseHariyamaPushOut(matches []string, text string) []core.Effect {
	return []core.Effect{{
		Type:        core.EffectForceSwitch,
		Target:      core.TargetOpponentActive,
		Description: text,
	}}
}
