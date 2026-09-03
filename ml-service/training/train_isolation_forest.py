import pandas as pd

from sklearn.ensemble import IsolationForest

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


# Create model
model = IsolationForest(
    n_estimators=100,
    contamination=0.25,
    random_state=42,
)


# Train model
model.fit(X)


# Generate raw scores
predictions = model.predict(X)
scores = model.decision_function(X)


# Fit calibration using training scores
calibration = fit_min_max(scores)


# Display results
print("Isolation Forest Results:")
print()

for i, (prediction, score) in enumerate(
    zip(predictions, scores)
):
    result = "Normal" if prediction == 1 else "Anomaly"

    print(
        f"Event {i + 1:2d} | "
        f"{result:7s} | "
        f"Raw Score: {score:.4f}"
    )


# Save artifacts
save_artifact(
    model,
    "isolation_forest.joblib",
)

save_artifact(
    calibration,
    "isolation_calibration.joblib",
)