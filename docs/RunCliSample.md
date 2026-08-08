# Run CLI Examples

## Description
How to run the CLI examples provided in the repository to understand the tool's behavior from the terminal.

---

## Run CLI Examples

CLI examples are shell scripts that demonstrate how to use the `agnos` command-line tool. They build the binary dynamically and run in an isolated scratch space, so they won't affect your actual data.

### Workflow

1. Browse the `/examples/cliExamples/` directory for a script that matches the workflow you want to learn (e.g., `ManageCategories.sh`).
2. Run the script from the project root:
   ```bash
   bash ./examples/cliExamples/ManageCategories.sh
   ```
3. Read the output. The script will output comments indicating what it's doing before executing the CLI commands.
4. To run all CLI examples sequentially:
   ```bash
   for script in ./examples/cliExamples/*.sh; do bash "$script"; done
   ```
