package backoffice_db

import (
	api "{{.Module}}/sandbox/api"
)

// RemoveBackofficeUserSession deletes one session of one backoffice-user
// record. The declaration generates Add and List for a nested collection but no
// Remove, so it is written here. Removing a session that is already gone is not
// an error.
func RemoveBackofficeUserSession(sandbox *api.Sandbox, db *BackofficeDb, userId int64, sessionId int64) error {
	collection, err := sandbox.Deps.OpinionatedAgnosDatabase.Collection(db.database, "backoffice-user")
	if err != nil {
		return err
	}
	user, ok, failure := collection.FindByID(userId)
	if failure != nil {
		return sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	if !ok {
		return sandbox.Deps.StdDeps.Errorf("backoffice-user %d not found", userId)
	}
	sessions, failure := user.Nested("session")
	if failure != nil {
		return sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	session, ok, failure := sessions.FindByID(sessionId)
	if failure != nil {
		return sandbox.Deps.OpinionatedAgnosDatabase.Fail(failure)
	}
	if !ok {
		return nil
	}
	return sandbox.Deps.OpinionatedAgnosDatabase.Fail(session.Remove())
}
