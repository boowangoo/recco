import argparse
from download_dataset import download_dataset
from preprocess_dataset import preprocess_dataset
from upload_movie_data import upload_movie_data, upload_rating_conversion

def parse_args():
    parser = argparse.ArgumentParser(description="Load movie data into recco-db.")
    parser.add_argument('--recco_ip', required=True, help='Recco server IP address')
    parser.add_argument('--recco_lb_port', required=True, help='Recco load balancer port')
    parser.add_argument('--movie_titles_embedding_dim', required=True, help='Embedding dimension for movie titles')
    parser.add_argument('--movie_ratings_embedding_dim', required=True, help='Embedding dimension for movie ratings')
    return parser.parse_args()

if __name__ == "__main__":
    args = parse_args()
    title_vector_size = int(args.movie_titles_embedding_dim)
    ratings_vector_size = int(args.movie_ratings_embedding_dim)
    host = f"http://{args.recco_ip}:{args.recco_lb_port}"

    dataset_dir = "/app/dataset"
    download_dataset(dataset_dir)
    movies_parquet, ratings_table = preprocess_dataset(dataset_dir, ratings_vector_size)
    upload_movie_data(host, movies_parquet, title_vector_size, ratings_vector_size, batch_size=16)
    upload_rating_conversion(host, ratings_table)
    print("Movie data loaded successfully.")
