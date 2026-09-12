package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemovePageInternal deletes a page's route package and the html template it
// renders, on one open transaction.
//
// It refuses a route with no html beside it, which is the whole of what tells
// a page from any other route: remove-route is the editor for those, and this
// one is what also drops the content. The two are one editor per place a page
// holds something, so neither leaves the other half orphaned.
func RemovePageInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateRouteName(sandbox, name); err != nil {
		return err
	}

	page := utils.PageAsset(sandbox, name)
	if !io.IsFile(page) {
		return sandbox.Deps.Std.Errorf("page %q not found: %s does not exist (a route with no html is removed with remove-route)",
			utils.RouteIdentifier(sandbox, name), page)
	}

	if err := removeRouteAction.RemoveRoutePackage(sandbox, io, name); err != nil {
		return err
	}

	sandbox.Deps.Std.Log("remove-page removing %s \n", page)
	io.RemoveDir(page)

	return nil
}
