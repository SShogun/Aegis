from preprocessing import (
    load_dataset,
    get_features,
    create_scaler,
    transform_features,
)


DATA_PATH = "../data/events.csv"


data = load_dataset(DATA_PATH)

X = get_features(data)

scaler = create_scaler(X)

X_scaled = transform_features(X, scaler)


print("Original shape:", X.shape)
print("Scaled shape:", X_scaled.shape)
print()
print("Scaled first row:")
print(X_scaled[0])