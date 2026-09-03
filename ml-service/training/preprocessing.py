import pandas as pd
from sklearn.preprocessing import StandardScaler


FEATURES = [
    "hour",
    "failed_logins",
    "requests_per_minute",
    "files_downloaded",
]


def load_dataset(path):
    return pd.read_csv(path)


def get_features(data):
    return data[FEATURES]


def create_scaler(X):
    scaler = StandardScaler()
    scaler.fit(X)

    return scaler


def transform_features(X, scaler):
    return scaler.transform(X)