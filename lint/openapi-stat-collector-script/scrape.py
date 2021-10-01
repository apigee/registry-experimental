import os
import signal
import sys
from collections import defaultdict

# Maps map the rule_name : organization_name : violation_count
warning_design_rule_to_occurrence = defaultdict(lambda: defaultdict(lambda: 0))
error_design_rule_to_occurrence = defaultdict(lambda: defaultdict(lambda: 0))
organization_name_to_specs = defaultdict(lambda: 1)

def design_rule_to_occurrence_index(organization_name, output):
    for line in output.split('\n'):
        tokens = line.split()
        if len(tokens) < 3:
            continue

        type = tokens[1]
        rule_name = tokens[2]

    
        # Record metrics
        if type == "warning":
            warning_design_rule_to_occurrence[rule_name][organization_name] += 1
        elif type == "error":
            error_design_rule_to_occurrence[rule_name][organization_name] += 1


def dfs(cur_dir_tokens):
    
    # Iterate through all files in the current directory.
    for filename in os.listdir(os.path.join("", *cur_dir_tokens)):

        # Skip files that are redundant.
        if filename in {".", ".."}:
            continue

        # Construct the absolute path to the file that's being iterated on.
        next_dir_tokens = cur_dir_tokens + [filename]
        next_dir_str = os.path.join("", *next_dir_tokens)

        # If this file is a directory, then it should be recursed on.
        if os.path.isdir(next_dir_str):
            dfs(next_dir_tokens)

        # If this is a yaml file, it is likely a spec. We should run the linter on it.
        elif filename.endswith(".yaml"): 
            print(next_dir_str)

            # This file is extremely large and causes the program to hang. Ignore it.
            if next_dir_str == "./beezup.com/2.0/openapi.yaml":
                continue
            stream = os.popen(f"spectral lint {next_dir_str}")
            output = stream.read()
            if len(cur_dir_tokens) >= 2:
                organization_name = cur_dir_tokens[1]
                organization_name_to_specs[organization_name] += 1

                # Get the actual metrics
                design_rule_to_occurrence_index(organization_name, output)


def flush_to_csv(file_name, data_map):
    import csv

    with open(file_name, 'w') as f:

        # create the csv writer
        writer = csv.writer(f)
        writer.writerow(["name", "organization_name", "occurrences"])

        for tag in data_map.keys():
            for organization_name in data_map[tag]:
                violations = data_map[tag][organization_name]/organization_name_to_specs[organization_name]
                writer.writerow([tag, organization_name, violations])
    

def flush_data():
    flush_to_csv("./warnings.csv", warning_design_rule_to_occurrence)
    flush_to_csv("./errors.csv", error_design_rule_to_occurrence)


def signal_handler(signal, frame):
     # join all processes
    flush_data()
    sys.exit(0)


if __name__ == "__main__":
    # Register signal handler
    signal.signal(signal.SIGINT, signal_handler)

    dfs(['.'])
    flush_data()

