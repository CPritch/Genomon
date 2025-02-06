import json

def normalize_id(card_id):
    """Normalize ID formats to match across datasets.
    Example: 'P-A-001' → 'PROMO-001'.
    """
    if card_id.startswith("P-"):
        return card_id.replace("P-A-", "PROMO-")
    return card_id

def parse_attack_info(info):
    """Convert attack 'info' strings (e.g., '{GC} Vine Whip 40') into structured attacks."""
    if not info:
        return []
    
    # Split cost symbols and attack details
    cost_part, *rest = info.split('}', 1)
    cost_symbols = cost_part.strip('{')
    attack_details = rest[0].strip() if rest else ""
    
    # Split attack name and damage (assume damage is the last number)
    name, damage = None, 0
    if attack_details:
        parts = attack_details.rsplit(' ', 1)
        if len(parts) == 2 and parts[1].isdigit():
            name, damage = parts[0], int(parts[1])
    
    # Map symbols to full type names (e.g., 'G' → 'Grass')
    symbol_to_type = {
        'G': 'Grass', 'C': 'Colorless', 'F': 'Fire', 
        'W': 'Water', 'P': 'Psychic', 'L': 'Lightning',
        'D': 'Darkness', 'M': 'Metal'
    }
    cost = [symbol_to_type.get(c, c) for c in cost_symbols]
    
    return [{
        "name": name or "Unknown Attack",
        "cost": cost,
        "damage": damage,
        "effects": []
    }]

def main():
    # Load datasets
    with open('./data/processed-pokemon.json', 'r', encoding='utf-8') as f:
        final_pokemon = json.load(f)
    with open('./data/cards.json', 'r', encoding='utf-8') as f:
        cards = json.load(f)
    
    # Create ID-mapped dictionaries
    final_pokemon_dict = {normalize_id(card['id']): card for card in final_pokemon}
    cards_dict = {normalize_id(card['id']): card for card in cards}
    
    # Merge existing cards
    missing_in_cards = []
    for card in final_pokemon:
        card_id = normalize_id(card['id'])
        if card_id in cards_dict:
            matched = cards_dict[card_id]
            card.update({
                "health": int(matched['hp']),
                "requires": matched.get('prew_stage_name'),
                "rarity": matched.get('rarity', card['rarity']),
                "dex": matched.get('dex'),
                "rule": matched.get('rule'),
                "weakness": matched.get('weakness')
            })
        else:
            missing_in_cards.append(card['id'])
    
    # Add missing cards from cards.json
    missing_in_final_pokemon = []
    for card_id, card_data in cards_dict.items():
        if card_id not in final_pokemon_dict:
            # Create a new entry in final-pokemon format
            new_card = {
                "id": card_data['id'].replace("PROMO-", "P-A-"),
                "name": card_data['name'],
                "rarity": card_data.get('rarity', '◇'),
                "type": card_data.get('color'),
                "stage": card_data.get('stage'),
                "requires": card_data.get('prew_stage_name'),
                "retreatCost": int(card_data.get('retreat', 0)),
                "attacks": parse_attack_info(card_data.get('attack', [{}])[0].get('info', '')),
                "health": int(card_data.get('hp', 0)),
                "dex": card_data.get('dex'),
                "rule": card_data.get('rule'),
                "weakness": card_data.get('weakness')
            }
            final_pokemon.append(new_card)
            missing_in_final_pokemon.append(card_data['id'])
    
    # Print warnings
    if missing_in_cards:
        print(f"Warning: Missing in cards.json: {missing_in_cards}")
    if missing_in_final_pokemon:
        print(f"Warning: Added new cards from cards.json: {missing_in_final_pokemon}")
    
    # Save merged data
    with open('./data/merged-cards.json', 'w', encoding='utf-8') as f:
        json.dump(final_pokemon, f, indent=2, ensure_ascii=False)
    print("Merged data saved to /data/merged-cards.json")


main()