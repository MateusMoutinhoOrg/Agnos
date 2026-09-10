package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/apishape"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// remoteApiAlias is the import alias the shim reads the remote contract under.
// The copied contract is imported under its own package name, so the two
// spellings of one type never collide in the generated file.
const remoteApiAlias = "remoteapi"

// ShimProps describes the adapter to generate for one remote dep.
type ShimProps struct {
	// Dep is the local name of the copied contract: the directory under
	// sandbox/deps/, the package of the adapter, and the Deps field
	// title-cased.
	Dep string

	// Module is the remote module path.
	Module string

	// Available is the available of the remote repo the shim builds its
	// sandbox from.
	Available string

	// HasDeps reports whether the remote repo has a deps layer at all — a repo
	// with none has `New()` instead of `New(deps)`, and no available to build.
	HasDeps bool
}

// GenerateShim writes adapters/libs/<dep>/<dep>.go: the adapter that builds the
// remote sandbox out of the remote repo's own adapters and hands it over as the
// local contract, plus the converters that carry every value between the two
// copies of the api.
func GenerateShim(deps *deps.Deps, io *smartio.SmartIO, remote *apishape.Api, props ShimProps) error {

	plan, err := apishape.Converters(deps, remote, props.Dep, remoteApiAlias)
	if err != nil {
		return err
	}

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}

	vars := map[string]any{
		"Adapter":      props.Dep,
		"Dep":          props.Dep,
		"Field":        utils.DepField(deps, props.Dep),
		"Module":       module_conf.Module,
		"RemoteModule": props.Module,
		"RemoteApi":    remoteApiAlias,
		"Available":    props.Available,
		"HasDeps":      props.HasDeps,
		"Entry":        plan.Entry,
		"Converters":   plan.Converters,
	}

	return utils.RenderTemplateToDest(deps, io, "templates/remote_shim.go", vars,
		utils.AdapterDir(props.Dep)+"/"+props.Dep+".go")
}
