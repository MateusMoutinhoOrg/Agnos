package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// modulePath is the current module path, replaced by the rename command.
const modulePath = "github.com/MateusMoutinhoOrg/Agnos-Cli"

// imagesDir is the directory holding the Dockerfiles for each build target.
const imagesDir = "build/images"

// releaseDir is the directory build artifacts are written to.
const releaseDir = "release"

// targetOutput maps each build target name to the filename it produces under
// releaseDir.
var targetOutput = map[string]string{
	"i32":     "windows-i32.exe",
	"rpm":     "x86-rpm.rpm",
	"deb":     "x86-deb.deb",
	"linux86": "linux86.out",
}

// detectProvider probes the host for a working container runtime, trying
// docker first and then podman. It returns the name of the first one whose
// "info" command succeeds, or an error when neither is available.
func detectProvider() (string, error) {
	for _, name := range []string{"docker", "podman"} {
		cmd := exec.Command(name, "info")
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err == nil {
			return name, nil
		}
	}
	return "", fmt.Errorf("neither docker nor podman is installed — install one or pass --provider")
}

// buildTargets validates the requested targets, resolves the container
// provider, creates the release directory, and builds each target in order.
func buildTargets(targets []string, provider string) error {
	for _, t := range targets {
		if _, ok := targetOutput[t]; !ok {
			return fmt.Errorf("unknown target %q (valid: i32, rpm, deb, linux86)", t)
		}
	}

	if provider == "" {
		detected, err := detectProvider()
		if err != nil {
			return err
		}
		provider = detected
		fmt.Printf("auto-detected provider: %s\n", provider)
	} else if provider != "docker" && provider != "podman" {
		return fmt.Errorf("--provider must be docker or podman, got %q", provider)
	}

	if err := os.MkdirAll(releaseDir, 0o755); err != nil {
		return fmt.Errorf("creating release directory: %w", err)
	}

	for _, t := range targets {
		output := targetOutput[t]
		dockerfile := filepath.Join(imagesDir, t+".dockerfile")
		fmt.Printf("building %s...\n", t)
		if err := buildSingleTarget(provider, t, dockerfile, output); err != nil {
			return fmt.Errorf("target %s: %w", t, err)
		}
		fmt.Printf("  -> %s/%s\n", releaseDir, output)
	}

	return nil
}

// buildSingleTarget builds the container image from the target's Dockerfile,
// creates a temporary container without starting it, copies the build
// artifact out to the release directory, and removes the container. The image
// is kept for caching on subsequent runs.
func buildSingleTarget(provider, target, dockerfile, outputFile string) error {
	imageName := "agnos-build-" + target
	containerName := "agnos-tmp-" + target

	// Remove any leftover container from a previous failed run.
	cleanup := exec.Command(provider, "rm", "-f", containerName)
	cleanup.Stdout = nil
	cleanup.Stderr = nil
	_ = cleanup.Run()

	// Build the image from the target's Dockerfile, with the entire
	// repository as the build context so COPY . . works.
	build := exec.Command(provider, "build",
		"-f", dockerfile,
		"-t", imageName,
		".",
	)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("building image: %w", err)
	}

	// Create a container (not started) so we can copy the artifact out.
	create := exec.Command(provider, "create", "--name", containerName, imageName)
	create.Stdout = nil
	create.Stderr = os.Stderr
	if err := create.Run(); err != nil {
		return fmt.Errorf("creating container: %w", err)
	}

	// Copy the built artifact from the container to the release directory.
	src := containerName + ":/output/" + outputFile
	dst := filepath.Join(releaseDir, outputFile)
	cp := exec.Command(provider, "cp", src, dst)
	cp.Stdout = os.Stdout
	cp.Stderr = os.Stderr
	if err := cp.Run(); err != nil {
		_ = exec.Command(provider, "rm", "-f", containerName).Run()
		return fmt.Errorf("copying artifact: %w", err)
	}

	// Clean up the temporary container.
	_ = exec.Command(provider, "rm", "-f", containerName).Run()
	return nil
}

// renameModule replaces every occurrence of the current module path with
// newPath in go.mod and every .go file under the repository root, skipping
// hidden directories, the release directory, and vendor.
func renameModule(newPath string) error {
	fmt.Printf("renaming %s -> %s\n", modulePath, newPath)

	count := 0
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			base := filepath.Base(path)
			if strings.HasPrefix(base, ".") && path != "." {
				return filepath.SkipDir
			}
			if base == releaseDir || base == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		base := filepath.Base(path)
		if ext != ".go" && base != "go.mod" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("reading %s: %w", path, readErr)
		}

		content := string(data)
		if !strings.Contains(content, modulePath) {
			return nil
		}

		replaced := strings.ReplaceAll(content, modulePath, newPath)
		if writeErr := os.WriteFile(path, []byte(replaced), info.Mode()); writeErr != nil {
			return fmt.Errorf("writing %s: %w", path, writeErr)
		}

		count++
		fmt.Printf("  updated: %s\n", path)
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("renamed in %d file(s)\n", count)
	return nil
}
