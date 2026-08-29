import os

for root, _, files in os.walk('.'):
    for file in files:
        if not file.endswith('.go'): continue
        path = os.path.join(root, file)
        with open(path, 'r') as f:
            content = f.read()
        
        # replace the old import with the new one
        new_content = content.replace('"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/sandbox"', '"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"')
        
        if content != new_content:
            with open(path, 'w') as f:
                f.write(new_content)
