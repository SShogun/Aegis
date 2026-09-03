import sys

import pandas as pd

sys.path.append("../inference")

from predict import predict_event


DATA_PATH = "../data/events.csv"


# --------------------------------
# Load dataset
# --------------------------------

data = pd.read_csv(DATA_PATH)


# --------------------------------
# Evaluate every event
# --------------------------------

results = []


for index, row in data.iterrows():

    event = {
        "hour": row["hour"],
        "failed_logins": row["failed_logins"],
        "requests_per_minute": row["requests_per_minute"],
        "files_downloaded": row["files_downloaded"],
    }


    assessment = predict_event(event)


    results.append({
        "event": index + 1,
        "actual": int(row["is_suspicious"]),

        "classifier": assessment[
            "signals"
        ]["classifier"],

        "isolation_forest": assessment[
            "signals"
        ]["isolation_forest"],

        "behavioral": assessment[
            "signals"
        ]["behavioral"],

        "autoencoder": assessment[
            "signals"
        ]["autoencoder"],

        "hybrid": assessment[
            "hybrid_score"
        ],

        "interpretation": assessment[
            "interpretation"
        ],
    })


# --------------------------------
# Create results table
# --------------------------------

results_df = pd.DataFrame(results)


# --------------------------------
# Display results
# --------------------------------

print()
print("=" * 100)
print("AEGIS ML EVALUATION")
print("=" * 100)
print()

print(
    results_df.to_string(
        index=False,
        formatters={
            "classifier": "{:.3f}".format,
            "isolation_forest": "{:.3f}".format,
            "behavioral": "{:.3f}".format,
            "autoencoder": "{:.3f}".format,
            "hybrid": "{:.3f}".format,
        },
    )
)


# --------------------------------
# Summary
# --------------------------------

print()
print("=" * 100)
print("SUMMARY")
print("=" * 100)

print()
print(
    "Average Classifier Signal:",
    f"{results_df['classifier'].mean():.3f}",
)

print(
    "Average Isolation Forest Signal:",
    f"{results_df['isolation_forest'].mean():.3f}",
)

print(
    "Average Behavioral Signal:",
    f"{results_df['behavioral'].mean():.3f}",
)

print(
    "Average Autoencoder Signal:",
    f"{results_df['autoencoder'].mean():.3f}",
)

print(
    "Average Hybrid Score:",
    f"{results_df['hybrid'].mean():.3f}",
)


# --------------------------------
# Compare normal vs suspicious
# --------------------------------

print()
print("=" * 100)
print("NORMAL VS SUSPICIOUS")
print("=" * 100)

grouped = results_df.groupby("actual")[
    [
        "classifier",
        "isolation_forest",
        "behavioral",
        "autoencoder",
        "hybrid",
    ]
].mean()

grouped.index = [
    "Normal",
    "Suspicious",
]

print()

print(
    grouped.to_string(
        float_format=lambda value: f"{value:.3f}"
    )
)