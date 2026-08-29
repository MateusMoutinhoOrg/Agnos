import os
import shutil
import glob

# Create dirs
os.makedirs("sandbox/contracts/lib/cli", exist_ok=True)
os.makedirs("sandbox/contracts/lib/config", exist_ok=True)
os.makedirs("sandbox/contracts/lib/core", exist_ok=True)
os.makedirs("sandbox/contracts/sandbox", exist_ok=True)

# Move files
if os.path.exists("sandbox/contracts/lib/cli.go"):
    shutil.move("sandbox/contracts/lib/cli.go", "sandbox/contracts/lib/cli/cliApi.go")
if os.path.exists("sandbox/contracts/lib/config.go"):
    shutil.move("sandbox/contracts/lib/config.go", "sandbox/contracts/lib/config/configApi.go")
if os.path.exists("sandbox/contracts/lib/core.go"):
    shutil.move("sandbox/contracts/lib/core.go", "sandbox/contracts/lib/core/coreApi.go")
if os.path.exists("sandbox/contracts/lib/sandbox.go"):
    shutil.move("sandbox/contracts/lib/sandbox.go", "sandbox/contracts/sandbox/sandbox.go")

# Replace package names
def replace_in_file(path, old, new):
    with open(path, 'r') as f:
        content = f.read()
    content = content.replace(old, new)
    with open(path, 'w') as f:
        f.write(content)

replace_in_file("sandbox/contracts/lib/cli/cliApi.go", "package lib", "package cli")
replace_in_file("sandbox/contracts/lib/config/configApi.go", "package lib", "package config")
replace_in_file("sandbox/contracts/lib/core/coreApi.go", "package lib", "package core")
replace_in_file("sandbox/contracts/sandbox/sandbox.go", "package lib", "package sandbox")

# Update imports in sandbox.go
sandbox_go = """package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
)

type SandBox struct {
	Cli    cli.CliApi
	Config config.Config
	Core   core.CoreApi
	Deps   deps.Deps
}
"""
with open("sandbox/contracts/sandbox/sandbox.go", "w") as f:
    f.write(sandbox_go)

# Update imports in cliApi.go
replace_in_file("sandbox/contracts/lib/cli/cliApi.go", "sandbox *SandBox", "sandbox *sandbox.SandBox")
# We need to add import for sandbox in cliApi.go
cli_content = open("sandbox/contracts/lib/cli/cliApi.go").read()
if "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/sandbox" not in cli_content:
    cli_content = cli_content.replace("package cli\n", "package cli\n\nimport \"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/sandbox\"\n")
    with open("sandbox/contracts/lib/cli/cliApi.go", "w") as f:
        f.write(cli_content)

