package backofficedb

import (
	api "{{.Module}}/sandbox/api"
)

// RemoveBackofficeuserSessions deletes one sessions record of one
// backofficeuser record. The declaration generates Add and List for a nested
// collection but no Remove, so it is written here. Removing a record that is
// already gone is not an error.
func RemoveBackofficeuserSessions(sandbox *api.Sandbox, self *Backofficedb, parent_id int64, id int64) error {
	schema, err := sandbox.Deps.OpinatedAgnosDatabase.Schema(self.handle, "backofficeuser")
	if err != nil {
		return err
	}
	parent, ok := schema.FindById(parent_id)
	if !ok {
		return sandbox.Deps.Std.Errorf("backofficeuser %d not found", parent_id)
	}
	for _, item := range parent.ListAll("sessions") {
		if item.Id == id {
			return sandbox.Deps.OpinatedAgnosDatabase.Fail(item.Remove())
		}
	}
	return nil
}
