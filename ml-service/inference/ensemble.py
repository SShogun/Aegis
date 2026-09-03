def calculate_hybrid_score(
    classifier_signal,
    isolation_signal,
    behavioral_signal,
    autoencoder_signal,
):
    weights = {
        "classifier": 0.30,
        "isolation": 0.25,
        "behavioral": 0.25,
        "autoencoder": 0.20,
    }

    score = (
        classifier_signal * weights["classifier"]
        + isolation_signal * weights["isolation"]
        + behavioral_signal * weights["behavioral"]
        + autoencoder_signal * weights["autoencoder"]
    )

    return score