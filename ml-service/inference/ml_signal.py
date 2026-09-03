def interpret_hybrid_score(score):
    if score >= 0.75:
        return "high_anomaly"

    if score >= 0.50:
        return "moderate_anomaly"

    return "low_anomaly"


def create_ml_assessment(
    classifier_signal,
    isolation_signal,
    behavioral_signal,
    autoencoder_signal,
    hybrid_score,
):
    return {
        "signals": {
            "classifier": float(classifier_signal),
            "isolation_forest": float(isolation_signal),
            "behavioral": float(behavioral_signal),
            "autoencoder": float(autoencoder_signal),
        },
        "hybrid_score": float(hybrid_score),
        "interpretation": interpret_hybrid_score(hybrid_score),
    }