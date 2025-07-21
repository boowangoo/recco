import os
import shutil
import urllib.request
import zipfile

def check_dataset_exists(dataset_dir, datasets):
    dataset_exists = True
    if os.path.exists(dataset_dir):
        for dataset in datasets:
            if not os.path.exists(os.path.join(dataset_dir, dataset)):
                dataset_exists = False
                break
    else:
        dataset_exists = False
    return dataset_exists

# Download and extract dataset if not present
def download_dataset(dataset_dir):
    datasets = ["movies.csv", "ratings.csv", "links.csv", "tags.csv"]
    if check_dataset_exists(dataset_dir, datasets):
        print("Dataset already exists, skipping download.")
        return
    print("Downloading and extracting dataset...")
    zip_path = os.path.join(dataset_dir, 'ml-32m.zip')

    # Download the dataset zip file
    # TODO host and retrieve the dataset from the cloud
    dataset_url = 'https://files.grouplens.org/datasets/movielens/ml-32m.zip'
    try:
        urllib.request.urlretrieve(dataset_url, zip_path)
    except Exception:
        raise Exception("Dataset download failed.")
    # Extract the zip file
    with zipfile.ZipFile(zip_path, 'r') as zip_ref:
        zip_ref.extractall(dataset_dir)
    extracted_dir = os.path.join(dataset_dir, 'ml-32m')
    if not os.path.exists(extracted_dir):
        raise Exception("Dataset extraction failed, extracted folder not found.")
    os.remove(zip_path)
    shutil.copytree(src=extracted_dir, dst=dataset_dir, dirs_exist_ok=True)
    shutil.rmtree(extracted_dir)
    
    # Confirm that all datasets are present
    if not check_dataset_exists(dataset_dir, datasets):
        raise Exception("Some datasets are missing.")
    print("All datasets downloaded and extracted successfully.")