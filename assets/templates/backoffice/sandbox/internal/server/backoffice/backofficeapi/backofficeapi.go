package backofficeapi

import (
	"{{.Module}}/sandbox/api"
	serializabledeps "{{.Module}}/sandbox/deps/serializabledeps"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// The documents below are what the /api/admin routes answer, each the JSON
// twin of what sandbox/internal/server/backoffice/backofficerender shows on a page. AddItemToObject and
// AddItemToArray copy a child as it is when added, so every child is built
// whole before it is added to its parent.

// UserJSON is the JSON object a backoffice user is answered as: its id, username,
// email and role by name. The password hash never leaves the server.
func UserJSON(sandbox *api.Sandbox, user backoffice_db.BackofficeUserRecord) *serializabledeps.SerializableObject {
	object := sandbox.Deps.SerializableDeps.CreateObject()
	object.AddItemToObject("id", user.Id)
	object.AddItemToObject("username", user.Username)
	object.AddItemToObject("email", user.Email)
	object.AddItemToObject("role", backofficeauth.RoleName(sandbox, backofficeauth.Role(user.Role)))
	return object
}

// UserResponseJSON is {"user": UserJSON(user)}.
func UserResponseJSON(sandbox *api.Sandbox, user backoffice_db.BackofficeUserRecord) *serializabledeps.SerializableObject {
	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("user", UserJSON(sandbox, user))
	return document
}

// UserListJSON is one page of the user list: the users themselves, the search and
// role it was filtered by, and where the page sits among every page.
func UserListJSON(sandbox *api.Sandbox, listing backofficeusers.Listing) *serializabledeps.SerializableObject {
	users := sandbox.Deps.SerializableDeps.CreateArray()
	for _, user := range listing.Users {
		users.AddItemToArray(UserJSON(sandbox, user))
	}

	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("users", users)
	document.AddItemToObject("search", listing.Search)
	document.AddItemToObject("role", listing.Role)
	document.AddItemToObject("total", listing.Total)
	document.AddItemToObject("page", listing.Page)
	document.AddItemToObject("pages", listing.Pages)
	document.AddItemToObject("limit", listing.Limit)
	return document
}

// OkJSON is {"status": "ok"}, what an action with nothing else to say answers.
func OkJSON(sandbox *api.Sandbox) *serializabledeps.SerializableObject {
	document := sandbox.Deps.SerializableDeps.CreateObject()
	document.AddItemToObject("status", "ok")
	return document
}
