import os
import re

for root, _, files in os.walk('sandbox/cli/commands'):
    for file in files:
        if not file.endswith('.go'): continue
        path = os.path.join(root, file)
        with open(path, 'r') as f:
            content = f.read()
        
        # fix the argument name in CommandHandler
        new_content = re.sub(r'func CommandHandler\(sandbox any, entries cli\.CliEntrys\) int \{\n\tsandbox := sandbox\.\(\*sandbox\.SandBox\)', 
                             r'func CommandHandler(sb any, entries cli.CliEntrys) int {\n\tsandbox := sb.(*sandbox.SandBox)', content)
        
        # also for CliMainFactory.go
        if content != new_content:
            with open(path, 'w') as f:
                f.write(new_content)

# CliMainFactory.go fix
with open('sandbox/cli/CliMainFactory.go', 'r') as f:
    content = f.read()
new_content = content.replace("command.Handler(sandbox, entries)", "command.Handler(sandbox, entries)") # This is already correct
with open('sandbox/cli/CliMainFactory.go', 'w') as f:
    f.write(new_content)

