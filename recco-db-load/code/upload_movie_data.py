import json
import pandas as pd
import re
import requests
from tqdm import trange

JSON_HEADER = {'Content-Type': 'application/json'}

def check_collection_exists(collection_url, collection_name):
    # Check if the collection already exists
    exists_response = requests.get(f"{collection_url}/{collection_name}/exists",
                                   headers=JSON_HEADER)
    if exists_response.status_code != 200:
        raise Exception(f"Request failed with status code {exists_response.status_code}")
    exists = exists_response.json().get("result", {}).get("exists", False)
    return exists

def create_collection(collection_url, collection_name, vectors_config):
    empty_collection = True
    if check_collection_exists(collection_url, collection_name):
        print(f"Collection '{collection_name}' already exists.")
        # Check if the collection is empty
        cnt_response = requests.post(
            f"{collection_url}/{collection_name}/points/count",
            json={},
            headers=JSON_HEADER
        )
        if cnt_response.status_code != 200:
            raise Exception(f"Failed to count points in collection: {cnt_response.json()}")

        results = cnt_response.json().get("result", None)
        if results is None:
            raise Exception("Count response is empty or malformed.")
        count = results.get("count", 0)
        empty_collection = (count == 0)
    else:
        print(f"Collection '{collection_name}' does not exist. Creating a new one.")
        payload = {
            "name": collection_name,
            "vectors": vectors_config
        }
        response = requests.put(f"{collection_url}/{collection_name}",
                                json=payload, headers=JSON_HEADER)
        if response.status_code != 200:
            raise Exception(f"Failed to create collection: {response.json()}")
        print(f"Collection '{collection_name}' created successfully.")
    return empty_collection

def upload_movie_data(host, movies_parquet, title_vector_size, ratings_vector_size, batch_size=32):
    collection_url = f"{host}/collections"
    collection_name = "movies"

    vectors_config = {
        "title": {
            "size": title_vector_size,
            "distance": "Cosine"
        },
        "features": {
            "size": ratings_vector_size,
            "distance": "Cosine"
        }
    }

    if not create_collection(collection_url, collection_name, vectors_config):
        print(f"Collection '{collection_name}' is already populated.")
        return

    movies = pd.read_parquet(movies_parquet, columns=["movieId", "title", "genres", "year", "average_rating", "features"])
    n_movies = len(movies)

    for i in trange(0, n_movies, batch_size, desc="Uploading movies"):
        j = min(i+batch_size, n_movies)
        n_batch = j - i
        batch_rows = movies[i:j]
        batch_title = batch_rows["title"].tolist()
        batch_genres = [g.tolist() for g in batch_rows["genres"].tolist()]
        batch_year = batch_rows["year"].tolist()
        batch_movieIds = batch_rows["movieId"].tolist()
        batch_features = [f.tolist() for f in batch_rows["features"].tolist()]
        batch_average_rating = batch_rows["average_rating"].tolist()

        response = requests.post(
            f"{host}/embed",
            json={"inputs": batch_title},
            headers=JSON_HEADER
        )
        if response.status_code == 413:
            print(f"payload too large! inputs: {batch_title}")
            raise Exception("Payload too large for embedding service.")
        if response.status_code != 200:
            raise Exception(f"Embedding request failed with status code {response.status_code}: {response.json() if response else ''}")
        batch_title_vectors = response.json()

        batch_data = {
            "batch": {
                "ids": batch_movieIds,
                "vectors": {
                    "title": batch_title_vectors,
                    "features": batch_features
                },
                "payloads": [{
                    "title": batch_title[i],
                    "genres": batch_genres[i],
                    "year": batch_year[i],
                    "average_rating": batch_average_rating[i]
                } for i in list(range(n_batch))]
            }
        }
        
        response = requests.put(
            f"{collection_url}/{collection_name}/points",
            json=batch_data,
            headers=JSON_HEADER
        )
        if response.status_code != 200:
            raise Exception(f"Failed to insert batch into collection: {response.json()}")
    
    print(f"Successfully loaded {n_movies} titles into '{collection_name}'.")

def upload_ratings_table(host, ratings_table):
    collection_url = f"{host}/collections"
    collection_name = "ratings"

    vector_conf = { "size": len(ratings_table["original"]), "distance": "Cosine" }
    vectors_config = {
        "original": vector_conf,
        "refit": vector_conf
    }

    if not create_collection(collection_url, collection_name, vectors_config):
        print(f"Collection '{collection_name}' is already populated.")
        return

    response = requests.put(
        f"{collection_url}/{collection_name}/points",
        json={
            "points": [{
                "id": 0,
                "vector": ratings_table
            }]
        },
        headers=JSON_HEADER
    )
    if response.status_code != 200:
        raise Exception(f"Failed to insert batch into collection: {response.json()}")
    print(f"Successfully uploaded rating conversion data to '{collection_name}'.")
