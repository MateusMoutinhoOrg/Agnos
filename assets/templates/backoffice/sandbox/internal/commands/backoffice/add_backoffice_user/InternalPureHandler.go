package add_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// InternalPureHandler answers `add-backoffice-user`. It generates the user's
// password — so none travels on the command line, where every user of the
// machine and the shell history read it — and adds the user through
// backofficeusers.Add, refused on the same grounds as the add form: a username
// or email already in use, an invalid email. The password is printed once.
func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	role, ok := backofficeauth.ParseRole(sandbox, entries.Role)
	if !ok {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "unknown role "+entries.Role)
	}
	password, err := backofficeusers.GeneratePassword(sandbox)
	if err != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "failed to generate a password: "+err.Error())
	}

	user, message, err := backofficeusers.Add(sandbox, backofficeusers.Fields{
		Username: entries.Username,
		Email:    entries.Email,
		Password: password,
		Role:     int64(role),
	})
	if err != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", "failed to add backoffice user: "+err.Error())
	}
	if message != "" {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", message)
	}

	response.Printf("backoffice user %s <%s> created with the %s role\n", user.Username, user.Email, backofficeauth.RoleName(sandbox, role))
	response.Printf("password: %s\n", password)
	response.Printf("It is shown only this once. Change it on /admin/root/edit-backoffice-user/%d.\n", user.Id)
	return nil
}
