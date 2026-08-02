# Rename the Module

## Description
Covers renaming the Go module path, typically the first step when using this repository as a template for a new project.

---

## Workflow
1. Open the `go.mod` file at the project root.
2. Change the module path on the first line (e.g., `module github.com/MateusMoutinhoOrg/Agnos`) to the new module path.
3. Replace every internal import path. Use the IDE's global "Find and Replace", or run:
   - **macOS:**
     ```sh
     find . -type f \( -name '*.go' -o -name '*.md' \) -exec sed -i '' 's|github.com/MateusMoutinhoOrg/Agnos|<your-new-module-path>|g' {} +
     ```
   - **Linux:**
     ```sh
     find . -type f \( -name '*.go' -o -name '*.md' \) -exec sed -i 's|github.com/MateusMoutinhoOrg/Agnos|<your-new-module-path>|g' {} +
     ```
4. Verify and clean up dependencies:
   ```bash
   go mod tidy
   ```
5. Build the project:
   ```bash
   go build ./...
   ```
