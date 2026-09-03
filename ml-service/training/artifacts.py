from pathlib import Path
import joblib


MODELS_DIR = Path("../models")


def save_artifact(obj, filename):
    MODELS_DIR.mkdir(parents=True, exist_ok=True)

    path = MODELS_DIR / filename

    joblib.dump(obj, path)

    print(f"Saved: {path}")


def load_artifact(filename):
    path = MODELS_DIR / filename

    return joblib.load(path)