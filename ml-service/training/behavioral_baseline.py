import pandas as pd

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


# Calculate behavioral baseline
baseline = X.mean()

# Calculate normal variation
feature_std = X.std()


# Calculate deviation from baseline
deviation = (X - baseline).abs()

# Normalize each feature's deviation
normalized_deviation = deviation / feature_std

# Create one score per event
behavioral_scores = normalized_deviation.mean(axis=1)


# Fit calibration
calibration = fit_min_max(
    behavioral_scores
)


# Display results
print("Behavioral Baseline:")
print(baseline)

print()
print("Behavioral Anomaly Scores:")

for i, score in enumerate(behavioral_scores):
    print(
        f"Event {i + 1:2d} | "
        f"Raw Score: {score:.4f}"
    )


# Save behavioral model information
artifact = {
    "baseline": baseline,
    "feature_std": feature_std,
    "calibration": calibration,
}


save_artifact(
    artifact,
    "behavioral_baseline.joblib",
)