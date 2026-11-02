import subprocess
import sys

r = subprocess.run(
    [sys.executable, "-m", "unittest", "-v", "test_duration"],
    capture_output=True,
    text=True,
)
print(r.stdout)
print(r.stderr)
print("exit code:", r.returncode)
