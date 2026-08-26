package standard

import (
	"os"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// New creates a deps.Deps backed by the standard adapter, ready for lib.New.
// The embedded Keep library writes the tracker's categories and transactions
// under the provided basePath, one file per key; the embedded Verb library
// parses the process's own command line, os.Args[1:]; and Printf writes to
// the process's standard output — this adapter is the opinionated one, so it
// picks the argument vector and the stream itself. Handing the same
// os.Args[1:] to api.Lib.Sandboxmain is what keeps the interface's view of
// the command line and the parser's in agreement. Every asset the library
// asks for is served from the whole compiled-in asset tree, so nothing has to
// exist on disk beside the binary.
//
// It builds the adapter instance and runs every field factory over it, so
// each closure reads the adapter's state at call time. Adding a field to
// deps.Deps means adding its factory call here.
func New(basePath string) deps.Deps {
	adapter := &StandardAdapter{
		args:         os.Args[1:],
		output:       os.Stdout,
		keepBasePath: basePath,
	}
	adapter.Deps.Now = NowFactory(adapter)
	adapter.Deps.Printf = PrintfFactory(adapter)
	adapter.Deps.Error = ErrorFactory(adapter)
	adapter.Deps.Errorf = ErrorfFactory(adapter)
	adapter.Deps.VerbLib = VerbLibFactory(adapter)
	adapter.Deps.KeepLib = KeepLibFactory(adapter)
	adapter.Deps.EmbedDeps = EmbedDepsFactory(adapter)
	adapter.Deps.IoLib = IoLibFactory(adapter)
	adapter.Deps.NewRequest = NewRequestFactory(adapter)
	adapter.Deps.SerializeLib = SerializeLibFactory(adapter)
	return adapter.Deps
}
