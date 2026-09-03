import subprocess
import sys


scripts = [
    "train_classifier.py",
    "train_isolation_forest.py",
    "behavioral_baseline.py",
    "train_autoencoder.py",
]


for script in scripts:
    print()
    print("=" * 50)
    print(f"Running {script}")
    print("=" * 50)

    subprocess.run(
        [sys.executable, script],
        check=True,
    )


print()
print("=" * 50)
print("ALL MODELS TRAINED SUCCESSFULLY")
print("=" * 50)