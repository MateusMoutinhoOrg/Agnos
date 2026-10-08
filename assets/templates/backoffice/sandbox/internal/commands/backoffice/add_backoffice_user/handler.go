package add_backoffice_user

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// Handle answers `add-backoffice-user`. It generates the user's
// password — so none travels on the command line, where every user of the
// machine and the shell history read it — and adds the user through
// backofficeusers.Add, refused on the same grounds as the add form: a username
// or email already in use, an invalid email. The password is printed once.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	role, ok := backofficeauth.ParseRole(sandbox, input.Role)
	if !ok {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "unknown role "+input.Role)
	}
	password, err := backofficeusers.GeneratePassword(sandbox)
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "failed to generate a password: "+err.Error())
	}

	user, message, err := backofficeusers.Add(sandbox, backofficeusers.Fields{
		Username: input.Username,
		Email:    input.Email,
		Password: password,
		Role:     int64(role),
	})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "failed to add backoffice user: "+err.Error())
	}
	if message != "" {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", message)
	}

	response.Printf("backoffice user %s <%s> created with the %s role\n", user.Username, user.Email, backofficeauth.RoleName(sandbox, role))
	response.Printf("password: %s\n", password)
	response.Printf("It is shown only this once. Change it on /admin/root/set-backoffice-user/%d.\n", user.Id)
	return nil
}
