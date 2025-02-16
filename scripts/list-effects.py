import json

with open("./data/cleaned-cards.json", "r", encoding="utf-8") as f:
    cards = json.load(f)

effect_types = set()

for card in cards:
    if "ability" in card and "effects" in card["ability"]:
        for effect in card["ability"]["effects"]:
            effect_types.add(effect["type"])

    if "attacks" in card:
        for attack in card["attacks"]:
            if "effects" in attack:
                for effect in attack["effects"]:
                    effect_types.add(effect["type"])

print("Unique effect types:", sorted(effect_types))
