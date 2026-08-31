package themesconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func NewEmpty(deps *deps.Deps) *ThemesConf {
	themes_conf := &ThemesConf{
		Themes: make([]Theme, 0),
	}
	BindMethods(deps, themes_conf)
	return themes_conf
}
