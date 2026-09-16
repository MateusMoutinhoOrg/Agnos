package enable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// EnableExtensionInternal turns one mechanic on in the project's declaration.
// It writes nothing else: the build that follows renders every group the
// declaration turns on, so what the mechanic owns appears from there.
func EnableExtensionInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	return utils.SetExtension(sandbox, io, name, true)
}
