package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/parsables"
)

func createThemes(sandbox *api.SandBox, props api.StartProps, configDir string) error {
	themes_dest := configDir + "/themes.yaml"

	var themes_conf *parsables.ThemesConf
	var err error

	if sandbox.Deps.IoLib.IsFile(themes_dest) {
		content_bytes, fileerror := sandbox.Deps.IoLib.ReadFile(themes_dest)
		if fileerror != nil {
			return fileerror
		}
		themes_conf, err = parsables.NewThemesConf(sandbox, string(content_bytes))
		if err != nil {
			return err
		}
	} else {
		themes_conf = parsables.NewThemesConfEmpty(sandbox)
	}

	err = themes_conf.AddTheme("LibUsage", "lib-usage", "Documentation explaning how to use the lib")
	if err != nil {
		return err
	}

	err = themes_conf.AddTheme("Development", "development", "Documentation explaning how to to. build the project, and how to modify the project")
	if err != nil {
		return err
	}

	rendered := themes_conf.Render()
	err = sandbox.Deps.IoLib.WriteFile(themes_dest, []byte(rendered))
	if err != nil {
		return err
	}
	return nil
}
