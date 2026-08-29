import os
import re

for root, _, files in os.walk('.'):
    for file in files:
        if not file.endswith('.go'): continue
        path = os.path.join(root, file)
        with open(path, 'r') as f:
            content = f.read()
        
        new_content = content.replace('lib.CliValue', 'cli.CliValue')
        
        if content != new_content:
            with open(path, 'w') as f:
                f.write(new_content)

