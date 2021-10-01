# Usage
- Ensure that Spectral is installed, Python3 is installed on the PATH.
- Ensure that `.spectral.yaml` and `scrape.py` are in this directory: https://github.com/APIs-guru/openapi-directory/tree/main/APIs
- Run the script `python3 scrape.py`.
- If early termination is required, send a `SIGINT` to the process (Ctrl-C on most terminal emulators), and stats will be collected in `errors.csv` and `warnings.csv` in the same directory.
