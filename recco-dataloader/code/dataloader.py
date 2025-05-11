import json
import requests
import pandas as pd

headers = {'Content-Type': 'application/json'}
collection_name = "movie_titles"
collection_url = f"http://127.0.0.1:6333/collections/{collection_name}"

try:
    # Check if the collection already exists
    exists_response = requests.get(f"{collection_url}/exists", headers=headers)
    if exists_response.status_code == 200:
        print(f"Collection '{collection_name}' already exists.")
    else:
        print(f"Collection '{collection_name}' does not exist. Creating a new one.")
except Exception as e:
    print(f"Error checking collection {collection_name} existence: {e}")
    exists_response = None

# TODO Adjust the vector size automatically according to the recco-embed model
payload = {
    "vectors": {
        "size": 768,
        "distance": "Cosine"
    }
}

response = requests.put(collection_url, json=payload, headers=headers)
print(response.json())

movies_df = pd.read_csv('/dataset/movies_metadata.csv', usecols=['id', 'title'])
movies = movies_df[["id", "title"]].dropna()
ids = movies["id"].tolist()
titles = movies["title"].tolist()

batch_size = 32

for i in range(0, len(titles), batch_size):
    print(f"Processing batch {i//batch_size + 1}/{len(titles)//batch_size + 1}")
    batch = titles[i:i+batch_size]
    data = json.dumps({"inputs": batch})
    response = requests.post(
        'http://recco-embed/embed',
        data=data,
        headers=headers
    )
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
        f"http://recco-db/collections/{collection_name}/points",
        json=batch_data,
        headers=headers
    )
    print(response.json())