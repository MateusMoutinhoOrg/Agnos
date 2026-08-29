import os
import glob
import re

# We need to replace `lib.` with appropriate packages based on the type.
# types mapping:
# SandBox -> sandbox.SandBox
# CliCommand -> cli.CliCommand
# CliApi -> cli.CliApi
# CliEntrys -> cli.CliEntrys
# CliArg -> cli.CliArg
# Cliflag -> cli.Cliflag
# CliTypeString, CliTypeInt, CliTypeBool, CliTypeFloat -> cli.CliType...
# ExitOk, ExitUsage, ExitFailure -> cli.ExitOk...
# Config -> config.Config
# CoreApi -> core.CoreApi
# StartProps, BuildProps, InstallProps, UninstallProps, ListProps, ExtensionHelpProps -> core.StartProps...

# Also update imports. Remove "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
# and add the needed ones.

def get_needed_imports(content):
    imports = set()
    if re.search(r'\bcli\.', content):
        imports.add('"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"')
    if re.search(r'\bconfig\.', content):
        imports.add('"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/config"')
    if re.search(r'\bcore\.', content):
        imports.add('"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"')
    if re.search(r'\bsandbox\.', content) and not 'package sandbox' in content:
        # Wait, if it's package sandbox, it doesn't need to import its own contract unless it's not the same.
        imports.add('"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/sandbox"')
    return imports

def replace_lib_usages(content):
    replacements = {
        r'\blib\.SandBox\b': 'sandbox.SandBox',
        r'\blib\.CliCommand\b': 'cli.CliCommand',
        r'\blib\.CliApi\b': 'cli.CliApi',
        r'\blib\.CliEntrys\b': 'cli.CliEntrys',
        r'\blib\.CliArg\b': 'cli.CliArg',
        r'\blib\.Cliflag\b': 'cli.Cliflag',
        r'\blib\.CliTypeString\b': 'cli.CliTypeString',
        r'\blib\.CliTypeInt\b': 'cli.CliTypeInt',
        r'\blib\.CliTypeFloat\b': 'cli.CliTypeFloat',
        r'\blib\.CliTypeBool\b': 'cli.CliTypeBool',
        r'\blib\.ExitOk\b': 'cli.ExitOk',
        r'\blib\.ExitUsage\b': 'cli.ExitUsage',
        r'\blib\.ExitFailure\b': 'cli.ExitFailure',
        r'\blib\.Config\b': 'config.Config',
        r'\blib\.CoreApi\b': 'core.CoreApi',
        r'\blib\.StartProps\b': 'core.StartProps',
        r'\blib\.BuildProps\b': 'core.BuildProps',
        r'\blib\.InstallProps\b': 'core.InstallProps',
        r'\blib\.UninstallProps\b': 'core.UninstallProps',
        r'\blib\.ListProps\b': 'core.ListProps',
        r'\blib\.ExtensionHelpProps\b': 'core.ExtensionHelpProps',
    }
    for old, new in replacements.items():
        content = re.sub(old, new, content)
    return content

for root, _, files in os.walk('.'):
    for file in files:
        if not file.endswith('.go'): continue
        path = os.path.join(root, file)
        with open(path, 'r') as f:
            content = f.read()

        if 'github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib' in content or re.search(r'\blib\.', content):
            new_content = replace_lib_usages(content)
            
            # Remove old lib import
            new_content = re.sub(r'(?m)^\s*"github\.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"\n', '', new_content)
            
            # Handlers update
            if 'CommandHandler(' in new_content:
                # signature change
                new_content = re.sub(r'func CommandHandler\(([\w]+) \*sandbox\.SandBox, ([\w]+) cli\.CliEntrys\) int \{', 
                                     r'func CommandHandler(\1 any, \2 cli.CliEntrys) int {\n\tsandbox := \1.(*sandbox.SandBox)', new_content)
            
            if 'package cli' in new_content and 'Handler func' in new_content:
                new_content = re.sub(r'Handler\s+func\(sandbox\s+\*sandbox\.SandBox,\s+entries\s+CliEntrys\)\s+int', 
                                     r'Handler func(sandbox any, entries CliEntrys) int', new_content)
                new_content = re.sub(r'(?m)^\s*"github\.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/sandbox"\n', '', new_content)
            
            # Add new imports
            needed = get_needed_imports(new_content)
            if needed:
                # Find import block
                import_match = re.search(r'import\s+\((.*?)\)', new_content, re.DOTALL)
                if import_match:
                    imports_str = import_match.group(1)
                    for imp in needed:
                        if imp not in imports_str:
                            imports_str += f'\n\t{imp}'
                    new_content = new_content[:import_match.start(1)] + imports_str + new_content[import_match.end(1):]
                else:
                    # Single imports or no imports
                    import_match = re.search(r'import\s+"[^"]+"', new_content)
                    if import_match:
                        # Convert to block
                        pass # too complex, let's just use goimports later
            
            if content != new_content:
                with open(path, 'w') as f:
                    f.write(new_content)

