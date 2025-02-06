import os
import requests
from bs4 import BeautifulSoup

URL = "https://game8.co/games/Pokemon-TCG-Pocket/archives/482685"
SAVE_DIR = ".\data\img"

# Ensure save directory exists
os.makedirs(SAVE_DIR, exist_ok=True)

def download_image(url, filename):
    response = requests.get(url, stream=True)
    if response.status_code == 200:
        filepath = os.path.join(SAVE_DIR, filename)
        with open(filepath, 'wb') as file:
            for chunk in response.iter_content(1024):
                file.write(chunk)
        print(f"Downloaded: {filename}")
    else:
        print(f"Failed to download: {url}")

def scrape_and_download():
    response = requests.get(URL)
    if response.status_code != 200:
        print("Failed to fetch the webpage")
        return
    
    soup = BeautifulSoup(response.text, "html.parser")
    rows = soup.select("div.archive-style-wrapper > div.scroll--table.table-header--fixed > table tr")
    
    for index, row in enumerate(rows, start=1):
        cells = row.find_all("td")
        if len(cells) < 3:
            continue
        
        img_tag = cells[2].find("img")
        if not img_tag or not img_tag.get("data-src"):
            continue
        
        image_url = img_tag["data-src"]
        filename = f"card-{str(index).zfill(3)}.png"
        
        # Optional: Extract name from another cell (e.g., cell[1])
        if len(cells) > 1 and cells[1].text.strip():
            name = cells[1].text.strip().replace(" ", "_").replace("/", "-")
            filename = f"{name}.png"
        
        download_image(image_url, filename)

scrape_and_download()
