package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// triggerDest is the one file of sandbox/internal/generated/trigger: the
// MatchTrigger the route and the command matchers share.
const triggerDest = utils.TriggerDir + "/MatchTrigger.go"

// GenerateTrigger renders assets/templates/match_trigger.go into
// sandbox/internal/generated/trigger/MatchTrigger.go on every build of a
// project whose cli or server layer is on — the two layers that match on an
// api.Trigger. It is not an asset of the sandbox group because it reads
// Deps.Stringsdeps, which a sandbox without its deps layer does not have, and
// not one of either layer's groups because both of them need it.
func GenerateTrigger(sandbox *api.Sandbox, io *smartio.SmartIO, module string) error {
	vars := map[string]any{
		"Module":        module,
		"GeneratorName": generatorName(sandbox),
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/match_trigger.go", vars, triggerDest)
}
