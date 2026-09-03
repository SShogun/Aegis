import pandas as pd

from sklearn.neural_network import MLPRegressor
from sklearn.preprocessing import StandardScaler

from preprocessing import load_dataset, get_features
from artifacts import save_artifact

import sys
sys.path.append("../inference")

from signal_calibration import fit_min_max


DATA_PATH = "../data/events.csv"


# Load dataset
data = load_dataset(DATA_PATH)

# Select features
X = get_features(data)


# Scale features
scaler = StandardScaler()

X_scaled = scaler.fit_transform(X)


# Create autoencoder
model = MLPRegressor(
    hidden_layer_sizes=(2,),
    max_iter=2000,
    random_state=42,
)


# Train
model.fit(
    X_scaled,
    X_scaled,
)


# Reconstruct
reconstructed = model.predict(
    X_scaled
)


# Calculate reconstruction error
reconstruction_errors = (
    (X_scaled - reconstructed) ** 2
).mean(axis=1)


# Fit calibration
calibration = fit_min_max(
    reconstruction_errors
)


# Display results
print("Autoencoder Reconstruction Error:")
print()

for i, error in enumerate(
    reconstruction_errors
):
    print(
        f"Event {i + 1:2d} | "
        f"Error: {error:.4f}"
    )


# Save artifacts
save_artifact(
    model,
    "autoencoder.joblib",
)

save_artifact(
    scaler,
    "autoencoder_scaler.joblib",
)

save_artifact(
    calibration,
    "autoencoder_calibration.joblib",
)