package build

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Build(deps *deps.Deps, path string) error {
	deps.Printf("build started with path %s \n", path)

	vars := map[string]string{
		"Name": "Agnos",
	}

	content, err := deps.EmbedDeps.RenderTemplate("teste.tmpl", vars)
	if err != nil {
		deps.Printf("failed to render template: %v\n", err)
		return err
	}

	outPath := "teste.txt"
	if path != "" && path != "." {
		outPath = fmt.Sprintf("%s/teste.txt", path)
	}

	err = deps.IoLib.WriteFile(outPath, content)
	if err != nil {
		deps.Printf("failed to write teste.txt: %v\n", err)
		return err
	}

	deps.Printf("successfully rendered template to %s\n", outPath)
	return nil
}
