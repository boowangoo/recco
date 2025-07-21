import json
import numpy as np
import os
import pandas as pd
from sklearn.preprocessing import QuantileTransformer
from scipy.sparse import coo_matrix
from scipy.sparse.linalg import svds

def ratings_refit(ratings, qt=None):
    # Refitting ratings to the range [-1, +1]
    ratings_input = ratings.reshape(-1, 1)
    if qt is None:
        qt = QuantileTransformer(output_distribution="normal")
        refitted = qt.fit_transform(ratings_input).flatten()
    else:
        refitted = qt.transform(ratings_input).flatten()
    ratings_range = max(abs(refitted.min()), refitted.max())
    return refitted / ratings_range, qt

def movie_features(ratings_data, k):
    n_users = ratings_data.userId.max() + 1
    n_movies = ratings_data.movieId.max() + 1
    print(f"Calculating movie features for {n_users} users and {n_movies}")
    # Create sparse matrix, and compute SVD
    R = coo_matrix((ratings_data.rating, (ratings_data.userId, ratings_data.movieId)), shape=(n_users, n_movies)).tocsr()
    _, _, Vt = svds(R, k)
    return Vt.T

def preprocess_dataset(data_dir, ratings_vector_size):
    preprocessed_data = os.path.join(data_dir, "movies.parquet")
    if os.path.exists(preprocessed_data):
        print(f"Preprocessed dataset already exists at {preprocessed_data}")
        return preprocessed_data, None
    print("Preprocessing dataset")
    # Load and sort the data
    movies_data = pd.read_csv(os.path.join(data_dir, "movies.csv"), usecols=["movieId", "title", "genres"])
    ratings_data = pd.read_csv(os.path.join(data_dir, "ratings.csv"), usecols=["userId", "movieId", "rating"])
    movies_data = movies_data.sort_values("movieId", ascending=True)
    ratings_data = ratings_data.sort_values(["movieId", "userId"], ascending=[True, True])
    # Merge ratings_data with popular_movies
    movie_counts = ratings_data.movieId.value_counts()
    popular_movies = movie_counts[movie_counts >= 5].index
    movies_data = movies_data[movies_data.movieId.isin(popular_movies)]
    # Separate the year from the title
    movies_data = movies_data[movies_data.title.str.contains(r"\(\d{4}\)")]
    movies_data["year"] = movies_data.title.str.extract(r"\((\d{4})\)").astype(int)
    movies_data.title = movies_data.title.str.replace(r"\s*\(\d{4}\)", "", regex=True)
    movies_data.title = movies_data.title.str.strip()
    movies_data = movies_data[movies_data.title.str.len() > 0]
    # Split genres
    movies_data.genres = movies_data.genres.str.split("|")
    movies_data["genres"] = movies_data["genres"].apply(lambda x: [g.strip().lower() for g in x if g.strip()])
    movies_data = movies_data[movies_data.genres.str.len() > 0]
    # Calculate average ratings
    ratings_data = ratings_data[ratings_data.movieId.isin(movies_data.movieId)]
    avg_ratings = ratings_data.groupby("movieId").rating.mean()
    movies_data = movies_data.merge(avg_ratings.rename("average_rating"), left_on="movieId", right_index=True, how="left")
    # Readjust the ratings
    ratings_data.rating, qt = ratings_refit(ratings_data.rating.values.astype(float))
    ratings_original = [0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0]
    ratings_table = {
        "original": ratings_original,
        "refit": ratings_refit(np.array(ratings_original), qt)[0].tolist()
    }
    # Prepare the data to calculate movie features
    movies_data.movieId = pd.factorize(movies_data.movieId)[0]
    ratings_data.movieId = pd.factorize(ratings_data.movieId)[0]
    ratings_data.userId = pd.factorize(ratings_data.userId)[0]
    movies_features = movie_features(ratings_data, k=ratings_vector_size)
    movies_data["features"] = [movies_features[i] for i in range(len(movies_data))]
    # Save the processed data to parquet files
    movies_data.to_parquet(preprocessed_data)
    print(f"Finished preprocessing dataset to {preprocessed_data}")

    return preprocessed_data, ratings_table
