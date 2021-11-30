import os
import subprocess
import sys
import yaml

from collections import defaultdict

class StyleGuideGenerator(object):
    def __init__(self):
        pass
    def generate(self):

        # Create a mapping between the AIP Guideline ID to the Guideline object
        guideline_id_to_guideline = defaultdict(lambda: {
            "id": "",
            "rules": [],
            "status": "ACTIVE",
        })

        rules_path = "./api-linter/docs/rules"

        # Iterate over all the rule markdown files, which are contained in
        # folders. Each folder name is the name of the guideline the rules belong to.
        for dirpath, _, filenames in os.walk(rules_path):
            if dirpath == rules_path:
                continue
            
            # Populate the ID of the Guideline
            guideline_id = int(os.path.basename(dirpath))
            guideline_id_to_guideline[guideline_id]["id"] = f"aip{guideline_id}"

            # Iterate over all the rule docs in the guideline
            for rule_doc_name in filenames:

                # Skip the index.md file, because it does not contain any rule info
                if rule_doc_name == "index.md":
                    continue

                # Get the path of the rule doc
                rule_doc_path = os.path.join(dirpath, rule_doc_name)

                # Open the rule doc for reading, and read all lines
                rule_doc = open(rule_doc_path,'r')
                lines = rule_doc.readlines()
                
                try:
                    # Get the start and end lines of the yaml contained within the MD file
                    yaml_start = lines.index("---\n")
                    yaml_end = len(lines)-lines[::-1].index("---\n")-1

                    # Get the YAML text within the MD file and parse it
                    yaml_text = "".join(lines[yaml_start+1 : yaml_end])
                    parsed_yaml = yaml.safe_load(yaml_text)

                    # Get metadata from the parsed YAML, and use it to populate the rule metadata
                    rule_name = os.path.basename(parsed_yaml["permalink"])

                    # Append the rule to the appropriate guideline
                    guideline_id_to_guideline[guideline_id]["rules"].append(
                        {
                            "id": rule_name.strip(),
                            "description" : parsed_yaml["rule"]["summary"].replace('\n', ' ').strip(),
                            "linter": "api-linter",
                            "linter_rulename": rule_name.strip(),
                            "severity": "ERROR",
                            "doc_uri" : "linter.aip.dev"+parsed_yaml["permalink"].strip()
                        }
                    )
                except Exception as e:
                    print(rule_doc_path)
                    pass

        return {
            "id": "api-linter-styleguide",
            "mime_types": [
                "application/x.protobuf+zip",
            ],
            "guidelines": [
                guideline_id_to_guideline[guideline_id] 
                for guideline_id in guideline_id_to_guideline
            ],
            "linters": [
                {
                    "name": "api-linter",
                    "uri": "https://github.com/googleapis/api-linter",
                }
            ],
        }
        

class ApiLinterStyleGuideGenerator(object):
    def __init__(self):
        subprocess.run(["git", "clone", "https://github.com/googleapis/api-linter.git"])
    
    def __enter__(self):
        return StyleGuideGenerator()
  
    def __exit__(self, exc_type, exc_value, tb):
        subprocess.run(["rm", "-rf", "api-linter"])

def show_usage():
    print("Usage: generate <output_file>")

if __name__ == "__main__":
    # Validate CLI args
    args = sys.argv
    if len(args) < 2:
        show_usage()
        sys.exit(1)
    
    output_file_path = sys.argv[1]

    with ApiLinterStyleGuideGenerator() as generator:
        # Generate the style guide
        style_guide = generator.generate()

        # Dump the generated style guide into an output file
        try:
            with open(output_file_path, 'w') as file:
                yaml.safe_dump(style_guide, file, sort_keys=False)
        except Exception as e:
            print(e)
            show_usage()
            sys.exit(1)
