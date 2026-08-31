package deps

import (
	"time"

	keepdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/keepdeps"

	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/embeddeps"
	iodeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/iodeps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/requestdeps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/serializebles"
	verbdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/verbdeps"
)

type Deps struct {
	Now func() time.Time

	Printf     func(format string, a ...any) (n int, err error)
	Error      func(format string, a ...any) (n int, err error)
	Errorf     func(format string, a ...any) error
	NewVerbLib func(args []string) verbdeps.Lib
	NewKeepLib func(basePath string) keepdeps.Lib
	EmbedDeps  embeddeps.Lib

	IoLib iodeps.Lib

	NewRequest func(url string) requestdeps.Request

	SerializeLib serializibles.Lib
}
