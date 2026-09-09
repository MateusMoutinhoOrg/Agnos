package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
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
func RemovePageInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateRouteName(deps, name); err != nil {
		return err
	}

	page := utils.PageAsset(deps, name)
	if !io.IsFile(page) {
		return deps.Std.Errorf("page %q not found: %s does not exist (a route with no html is removed with remove-route)",
			utils.RouteIdentifier(deps, name), page)
	}

	if err := removeRouteAction.RemoveRoutePackage(deps, io, name); err != nil {
		return err
	}

	deps.Std.Log("remove-page removing %s \n", page)
	io.RemoveDir(page)

	return nil
}
