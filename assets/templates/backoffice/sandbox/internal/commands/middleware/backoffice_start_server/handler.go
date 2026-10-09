package backoffice_start_server

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/commandprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficehttp"
	"{{.Module}}/sandbox/internal/snapshots"
)

// Handle runs in front of `start-server` and answers nothing, so
// start-server runs after it. It reads the secret that signs the backoffice
// sessions from the environment — never from the command line — and the two
// flags the backoffice adds to start-server, onto sandbox.Config, where every
// route reads them. With no secret set it generates one for the run and warns
// that sessions end at the next restart; with one too short the server does
// not start. A snapshot the last run left being created is marked failed, and
// a restore it left unfinished put back, in the background, since no job of
// this run will finish either. It warns when --insecure-http is served on
// every interface, and when .gitignore lists neither store under the
// --database folder of the run.
//
// It is a middleware rather than an edit to start-server's own handler, so
// that file stays the project's and backoffice-purge has nothing to undo in it.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	secret, generated, err := backofficeauth.ReadSecret(sandbox)
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	if generated {
		sandbox.Deps.StdDeps.Eprintf("warning: %s is not set, so a random secret was generated for this run: every backoffice session ends when the server restarts, and no other instance accepts them; set it to a random secret of at least %d characters (openssl rand -hex 32) to keep them\n", backofficeauth.SecretEnv(sandbox), backofficeauth.MinSecretLength)
	}

	sandbox.Config.SessionSecret = secret
	sandbox.Config.AllowXForwardedFor = input.AllowXForwardedFor
	sandbox.Config.InsecureHttp = input.InsecureHttp

	if input.AllowXForwardedFor && backofficehttp.ListensEverywhere(sandbox, input.Addr) {
		sandbox.Deps.StdDeps.Eprintf("warning: X-Forwarded-For is trusted but the server listens on every interface: bind it to the address only the proxy reaches (--addr 127.0.0.1:3000) or firewall the port, or anyone reaching it can forge their ip\n")
	}

	if input.InsecureHttp && backofficehttp.ListensEverywhere(sandbox, input.Addr) {
		sandbox.Deps.StdDeps.Eprintf("warning: --insecure-http serves the backoffice over plain http, and the server listens on every interface: anyone on the network can read the passwords and sessions sent to /admin; for local development bind it to this machine alone (--addr 127.0.0.1:3000)\n")
	}
	warnUnignoredStores(sandbox)

	snapshots.StartRecover(sandbox)
	return nil
}

// warnUnignoredStores warns when the project's .gitignore, which
// backoffice-init wrote, ignores neither the users' store nor the backups'
// under the --database folder this run uses: both hold password hashes and
// every database, and a folder other than the default one is not listed.
func warnUnignoredStores(sandbox *api.Sandbox) {
	io := sandbox.Deps.IoDeps
	strings := sandbox.Deps.StringsDeps
	content, err := io.ReadFile(".gitignore")
	if err != nil {
		return
	}
	dir := strings.Trim(strings.ReplaceAll(sandbox.Config.DatabaseDir, "\\", "/"), "/")
	dir = strings.TrimPrefix(dir, "./")
	lines := map[string]bool{}
	for _, line := range strings.Split(string(content), "\n") {
		lines[strings.Trim(strings.TrimSpace(line), "/")] = true
	}
	for _, store := range []string{snapshots.BackofficeDir, snapshots.BackupDir} {
		if !lines[dir] && !lines[dir+"/"+store] {
			sandbox.Deps.StdDeps.Eprintf("warning: .gitignore does not list /%s/%s: the backoffice users and the backups under --database %s may end up committed; add /%s to it\n", dir, store, sandbox.Config.DatabaseDir, dir)
			return
		}
	}
}
