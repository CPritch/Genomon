import json
import re
import os

INPUT_FILE = "./data/finalised-cards.json"
OUTPUT_FILE = "./data/cleaned-cards.json"

def canonical_type(type_str):
    """
    Convert a type string to snake_case and map known variants to a canonical name.
    Also detect if the effect was originally conditional.
    """
    # Convert camelCase to snake_case
    s = re.sub(r'([a-z0-9])([A-Z])', r'\1_\2', type_str)
    s = s.lower()
    s = re.sub(r'_+', '_', s).strip('_')
    
    # Check for a conditional prefix
    is_conditional = False
    if s.startswith("conditional_"):
        is_conditional = True
        s = s[len("conditional_"):]
    
    # Remove underscores to use for lookup in our synonyms dictionary.
    s_key = s.replace("_", "")
    
    synonyms = {
        "attachenergy": "attach_energy",
        "discardcard": "discard",
        "discardcards": "discard",
        "damagebonus": "damage_modifier",
        "damagemodification": "damage_modifier",
        "damageprevention": "damage_reduction",
        "flipcoins": "coin_flip",
        "counterattack": "counter_attack",
        "preventattackusage": "prevent_attack",
        "selectpokemons": "select_pokemon",
        "statuseffect": "status_condition",
        "status": "status_condition",
        "switch": "switch_pokemon",
        "switchout": "switch_pokemon",
        "energyprovisionmodifier": "modify_energy_provided",
        "useattack": "use",
        "conditionaluse": "use",
        "prevent": "prevent_action",
    }
    
    canonical = synonyms.get(s_key, s)
    return canonical, is_conditional

def normalize_effect(effect):
    """
    Normalize a single effect dictionary:
      - Standardize its type using canonical_type.
      - If originally conditional, remove the 'conditional_' prefix and
        move non-base keys into a nested "condition" dict.
    """
    if "type" not in effect:
        return effect  # Nothing to normalize if there's no type

    canonical, was_conditional = canonical_type(effect["type"])
    
    # Define which keys we consider "base" for the effect and not part of its condition.
    base_keys = {"target", "amount", "damage", "resource", "destination",
                 "source", "filter", "base_damage", "cost", "text", "name", "each"}
    
    if was_conditional:
        new_effect = {"type": canonical}
        condition = {}
        # Loop over keys in the effect (except "type")
        for key, value in effect.items():
            if key == "type":
                continue
            if key in base_keys:
                new_effect[key] = value
            else:
                condition[key] = value
        if condition:
            new_effect["condition"] = condition
        return new_effect
    else:
        # Even non-conditional effects get their type updated
        effect["type"] = canonical
        return effect

def process_card_effects(card):
    """
    Normalize effects in a card's ability (if any) and in each attack.
    """
    if "ability" in card and card["ability"]:
        if "effects" in card["ability"] and isinstance(card["ability"]["effects"], list):
            card["ability"]["effects"] = [normalize_effect(eff) for eff in card["ability"]["effects"]]
    
    if "attacks" in card and isinstance(card["attacks"], list):
        for attack in card["attacks"]:
            if "effects" in attack and isinstance(attack["effects"], list):
                attack["effects"] = [normalize_effect(eff) for eff in attack["effects"]]
    return card

def main():
    if not os.path.exists(INPUT_FILE):
        print(f"Error: Input file {INPUT_FILE} not found.")
        return

    with open(INPUT_FILE, "r", encoding="utf-8") as f:
        cards = json.load(f)

    # Normalize all cards
    cleaned_cards = [process_card_effects(card) for card in cards]

    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        json.dump(cleaned_cards, f, indent=2)
    
    print(f"Cleaned data saved to {OUTPUT_FILE}")

if __name__ == "__main__":
    main()
