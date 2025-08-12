import json
import requests
import os

# Load your .har file
with open('path_to_your_file.har', 'r') as f:
    har_data = json.load(f)

# Directory to save images
os.makedirs('downloaded_images', exist_ok=True)

# Extract and download .png files
for entry in har_data['log']['entries']:
    url = entry['request']['url']
    if url.endswith('.png'):
        try:
            print(f"Downloading: {url}")
            response = requests.get(url, stream=True)
            filename = url.split('/')[-1].split('?')[0]
            with open(os.path.join('downloaded_images', filename), 'wb') as out_file:
                out_file.write(response.content)
        except Exception as e:
            print(f"Failed to download {url}: {e}")