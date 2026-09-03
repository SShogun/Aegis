from sklearn.model_selection import train_test_split
from sklearn.tree import DecisionTreeClassifier
from artifacts import save_artifact

from preprocessing import (
    load_dataset,
    get_features,
    create_scaler,
    transform_features,
)


DATA_PATH = "../data/events.csv"


data = load_dataset(DATA_PATH)

X = get_features(data)
y = data["is_suspicious"]


X_train, X_test, y_train, y_test = train_test_split(
    X,
    y,
    test_size=0.25,
    random_state=42,
    stratify=y,
)


scaler = create_scaler(X_train)

X_train_scaled = transform_features(X_train, scaler)
X_test_scaled = transform_features(X_test, scaler)


model = DecisionTreeClassifier(
    random_state=42
)

model.fit(X_train_scaled, y_train)


save_artifact(
    model,
    "classifier.joblib",
)

save_artifact(
    scaler,
    "classifier_scaler.joblib",
)
predictions = model.predict(X_test_scaled)

print("Predictions:")
print(predictions)

save_artifact(
    model,
    "classifier.joblib",
)

save_artifact(
    scaler,
    "classifier_scaler.joblib",
)