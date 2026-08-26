package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/parsables"
)

func createThemes(sandbox *api.SandBox, props api.StartProps, configDir string) error {
	themes_dest := configDir + "/themes.yaml"
	themes_conf, err := parsables.NewThemesConf(sandbox, themes_dest)
	if err != nil {
		return err
	}

	err = themes_conf.AddTheme("LibUsage", "lib-usage", "Documentation explaning how to use the lib")
	if err != nil {
		return err
	}

	err = themes_conf.AddTheme("Development", "development", "Documentation explaning how to to. build the project, and how to modify the project")
	if err != nil {
		return err
	}

	err = themes_conf.Persist()
	if err != nil {
		return err
	}
	return nil
}
