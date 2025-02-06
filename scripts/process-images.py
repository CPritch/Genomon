import os
import json
import easyocr
from PIL import Image

JSON_PATH = "./data/refined-pokemon.json"
JSON_OUT_PATH = "./data/processed-pokemon.json"
INPUT_DIR = "./data/img"
OUTPUT_DIR = "./data/img-crop"
CROP_PERCENTAGE_WIDTH = 0.25
CROP_PERCENTAGE_HEIGHT = 0.1

# Ensure output directory exists
os.makedirs(OUTPUT_DIR, exist_ok=True)
reader = easyocr.Reader(['en'])

def crop_top_right(image_path):
    with Image.open(image_path) as img:
        width, height = img.size
        crop_width = int(width * CROP_PERCENTAGE_WIDTH)
        crop_height = int(height * CROP_PERCENTAGE_HEIGHT)
        
        return img.crop((width - crop_width, 10, width * 0.90, crop_height))

def process_image(card):
    input_path = os.path.join(INPUT_DIR, f"{card['id'].replace("-","_").replace("P_","P-")}.png")
    output_path = os.path.join(OUTPUT_DIR, f"{card['id']}.png")
    img = crop_top_right(input_path)
    width, height = img.size
    img = img.convert("RGBA").resize((width * 2, height * 2))
    img.save(output_path)
    matches = reader.readtext(output_path)
    digits = [match for match in matches if match[1].isdigit()]
    if len(digits) > 0:
        health = int(digits[0][1])
        return guess_real_number(health)
    else:
        return

def guess_real_number(proposed): #Assuming out of bounds is close digits wise (HP always ends in a 0 currently in PTCGP)
    if proposed < 10:
        return proposed * 10
    if proposed > 500:
        return int(str(proposed)[:-1])
    return proposed

def process_json():
    with open(JSON_PATH) as f:
        pokemon = json.load(f)
        for card in pokemon:
            if card['type'] == "Supporter" or card['type'] == "Item" or card['type'] == "Pokemon Tool": 
                continue
            health = process_image(card)
            card['health'] = health
            print(f"{card['name']} - HP{card['health']}")
        with open(JSON_OUT_PATH, "w") as w:
            w.write(json.dumps(pokemon))

process_json()