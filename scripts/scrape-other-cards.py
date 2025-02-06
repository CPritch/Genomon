import os
import requests
from bs4 import BeautifulSoup

SAVE_DIR = ".\\data"
ROOT_URL = "https://ptcgpocket.gg/cards/"
# Ensure save directory exists
os.makedirs(SAVE_DIR, exist_ok=True)

def find_cards():
    response = requests.get(ROOT_URL, stream=True)
    if response.status_code != 200:
        print("Failed to fetch the webpage")
        return
    
    soup = BeautifulSoup(response.text, "html.parser")
    cards = soup.select("#rootOfProContentpokepocket > section > div")
    print(f"{len(cards)} found")
    print(response.text)
find_cards()