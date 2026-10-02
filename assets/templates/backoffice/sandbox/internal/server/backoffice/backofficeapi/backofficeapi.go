package backofficeapi

import (
	"{{.Module}}/sandbox/api"
	serializables "{{.Module}}/sandbox/deps/serializables"
	"{{.Module}}/sandbox/internal/databases/backofficedb"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// The documents below are what the /api/admin routes answer, each the JSON
// twin of what sandbox/internal/server/backoffice/backofficerender shows on a page. AddItemToObject and
// AddItemToArray copy a child as it is when added, so every child is built
// whole before it is added to its parent.

// User is the JSON object a backoffice user is answered as: its id, username,
// email and role by name. The password hash never leaves the server.
func User(sandbox *api.Sandbox, user backofficedb.BackofficeuserItem) *serializables.SerializibleObject {
	object := sandbox.Deps.Serializables.CreateObject()
	object.AddItemToObject("id", user.Id)
	object.AddItemToObject("username", user.Username)
	object.AddItemToObject("email", user.Email)
	object.AddItemToObject("role", backofficeauth.RoleName(sandbox, backofficeauth.Role(user.Role)))
	return object
}

// UserDocument is {"user": User(user)}.
func UserDocument(sandbox *api.Sandbox, user backofficedb.BackofficeuserItem) *serializables.SerializibleObject {
	document := sandbox.Deps.Serializables.CreateObject()
	document.AddItemToObject("user", User(sandbox, user))
	return document
}

// Listing is one page of the user list: the users themselves, the search and
// role it was filtered by, and where the page sits among every page.
func Listing(sandbox *api.Sandbox, listing backofficeusers.Listing) *serializables.SerializibleObject {
	users := sandbox.Deps.Serializables.CreateArray()
	for _, user := range listing.Users {
		users.AddItemToArray(User(sandbox, user))
	}

	document := sandbox.Deps.Serializables.CreateObject()
	document.AddItemToObject("users", users)
	document.AddItemToObject("search", listing.Search)
	document.AddItemToObject("role", listing.Role)
	document.AddItemToObject("total", listing.Total)
	document.AddItemToObject("page", listing.Page)
	document.AddItemToObject("pages", listing.Pages)
	document.AddItemToObject("limit", listing.Limit)
	return document
}

// Ok is {"status": "ok"}, what an action with nothing else to say answers.
func Ok(sandbox *api.Sandbox) *serializables.SerializibleObject {
	document := sandbox.Deps.Serializables.CreateObject()
	document.AddItemToObject("status", "ok")
	return document
}

// Role is the role column value of the role name a body sent, or a 400 on
// field "role" when no role has that name. The routes' json-schema already
// holds role to the names of backofficeauth.Roles, so the failure only
// guards a schema and a role list drifting apart.
func Role(sandbox *api.Sandbox, name string) (int64, error) {
	role, ok := backofficeauth.ParseRole(sandbox, name)
	if !ok {
		return 0, routeio.Fail(sandbox, api.StatusBadRequest, "role", "role must be root or viewer")
	}
	return int64(role), nil
}
