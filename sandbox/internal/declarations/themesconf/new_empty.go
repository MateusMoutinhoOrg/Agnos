package themesconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func NewEmpty(sandbox *api.Sandbox) *ThemesConf {
	themes_conf := &ThemesConf{
		Themes: make([]Theme, 0),
	}
	BindMethods(sandbox, themes_conf)
	return themes_conf
}
