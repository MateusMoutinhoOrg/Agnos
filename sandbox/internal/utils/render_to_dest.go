package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func RenderTemplateToDest(deps *deps.Deps, io *smartio.SmartIO, template_path string, vars interface{}, dest_path string) error {

	content, err := deps.EmbedDeps.RenderTemplate(template_path, vars)
	if err != nil {
		return err
	}

	err = io.WriteFile(dest_path, content)
	if err != nil {
		return err
	}

	return nil
}
