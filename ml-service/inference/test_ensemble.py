from ensemble import calculate_hybrid_score
from ml_signal import create_ml_assessment


classifier_signal = 0.90
isolation_signal = 0.85
behavioral_signal = 0.75
autoencoder_signal = 0.80


hybrid_score = calculate_hybrid_score(
    classifier_signal,
    isolation_signal,
    behavioral_signal,
    autoencoder_signal,
)


assessment = create_ml_assessment(
    classifier_signal,
    isolation_signal,
    behavioral_signal,
    autoencoder_signal,
    hybrid_score,
)


print("ML Assessment:")
print()

print(assessment)