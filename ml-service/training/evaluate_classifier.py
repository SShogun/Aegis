from sklearn.metrics import (
    accuracy_score,
    precision_score,
    recall_score,
    f1_score,
    confusion_matrix,
    classification_report,
)

from train_classifier import (
    X_test,
    y_test,
    predictions,
    model,
)


accuracy = accuracy_score(y_test, predictions)
precision = precision_score(y_test, predictions, zero_division=0)
recall = recall_score(y_test, predictions, zero_division=0)
f1 = f1_score(y_test, predictions, zero_division=0)

matrix = confusion_matrix(y_test, predictions)


print("Accuracy:", accuracy)
print("Precision:", precision)
print("Recall:", recall)
print("F1 Score:", f1)

print()
print("Confusion Matrix:")
print(matrix)

print()
print("Classification Report:")
print(classification_report(
    y_test,
    predictions,
    zero_division=0,
))

probabilities = model.predict_proba(X_test)

suspicious_probabilities = probabilities[:, 1]

for actual, prediction, probability in zip(
    y_test,
    predictions,
    suspicious_probabilities,
):
    print(
        f"Actual: {actual} | "
        f"Predicted: {prediction} | "
        f"Suspicious Probability: {probability:.4f}"
    )