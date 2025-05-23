import os
import json
import requests
import pandas as pd
from tqdm import tqdm
from dotenv import dotenv_values

recco_env = dotenv_values("./recco.env")

if "RECCO_IP" not in recco_env or "RECCO_LB_PORT" not in recco_env:
    raise Exception("RECCO_IP and RECCO_LB_PORT must be set in the recco.env file.")

lb_addr = f"{recco_env['RECCO_IP']}:{recco_env['RECCO_LB_PORT']}"
headers = {'Content-Type': 'application/json'}
collection_name = "movie_titles"
collection_url = f"http://{lb_addr}/collections/{collection_name}"


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
    # TODO Adjust the vector size automatically according to the recco-embed model
    payload = {
        "name": collection_name,
        "vectors": {
            "size": 1024,  # Adjust this size according to your embedding model
            "distance": "Cosine"
        }
    }
    response = requests.put(collection_url, json=payload, headers=headers)
    if response.status_code != 200:
        raise Exception(f"Failed to create collection: {response.json()}")
    print(f"Collection '{collection_name}' created successfully.")

# Load movie titles from the dataset
# TODO retrieve the dataset from the cloud
movies_df = pd.read_csv('./dataset/movies_metadata.csv', usecols=['id', 'title'])
movies = movies_df[["id", "title"]].dropna()
ids = movies["id"].tolist()
titles = movies["title"].tolist()

print(f"Loading {len(titles)} titles from the dataset:")

batch_size = 32

for i in tqdm(range(0, len(titles), batch_size)):
    batch = titles[i:i+batch_size]
    data = json.dumps({"inputs": batch})
    response = requests.post(
        f'http://{lb_addr}/embed',
        data=data,
        headers=headers
    )
    if response.status_code != 200:
        raise Exception(f"Embedding request failed with status code {response.status_code}: {response.json()}")
    vectors = response.json()


    batch_ids = [int(id) for id in ids[i:i+batch_size]]
    payloads = [{"title": title} for title in batch]

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