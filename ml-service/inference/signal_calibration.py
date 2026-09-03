import numpy as np


def fit_min_max(values):
    values = np.asarray(values, dtype=float)

    return {
        "minimum": float(values.min()),
        "maximum": float(values.max()),
    }


def apply_min_max(value, calibration):
    minimum = calibration["minimum"]
    maximum = calibration["maximum"]

    if maximum == minimum:
        return 0.0

    score = (
        (float(value) - minimum)
        / (maximum - minimum)
    )

    return float(np.clip(score, 0.0, 1.0))


def apply_inverted_min_max(value, calibration):
    return 1.0 - apply_min_max(
        value,
        calibration,
    )