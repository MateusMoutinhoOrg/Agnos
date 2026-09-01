package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func RenderTemplateToDest(deps *deps.Deps, template_path string, vars interface{}, dest_path string) error {

	content, err := deps.EmbedDeps.RenderTemplate(template_path, vars)
	if err != nil {
		return err
	}

	err = deps.IoLib.WriteFile(dest_path, content)
	if err != nil {
		return err
	}

	return nil
}
