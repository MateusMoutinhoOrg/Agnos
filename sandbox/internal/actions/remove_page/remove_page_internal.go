package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemovePageInternal deletes one page — assets/frontend/<name>.html — on an
// already open transaction, and refuses a name with no such file.
func RemovePageInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	if err := utils.ValidatePageName(sandbox, name); err != nil {
		return err
	}

	page := utils.PageAsset(sandbox, name)
	if !io.IsFile(page) {
		return sandbox.Deps.Std.Errorf("page %q not found: %s does not exist", utils.PageName(sandbox, name), page)
	}

	sandbox.Deps.Std.Log("remove-page removing %s \n", page)
	io.RemoveDir(page)

	return nil
}
