import sys

import pandas as pd

sys.path.append("../training")

from artifacts import load_artifact
from preprocessing import FEATURES

from signal_calibration import (
    apply_min_max,
    apply_inverted_min_max,
)

from ensemble import calculate_hybrid_score
from ml_signal import create_ml_assessment


# --------------------------------
# Load artifacts
# --------------------------------

classifier = load_artifact(
    "classifier.joblib"
)

classifier_scaler = load_artifact(
    "classifier_scaler.joblib"
)


isolation_forest = load_artifact(
    "isolation_forest.joblib"
)

isolation_calibration = load_artifact(
    "isolation_calibration.joblib"
)


behavioral_artifact = load_artifact(
    "behavioral_baseline.joblib"
)


autoencoder = load_artifact(
    "autoencoder.joblib"
)

autoencoder_scaler = load_artifact(
    "autoencoder_scaler.joblib"
)

autoencoder_calibration = load_artifact(
    "autoencoder_calibration.joblib"
)


# --------------------------------
# Prediction
# --------------------------------

def predict_event(event):

    X = pd.DataFrame(
        [[event[feature] for feature in FEATURES]],
        columns=FEATURES,
    )


    # --------------------------------
    # Classifier
    # --------------------------------

    X_classifier = (
        classifier_scaler.transform(X)
    )

    classifier_signal = (
        classifier.predict_proba(
            X_classifier
        )[0][1]
    )


    # --------------------------------
    # Isolation Forest
    # --------------------------------

    isolation_raw = (
        isolation_forest
        .decision_function(X)[0]
    )

    isolation_signal = (
        apply_inverted_min_max(
            isolation_raw,
            isolation_calibration,
        )
    )


    # --------------------------------
    # Behavioral Baseline
    # --------------------------------

    baseline = (
        behavioral_artifact["baseline"]
    )

    feature_std = (
        behavioral_artifact["feature_std"]
    )

    behavioral_score = 0.0

    for feature in FEATURES:

        value = X.loc[0, feature]

        deviation = abs(
            value - baseline[feature]
        )

        normalized_deviation = (
            deviation
            / feature_std[feature]
        )

        behavioral_score += (
            normalized_deviation
        )

    behavioral_score /= len(FEATURES)

    behavioral_signal = (
        apply_min_max(
            behavioral_score,
            behavioral_artifact[
                "calibration"
            ],
        )
    )


    # --------------------------------
    # Autoencoder
    # --------------------------------

    X_autoencoder = (
        autoencoder_scaler.transform(X)
    )

    reconstructed = (
        autoencoder.predict(
            X_autoencoder
        )
    )

    reconstruction_error = (
        (X_autoencoder - reconstructed) ** 2
    ).mean()

    autoencoder_signal = (
        apply_min_max(
            reconstruction_error,
            autoencoder_calibration,
        )
    )


    # --------------------------------
    # Hybrid score
    # --------------------------------

    hybrid_score = calculate_hybrid_score(
        classifier_signal,
        isolation_signal,
        behavioral_signal,
        autoencoder_signal,
    )


    # --------------------------------
    # Final ML assessment
    # --------------------------------

    return create_ml_assessment(
        classifier_signal,
        isolation_signal,
        behavioral_signal,
        autoencoder_signal,
        hybrid_score,
    )