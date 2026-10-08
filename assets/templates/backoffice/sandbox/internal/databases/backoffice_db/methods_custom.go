package backoffice_db

import (
	api "{{.Module}}/sandbox/api"
)

// RemoveBackofficeUserSession deletes one session of one backoffice-user
// record. The declaration generates Add and List for a nested collection but no
// Remove, so it is written here. Removing a session that is already gone is not
// an error.
func RemoveBackofficeUserSession(sandbox *api.Sandbox, db *BackofficeDb, userId int64, sessionId int64) error {
	schema, err := sandbox.Deps.OpinionatedAgnosDatabase.Schema(db.handle, "backoffice-user")
	if err != nil {
		return err
	}
	user, ok := schema.FindById(userId)
	if !ok {
		return sandbox.Deps.StdDeps.Errorf("backoffice-user %d not found", userId)
	}
	for _, item := range user.ListAll("session") {
		if item.Id == sessionId {
			return sandbox.Deps.OpinionatedAgnosDatabase.Fail(item.Remove())
		}
	}
	return nil
}
