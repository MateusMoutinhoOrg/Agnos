package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddPageInternal writes one page — assets/frontend/<name>.html — on an
// already open transaction. A page is a file and nothing else: the frontend
// route front-init wrote serves it as soon as it exists, so there is no route
// to declare and nothing for the build to generate.
//
// An existing file is refused rather than kept: everything under
// assets/frontend/ is the project's content, and no editor writes over it.
func AddPageInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.PageProps) error {
	if err := utils.ValidatePageName(sandbox, props.Name); err != nil {
		return err
	}

	has_front, err := utils.ExtensionEnabled(sandbox, io, utils.ExtensionSandboxFront)
	if err != nil {
		return err
	}
	if !has_front {
		return sandbox.Deps.Std.Errorf("the project has no front layer: run front-init before declaring a page")
	}

	page := utils.PageAsset(sandbox, props.Name)
	if io.IsFile(page) {
		return sandbox.Deps.Std.Errorf("page %q already exists: %s", utils.PageName(sandbox, props.Name), page)
	}

	title := sandbox.Deps.Stringsdeps.TrimSpace(props.Title)
	if title == "" {
		title = utils.PageName(sandbox, props.Name)
	}

	return utils.WritePage(sandbox, io, props.Name, title)
}
