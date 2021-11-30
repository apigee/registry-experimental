## Summary
This Python3 script generates an API Linter Style Guide based on the [API Linter Rules Documentation](https://github.com/googleapis/api-linter/tree/main/docs/rules)

Requirements: `pyyaml` should be installed in the environment.

## Test
- Create a virtual environment for Python3.
- `pip3 install pyyaml`

Run the script:
`python3 generate.py styleguide.yaml`

Verify that it generates [the following API Style Guide] in `styleguide.yaml`. (https://gist.github.com/muhammadharis/9c6a1662e06c4687a231e8852ef13999).
