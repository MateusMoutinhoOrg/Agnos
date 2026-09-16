package disable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// DisableExtensionInternal turns one mechanic off in the project's
// declaration. It removes nothing: from here on agnos simply stops rendering
// what that mechanic owns, and whatever it wrote before is the project's to
// keep or to edit by hand. Deleting those files is what an <x>-purge is for.
func DisableExtensionInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	return utils.SetExtension(sandbox, io, name, false)
}
