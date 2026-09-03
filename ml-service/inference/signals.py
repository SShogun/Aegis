def create_signal(
    signal_type,
    score,
    source,
    description,
):
    return {
        "type": signal_type,
        "score": float(score),
        "source": source,
        "description": description,
    }