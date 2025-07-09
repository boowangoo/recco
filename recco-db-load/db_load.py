import os
import json
import requests
import pandas as pd
from tqdm import tqdm
import re
import zipfile
import shutil
import urllib.request

RECCO_IP = os.getenv("RECCO_IP")
RECCO_LB_PORT = os.getenv("RECCO_LB_PORT")
if not RECCO_IP or not RECCO_LB_PORT:
    raise Exception("RECCO_IP and RECCO_LB_PORT must be set in the environment variables.")

lb_addr = f"{RECCO_IP}:{RECCO_LB_PORT}"
headers = {'Content-Type': 'application/json'}
collection_name = "movie_titles"
collection_url = f"http://{lb_addr}/collections/{collection_name}"

print(f"lb_addr: {lb_addr}")

# Download and extract dataset if not present
dataset_dir = '/app/dataset'
movies_csv = os.path.join(dataset_dir, 'movies.csv')
if not os.path.exists(movies_csv):
    print("movies.csv not found, downloading and extracting dataset...")
    zip_path = os.path.join(dataset_dir, 'ml-32m.zip')
    os.makedirs(dataset_dir, exist_ok=True)
    url = 'https://files.grouplens.org/datasets/movielens/ml-32m.zip'
    urllib.request.urlretrieve(url, zip_path)
    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
        zip_ref.extractall(dataset_dir)
    # Find the extracted folder (should be ml-32m/)
    extracted_dir = os.path.join(dataset_dir, 'ml-32m')
    if os.path.exists(extracted_dir):
        for filename in os.listdir(extracted_dir):
            shutil.move(os.path.join(extracted_dir, filename), dataset_dir)
        shutil.rmtree(extracted_dir)
    os.remove(zip_path)
    print("Dataset downloaded and extracted.")


# Check if the collection already exists
exists_response = requests.get(f"{collection_url}/exists", headers=headers)
if exists_response.status_code != 200:
    raise Exception(f"Request failed with status code {exists_response.status_code}")
exists = exists_response.json().get("result", {}).get("exists", False)

if exists:
    print(f"Collection '{collection_name}' already exists.")
    # Check if the collection is empty
    cnt_response = requests.post(
        f"{collection_url}/points/count", 
        headers=headers,
        json={}
    )
    if cnt_response.status_code != 200:
        raise Exception(f"Failed to count points in collection: {cnt_response.json()}")
    count = cnt_response.json().get("count", 0)
    if count > 0:
        raise Exception(f"Collection '{collection_name}' is not empty.")
    else:
        print(f"Collection '{collection_name}' is empty. Loading data:")
else:
    print(f"Collection '{collection_name}' does not exist. Creating a new one.")
    
    vector_size = os.getenv("EMBEDDING_DIM")
    if not vector_size:
        raise Exception("EMBEDDING_DIM environment variable must be set.")

    payload = {
        "name": collection_name,
        "vectors": {
            "size": int(vector_size),
            "distance": "Cosine"
        }
    }
    response = requests.put(collection_url, json=payload, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Failed to create collection: {response.json()}")
    print(f"Collection '{collection_name}' created successfully.")

# Load movie titles from the dataset
# TODO retrieve the dataset from the cloud
movies_df = pd.read_csv('/app/dataset/movies.csv', usecols=['movieId', 'title', 'genres'])
movies = movies_df[["movieId", "title"]].dropna()
ids = movies["movieId"].tolist()
titles = movies["title"].tolist()
title_year_pattern = re.compile(r"(.+)\s+\((\d{4})\)")
parsed = [title_year_pattern.match(ty) for ty in titles]
titles = [m.group(1).strip() if m else title for m, title in zip(parsed, titles)]
years = [int(m.group(2)) if m else None for m in parsed]

print(f"Loading {len(titles)} titles from the dataset:")

batch_size = 32

for i in tqdm(range(0, len(titles), batch_size)):
    batch_titles = titles[i:i+batch_size]
    batch_years = years[i:i+batch_size]
    data = json.dumps({"inputs": batch_titles})
    response = requests.post(
        f'http://{lb_addr}/embed',
        data=data,
        headers=headers
    )
    if response.status_code != 200:
        raise Exception(f"Embedding request failed with status code {response.status_code}: {response.json() if response else ''}")
    vectors = response.json()


    batch_ids = [int(id) for id in ids[i:i+batch_size]]
    payloads = [{"title": title, "year": year} for title, year in zip(batch_titles, batch_years)]

    batch_data = {
        "batch": {
            "ids": batch_ids,
            "vectors": vectors,
            "payloads": payloads,
        }
    }
    
    response = requests.put(
        f"http://{lb_addr}/collections/{collection_name}/points",
        json=batch_data,
        headers=headers
    )
    if response.status_code != 200:
        raise Exception(f"Failed to insert batch into collection: {response.json()}")

print(f"Successfully loaded {len(titles)} titles into '{collection_name}'.")
